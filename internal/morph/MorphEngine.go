package morph

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jhonroun/pullenti/internal/MorphForm"
	"github.com/jhonroun/pullenti/internal/morphinternal"
)

//go:embed data/*.dat
var embeddedFS embed.FS

type MorphEngine struct {
	mu           sync.Mutex
	mLazyBuf     *morphinternal.ByteArrayWrapper
	Language     MorphForm.MorphLang
	MRoot        *MorphTreeNode
	MRootReverce *MorphTreeNode
	mRules       []*MorphForm.MorphRule
	mMiscInfos   []*MorphForm.MorphMiscInfo
}

func NewMorphEngine() *MorphEngine {
	return &MorphEngine{
		MRoot:        &MorphTreeNode{},
		MRootReverce: &MorphTreeNode{},
		Language:     MorphForm.MorphLangUnknown,
	}
}

func (me *MorphEngine) GetLazyBuf() *morphinternal.ByteArrayWrapper {
	return me.mLazyBuf
}

func (me *MorphEngine) AddRule(r *MorphForm.MorphRule) {
	me.mRules = append(me.mRules, r)
}

func (me *MorphEngine) GetRule(id int) *MorphForm.MorphRule {
	if id > 0 && id <= len(me.mRules) {
		return me.mRules[id-1]
	}
	return nil
}

func (me *MorphEngine) GetMutRule(id int) *MorphForm.MorphRule {
	return me.GetRule(id)
}

func (me *MorphEngine) GetRuleVar(rid, vid int) *MorphForm.MorphRuleVariant {
	r := me.GetRule(rid)
	if r == nil {
		return nil
	}
	return r.FindVar(int16(vid))
}

func (me *MorphEngine) AddMiscInfo(mi *MorphForm.MorphMiscInfo) {
	if mi.Id == 0 {
		mi.Id = len(me.mMiscInfos) + 1
	}
	me.mMiscInfos = append(me.mMiscInfos, mi)
}

func (me *MorphEngine) GetMiscInfo(id int) *MorphForm.MorphMiscInfo {
	if id > 0 && id <= len(me.mMiscInfos) {
		return me.mMiscInfos[id-1]
	}
	return nil
}

func (me *MorphEngine) Initialize(lang MorphForm.MorphLang, lazy bool) bool {
	if !me.Language.IsUndefined() {
		return false
	}

	me.mu.Lock()
	defer me.mu.Unlock()

	if !me.Language.IsUndefined() {
		return false
	}

	me.Language = lang
	filename := fmt.Sprintf("m_%s.dat", lang.String())
	path := filepath.Join("../internal/data", filename)

	data, err := embeddedFS.ReadFile(path)
	if err != nil {
		me.Language = MorphForm.MorphLangUnknown // сбрасываем при неудаче
		return false
	}

	return me.Deserialize(bytes.NewReader(data), false, lazy)
}

// Deserialize загружает и инициализирует словарь морфологии
func (me *MorphEngine) Deserialize(stream io.Reader, ignoreRevTree, lazy bool) bool {
	if stream == nil {
		return false
	}

	// Декодируем сжатый GZIP-данный поток
	var rawData bytes.Buffer
	err := morphinternal.DeflateGzip(stream, &rawData)
	if err != nil {
		return false
	}

	// Оборачиваем в ByteArrayWrapper
	buf := rawData.Bytes()
	me.mLazyBuf = morphinternal.NewByteArrayWrapper(buf)
	pos := 0

	// Десериализация основного дерева
	if lazy {
		me.MRoot.DeserializeLazy(me.mLazyBuf, me, &pos)
	} else {
		me.MRoot.Deserialize(me.mLazyBuf, &pos)
	}

	if !ignoreRevTree {
		if lazy {
			me.MRootReverce.DeserializeLazy(me.mLazyBuf, me, &pos)
		} else {
			me.MRootReverce.Deserialize(me.mLazyBuf, &pos)
		}
	}

	return true
}

// loadTreeNode загружает ленивый MorphTreeNode из буфера при необходимости.
// Аналог метода _loadTreeNode в C#.
func (me *MorphEngine) loadTreeNode(tn *MorphTreeNode) {
	me.mu.Lock()
	defer me.mu.Unlock()

	pos := tn.LazyPos
	if pos > 0 {
		tn.DeserializeLazy(me.mLazyBuf, me, &pos)
	}
	tn.LazyPos = 0
}

// calcWordsCount рекурсивно подсчитывает количество слов (вариантов) в дереве MorphTreeNode.
func (me *MorphEngine) calcWordsCount(tn *MorphTreeNode) int64 {
	if tn.LazyPos > 0 {
		me.loadTreeNode(tn)
	}

	var res int64 = 0
	if tn.RuleIds != nil {
		for _, id := range tn.RuleIds {
			rule := me.GetRule(int(id))
			if rule != nil && rule.MorphVars != nil {
				for _, v := range rule.MorphVars {
					res += int64(len(v)) // v — это []MorphRuleVariant
				}
			}
		}
	}

	for _, child := range tn.Nodes {
		res += me.calcWordsCount(child)
	}

	return res
}

// CalcTotalWords возвращает общее количество слов в морфологическом дереве.
func (me *MorphEngine) CalcTotalWords() int64 {
	if me.MRoot == nil {
		return 0
	}
	return me.calcWordsCount(me.MRoot)
}

// GetWordform подбирает подходящую словоформу по заданным морфологическим условиям.
func (me *MorphEngine) GetWordform(word string, cla MorphForm.MorphClass, gender MorphForm.MorphGender, cas MorphForm.MorphCase, num MorphForm.MorphNumber, addInfo *MorphForm.MorphWordForm) string {
	tn := me.MRoot
	var find bool
	var res string
	maxCoef := -10

	for i := 0; i <= len(word); i++ {
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}
		if tn.RuleIds != nil {
			var wordBegin, wordEnd string
			if i > 0 {
				wordBegin = word[:i]
			} else {
				wordEnd = word
			}
			if i < len(word) {
				wordEnd = word[i:]
			} else {
				wordBegin = word
			}

			for _, rid := range tn.RuleIds {
				r := me.GetRule(rid)
				if r == nil || !r.ContainsVar(wordEnd) {
					continue
				}
				for _, li := range r.MorphVars {
					for _, v := range li {
						if (cla.Value&v.GetClass().Value) != 0 && v.NormalTail != "" {
							if cas.IsUndefined() {
								if !v.GetCase().IsNominative() && !v.GetCase().IsUndefined() {
									continue
								}
							} else if v.GetCase().And(cas).IsUndefined() {
								continue
							}

							sur := cla.IsProperSurname()
							sur0 := v.GetClass().IsProperSurname()
							if sur || sur0 {
								if sur != sur0 {
									continue
								}
							}

							find = true
							if gender != MorphForm.MorphGenderUndefined {
								if gender.And(v.Gender) == MorphForm.MorphGenderUndefined {
									if num != MorphForm.MorphNumberPlural {
										continue
									}
								}
							}
							if num != MorphForm.MorphNumberUndefined {
								if num.And(v.Number) == MorphForm.MorphNumberUndefined {
									continue
								}
							}

							re := wordBegin + v.Tail
							co := 0
							if addInfo != nil {
								co = me.calcEqCoef(v, addInfo)
							}
							if res == "" || co > maxCoef {
								res = re
								maxCoef = co
							}
							if maxCoef == 0 && wordBegin+v.NormalTail == word {
								return wordBegin + v.Tail
							}
						}
					}
				}
			}
		}
		if tn.Nodes == nil || i >= len(word) {
			break
		}
		ch := rune(word[i])
		next, ok := tn.Nodes[int16(ch)]
		if !ok {
			break
		}
		tn = next
	}
	if find {
		return res
	}

	tn = me.MRootReverce
	tn0 := me.MRootReverce
	for i := len(word) - 1; i >= 0; i-- {
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}
		ch := rune(word[i])
		if tn.Nodes == nil {
			break
		}
		next, ok := tn.Nodes[int16(ch)]
		if !ok {
			break
		}
		tn = next
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}
		if tn.ReverceVariants != nil {
			tn0 = tn
			break
		}
	}
	if tn0 == me.MRootReverce {
		return ""
	}

	for _, mvr := range tn0.ReverceVariants {
		rule := me.GetRule(mvr.RuleId)
		if rule == nil {
			continue
		}
		mv := rule.FindVar(mvr.VariantId)
		if mv == nil {
			continue
		}
		if (mv.GetClass().Value & cla.Value) != 0 {
			if mv.Tail != "" && !MorphForm.EndsWith(word, mv.Tail) {
				continue
			}
			wordBegin := word[:len(word)-len(mv.Tail)]
			for _, liv := range rule.MorphVars {
				for _, v := range liv {
					if (v.GetClass().Value & cla.Value) != 0 {
						sur := cla.IsProperSurname()
						sur0 := v.GetClass().IsProperSurname()
						if sur || sur0 {
							if sur != sur0 {
								continue
							}
						}
						if !cas.IsUndefined() {
							if v.GetCase().And(cas).IsUndefined() && !v.GetCase().IsUndefined() {
								continue
							}
						}
						if num != MorphForm.MorphNumberUndefined {
							if v.Number != MorphForm.MorphNumberUndefined && v.Number.And(num) == MorphForm.MorphNumberUndefined {
								continue
							}
						}
						if gender != MorphForm.MorphGenderUndefined {
							if v.Gender != MorphForm.MorphGenderUndefined && v.Gender.And(gender) == MorphForm.MorphGenderUndefined {
								continue
							}
						}
						if addInfo != nil {
							if me.calcEqCoef(v, addInfo) < 0 {
								continue
							}
						}
						res = wordBegin + v.Tail
						if res == word {
							return word
						}
						return res
					}
				}
			}
		}
	}
	if cla.IsProperSurname() {
		if gender == MorphForm.MorphGenderFeminine && !cas.IsUndefined() && !cas.IsNominative() {
			if strings.HasSuffix(word, "ВА") || strings.HasSuffix(word, "НА") {
				if cas.IsAccusative() {
					return word[:len(word)-1] + "У"
				}
				return word[:len(word)-1] + "ОЙ"
			}
		}
		if gender == MorphForm.MorphGenderFeminine {
			last := rune(word[len(word)-1])
			if last == 'А' || last == 'Я' || last == 'О' {
				return word
			}
			if MorphForm.IsCyrillicVowel(last) {
				return word[:len(word)-1] + "А"
			} else if last == 'Й' && len(word) > 1 {
				return word[:len(word)-2] + "АЯ"
			} else {
				return word + "А"
			}
		}
	}

	return res
}

// _calcEqCoef вычисляет степень соответствия варианта v и формы wf.
// Возвращает -1, если несовместимы, 0 — если есть частичное совпадение, 1 — если полное.
func (me *MorphEngine) calcEqCoef(v *MorphForm.MorphRuleVariant, wf *MorphForm.MorphWordForm) int {
	if wf.Class.Value != 0 {
		if v.GetClass().Value&wf.Class.Value == 0 {
			return -1
		}
	}

	if wf.Misc != nil && v.MiscInfoId != int16(wf.Misc.Id) {
		vi := me.GetMiscInfo(int(v.MiscInfoId))
		if vi == nil {
			return -1
		}

		if vi.Mood().IsDefined() && wf.Misc.Mood().IsDefined() {
			if vi.Mood() != wf.Misc.Mood() {
				return -1
			}
		}

		if vi.Tense().IsDefined() && wf.Misc.Tense().IsDefined() {
			if (vi.Tense() & wf.Misc.Tense()).IsUndefined() {
				return -1
			}
		}

		if vi.Voice().IsDefined() && wf.Misc.Voice().IsDefined() {
			if vi.Voice() != wf.Misc.Voice() {
				return -1
			}
		}

		if vi.Person().IsDefined() && wf.Misc.Person().IsDefined() {
			if (vi.Person() & wf.Misc.Person()).IsUndefined() {
				return -1
			}
		}

		return 0
	}

	if !v.CheckAccord(&wf.MorphBaseInfo, false, false) {
		return -1
	}

	return 1
}

// Comp1 сравнивает две формы слов по числу, роду, падежу и нормальной форме.
func Comp1(r1, r2 *MorphForm.MorphWordForm) bool {
	if r1.Number != r2.Number || r1.Gender != r2.Gender {
		return false
	}
	if !r1.Case.Equals(r2.Case) {
		return false
	}
	if r1.NormalCase != r2.NormalCase {
		return false
	}
	return true
}

// Sort сортирует список словоформ и удаляет дубликаты по определённым правилам.
func (me *MorphEngine) Sort(res *[]*MorphForm.MorphWordForm, word string) {
	if res == nil || len(*res) < 2 {
		return
	}

	// Сортировка пузырьком с использованием me.Compare
	for k := 0; k < len(*res); k++ {
		changed := false
		for i := 0; i < len(*res)-1; i++ {
			if me.Compare((*res)[i], (*res)[i+1]) > 0 {
				(*res)[i], (*res)[i+1] = (*res)[i+1], (*res)[i]
				changed = true
			}
		}
		if !changed {
			break
		}
	}

	// Удаление похожих форм по правилам
	for i := 0; i < len(*res)-1; i++ {
		for j := i + 1; j < len(*res); j++ {
			if Comp1((*res)[i], (*res)[j]) {
				r1 := (*res)[i]
				r2 := (*res)[j]
				if r1.Class.IsAdjective() && r2.Class.IsNoun() && !r1.IsInDictionary() && !r2.IsInDictionary() {
					*res = append((*res)[:j], (*res)[j+1:]...)
					j--
					continue
				}
				if r1.Class.IsNoun() && r2.Class.IsAdjective() && !r1.IsInDictionary() && !r2.IsInDictionary() {
					*res = append((*res)[:i], (*res)[i+1:]...)
					i--
					break
				}
				if r1.Class.IsAdjective() && r2.Class.IsPronoun() {
					*res = append((*res)[:i], (*res)[i+1:]...)
					i--
					break
				}
				if r1.Class.IsPronoun() && r2.Class.IsAdjective() {
					if r2.NormalFull == "ОДИН" || r2.NormalCase == "ОДИН" {
						continue
					}
					*res = append((*res)[:j], (*res)[j+1:]...)
					j--
					continue
				}
			}
		}
	}
}

func (me *MorphEngine) Compare(x, y *MorphForm.MorphWordForm) int {
	// Если одна форма из словаря, а другая нет
	if x.IsInDictionary() && !y.IsInDictionary() {
		return -1
	}
	if !x.IsInDictionary() && y.IsInDictionary() {
		return 1
	}

	// Сравнение по коэффициенту неопределённости
	if x.UndefCoef > 0 {
		if x.UndefCoef > y.UndefCoef*2 {
			return -1
		}
		if x.UndefCoef*2 < y.UndefCoef {
			return 1
		}
	}

	// Сравнение по морфологическому классу
	if !x.Class.Equals(y.GetClass()) {
		if x.Class.IsPreposition() || x.Class.IsConjunction() || x.Class.IsPronoun() || x.Class.IsPersonalPronoun() {
			return -1
		}
		if y.Class.IsPreposition() || y.Class.IsConjunction() || y.Class.IsPronoun() || y.Class.IsPersonalPronoun() {
			return 1
		}
		if x.Class.IsVerb() {
			return 1
		}
		if y.Class.IsVerb() {
			return -1
		}
		if x.Class.IsNoun() {
			return -1
		}
		if y.Class.IsNoun() {
			return 1
		}
	}

	// Сравнение по пользовательскому коэффициенту
	cx := me.calcCoef(x)
	cy := me.calcCoef(y)

	if cx > cy {
		return -1
	}
	if cx < cy {
		return 1
	}

	// Сравнение по числу
	if x.Number == MorphForm.MorphNumberPlural && y.Number != MorphForm.MorphNumberPlural {
		return 1
	}
	if y.Number == MorphForm.MorphNumberPlural && x.Number != MorphForm.MorphNumberPlural {
		return -1
	}

	return 0
}

// calcCoef вычисляет "вес" морфологической формы для сортировки.
// Чем выше значение, тем предпочтительнее форма. Аналог _calcCoef в C#.
func (me *MorphEngine) calcCoef(wf *MorphForm.MorphWordForm) int {
	k := 0

	// Учитываем наличие падежа
	if !wf.Case.IsUndefined() {
		k++
	}
	// Учитываем наличие рода
	if wf.Gender != MorphForm.MorphGenderUndefined {
		k++
	}
	// Учитываем наличие числа
	if wf.Number != MorphForm.MorphNumberUndefined {
		k++
	}
	// Снижаем приоритет для синонимичных форм
	if wf.Misc != nil && wf.Misc.IsSynonymForm() {
		k -= 3
	}

	// Если нормальная форма короткая или пустая — возвращаем текущий коэффициент
	if wf.NormalCase == "" || len(wf.NormalCase) < 4 {
		return k
	}

	runes := []rune(wf.NormalCase)
	last := runes[len(runes)-1]  // последний символ
	last1 := runes[len(runes)-2] // предпоследний символ

	// Специальная логика для прилагательных в единственном числе
	if wf.Class.IsAdjective() && wf.Number != MorphForm.MorphNumberPlural {
		ok := false
		switch wf.Gender {
		case MorphForm.MorphGenderFeminine:
			if last == 'Я' {
				ok = true
			}
		case MorphForm.MorphGenderMasculine:
			if last == 'Й' {
				if last1 == 'И' {
					k++
				}
				ok = true
			}
		case MorphForm.MorphGenderNeuter:
			if last == 'Е' {
				ok = true
			}
		}
		if ok && MorphForm.IsCyrillicVowel(last1) {
			k++
		}
	} else if wf.Class.IsAdjective() && wf.Number == MorphForm.MorphNumberPlural {
		// Специальная логика для прилагательных во множественном числе
		if last == 'Й' || last == 'Е' {
			k++
		}
	}

	return k
}

// GetAllWordforms возвращает все морфологические формы слова.
func (me *MorphEngine) GetAllWordforms(word string) []*MorphForm.MorphWordForm {
	var res []*MorphForm.MorphWordForm
	tn := me.MRoot

	// Основной проход по дереву
	for i := 0; i <= len(word); i++ {
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}

		if tn.RuleIds != nil {
			wordBegin := ""
			wordEnd := ""
			if i > 0 {
				wordBegin = word[:i]
			} else {
				wordEnd = word
			}
			if i < len(word) {
				wordEnd = word[i:]
			} else {
				wordBegin = word
			}

			for _, rid := range tn.RuleIds {
				rule := me.GetRule(rid)
				if rule != nil && rule.ContainsVar(wordEnd) {
					for _, vl := range rule.MorphVars {
						for _, v := range vl {
							wf := MorphForm.NewMorphWordFormFromVariant(v, "", me.GetMiscInfo(int(v.MiscInfoId)))
							wf.NormalCase = wordBegin + v.Tail
							wf.UndefCoef = 0

							// Проверка уникальности формы
							alreadyExists := false
							for _, ex := range res {
								if wf.HasMorphEquals([]*MorphForm.MorphWordForm{ex}) {
									alreadyExists = true
									break
								}
							}
							if !alreadyExists {
								res = append(res, wf)
							}
						}
					}
				}
			}
		}

		if tn.Nodes == nil || i >= len(word) {
			break
		}

		ch := int16(word[i])
		next, ok := tn.Nodes[ch]
		if !ok {
			break
		}
		tn = next
	}

	// Объединение по классу, роду, числу и нормальной форме — агрегируем падеж
	for i := 0; i < len(res); i++ {
		wf := res[i]
		if wf.ContainsAttr("инф.", nil) {
			continue
		}
		cas := wf.Case
		for j := i + 1; j < len(res); j++ {
			wf1 := res[j]
			if wf1.ContainsAttr("инф.", nil) {
				continue
			}
			if wf.Class.Equals(wf1.Class) && wf.Gender == wf1.Gender &&
				wf.Number == wf1.Number && wf.NormalCase == wf1.NormalCase {
				cas = cas.Or(wf1.Case)
				res = append(res[:j], res[j+1:]...)
				j--
			}
		}
		if !cas.Equals(wf.Case) {
			wf.Case = cas
		}
	}

	// Объединение по классу, падежу, числу и нормальной форме — агрегируем род
	for i := 0; i < len(res); i++ {
		wf := res[i]
		if wf.ContainsAttr("инф.", nil) {
			continue
		}
		for j := i + 1; j < len(res); j++ {
			wf1 := res[j]
			if wf1.ContainsAttr("инф.", nil) {
				continue
			}
			if wf.Class.Equals(wf1.Class) && wf.Case.Equals(wf1.Case) &&
				wf.Number == wf1.Number && wf.NormalCase == wf1.NormalCase {
				wf.Gender = wf.Gender.Or(wf1.Gender)
				res = append(res[:j], res[j+1:]...)
				j--
			}
		}
	}

	return res
}

// checkCorrVar проверяет наличие подходящего варианта слова в дереве морфемных правил.
// Аналог метода _checkCorrVar в C#.
func (me *MorphEngine) checkCorrVar(word string, tn *MorphTreeNode, i int) string {
	for ; i <= len(word); i++ {
		// Загружаем узел, если он отложенно десериализован
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}

		// Проверяем текущие правила узла
		if tn.RuleIds != nil {
			var wordBegin, wordEnd string
			if i > 0 {
				wordBegin = word[:i]
			} else {
				wordEnd = word
			}
			if i < len(word) {
				wordEnd = word[i:]
			} else {
				wordBegin = word
			}

			for _, rid := range tn.RuleIds {
				r := me.GetRule(rid)
				if r == nil {
					continue
				}
				// Прямое совпадение окончания
				if r.ContainsVar(wordEnd) {
					return wordBegin + wordEnd
				}

				// Поддержка шаблонов со звёздочками ('*')
				if strings.ContainsRune(wordEnd, '*') {
					for _, v := range r.Tails {
						if len(v) == len(wordEnd) {
							match := true
							for j := 0; j < len(v); j++ {
								if wordEnd[j] == '*' || wordEnd[j] == v[j] {
									continue
								}
								match = false
								break
							}
							if match {
								return wordBegin + v
							}
						}
					}
				}
			}
		}

		// Выход, если дошли до конца слова или нет дочерних узлов
		if tn.Nodes == nil || i >= len(word) {
			break
		}

		ch := rune(word[i])
		if ch != '*' {
			// Если обычный символ, идем по дереву
			child, ok := tn.Nodes[int16(ch)]
			if !ok {
				break
			}
			tn = child
			continue
		}

		// Если встретилась *, перебираем все возможные подстановки
		for k, child := range tn.Nodes {
			// Заменяем * на текущий символ и продолжаем поиск рекурсивно
			ww := []rune(word)
			ww[i] = rune(k)
			res := me.checkCorrVar(string(ww), child, i+1)
			if res != "" {
				return res
			}
		}
		break
	}
	return ""
}

// CorrectWordByMorph пытается скорректировать слово с использованием морфологического дерева.
// Аналог метода CorrectWordByMorph в C#.
func (me *MorphEngine) CorrectWordByMorph(word string, oneVar bool) []string {
	vars := []string{}
	runes := []rune(word)
	tmp := make([]rune, len(runes)+1) // запас на случай вставки

	// 1. Заменяем каждый символ на '*'
	for ch := 0; ch < len(runes); ch++ {
		copy(tmp, runes)
		tmp[ch] = '*'
		variant := me.checkCorrVar(string(tmp[:len(runes)]), me.MRoot, 0)
		if variant != "" && !morphinternal.Contains(vars, variant) {
			vars = append(vars, variant)
			if oneVar {
				return vars
			}
		}
	}

	// 2. Пробуем вставить '*' в каждую позицию
	if len(vars) == 0 || !oneVar {
		for ch := 1; ch < len(runes); ch++ {
			copy(tmp, runes[:ch])
			tmp[ch] = '*'
			copy(tmp[ch+1:], runes[ch:])
			variant := me.checkCorrVar(string(tmp[:len(runes)+1]), me.MRoot, 0)
			if variant != "" && !morphinternal.Contains(vars, variant) {
				vars = append(vars, variant)
				if oneVar {
					return vars
				}
			}
		}
	}

	// 3. Пробуем удалить каждый символ
	if len(vars) == 0 || !oneVar {
		for ch := 0; ch < len(runes)-1; ch++ {
			copy(tmp, runes[:ch])
			copy(tmp[ch:], runes[ch+1:])
			variant := me.checkCorrVar(string(tmp[:len(runes)-1]), me.MRoot, 0)
			if variant != "" && !morphinternal.Contains(vars, variant) {
				vars = append(vars, variant)
				if oneVar {
					return vars
				}
			}
		}
	}

	if len(vars) == 0 {
		return nil
	}
	return vars
}

// ProcessProperVariants ищет варианты собственных имён (имена или гео-объекты) по слову в обратном дереве.
func (me *MorphEngine) ProcessProperVariants(word string, res *[]*MorphForm.MorphWordForm, geo bool) {
	tn := me.MRootReverce
	var nodesWithVars []*MorphTreeNode

	// Обход слова с конца — поиск узлов, содержащих ReverceVariants.
	for i := len(word) - 1; i >= 0; i-- {
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}

		ch := int16(word[i])
		if tn.Nodes == nil {
			break
		}
		next, ok := tn.Nodes[ch]
		if !ok {
			break
		}
		tn = next
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}
		if tn.ReverceVariants != nil {
			nodesWithVars = append(nodesWithVars, tn)
		}
	}

	// Если не найдено ни одного узла с вариантами — завершить
	if len(nodesWithVars) == 0 {
		return
	}

	// Перебор найденных узлов, начиная с самого глубокого
	for i := len(nodesWithVars) - 1; i >= 0; i-- {
		tn := nodesWithVars[i]
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}

		ok := false
		for _, vr := range tn.ReverceVariants {
			v := me.GetRuleVar(vr.RuleId, int(vr.VariantId))
			if v == nil {
				continue
			}
			// Фильтрация по типу: гео-объект или фамилия
			if geo && !v.GetClass().IsProperGeo() {
				continue
			}
			if !geo && !v.GetClass().IsProperSurname() {
				continue
			}

			// Создание формы слова
			r := MorphForm.NewMorphWordFormFromVariant(v, word, me.GetMiscInfo(int(v.MiscInfoId)))
			if !r.HasMorphEquals(*res) {
				r.UndefCoef = vr.Coef
				*res = append(*res, r)
			}
			ok = true
		}
		if ok {
			break // если нашёлся хотя бы один вариант — остальные не смотрим
		}
	}
}

// ProcessSurnameVariants вызывает ProcessProperVariants с флагом geo = false.
// Используется для обработки вариантов фамилий.
func (me *MorphEngine) ProcessSurnameVariants(word string, res *[]*MorphForm.MorphWordForm) {
	me.ProcessProperVariants(word, res, false)
}

// ProcessGeoVariants вызывает ProcessProperVariants с флагом geo = true.
// Используется для обработки географических названий.
func (me *MorphEngine) ProcessGeoVariants(word string, res *[]*MorphForm.MorphWordForm) {
	me.ProcessProperVariants(word, res, true)
}

// _getAllWordsByClass рекурсивно собирает все слова, соответствующие заданному классу (cla),
// начиная с узла tn, и добавляет их в результат res. tmp — накопитель текущего слова.
func (me *MorphEngine) _getAllWordsByClass(
	cla MorphForm.MorphClass,
	tn *MorphTreeNode,
	res *[]*MorphForm.MorphWordForm,
	tmp *strings.Builder,
) {
	// Загружаем узел из ленивого дерева, если требуется
	if tn.LazyPos > 0 {
		me.loadTreeNode(tn)
	}

	// Если есть правила, проверяем каждое
	if tn.RuleIds != nil {
		for _, rid := range tn.RuleIds {
			r := me.GetRule(rid)
			if r == nil {
				continue
			}
			for i, v := range r.MorphVars {
				var wf *MorphForm.MorphWordForm
				for _, vv := range v {
					// Проверяем соответствие классу
					if vv.GetClass().And(cla).IsUndefined() {
						continue
					}
					// Пропускаем, если падеж задан и не именительный
					if !vv.GetCase().IsUndefined() && !vv.GetCase().IsNominative() {
						continue
					}
					// Пропускаем, если число не единственное и определено
					if vv.Number != MorphForm.MorphNumberUndefined && vv.Number.And(MorphForm.MorphNumberSingular) == MorphForm.MorphNumberUndefined {
						continue
					}

					// Создаём словоформу, копируя данные из варианта
					wf = &MorphForm.MorphWordForm{}
					wf.CopyFromVariant(vv)
					wf.NormalCase = tmp.String()
					if i < len(r.Tails) && r.Tails[i] != "" {
						wf.NormalCase += r.Tails[i]
					}
					*res = append(*res, wf)
					break // только первый подходящий вариант
				}
			}
		}
	}

	// Рекурсивно обходим потомков
	if tn.Nodes != nil {
		for ch, child := range tn.Nodes {
			tmp.WriteRune(rune(ch))
			me._getAllWordsByClass(cla, child, res, tmp)
			s := tmp.String()
			tmp.Reset()
			tmp.WriteString(s[:len(s)-1])
		}
	}
}

// GetAllWordsByClass перебирает всё морфологическое дерево, начиная с корня,
// и находит все слова, соответствующие заданному грамматическому классу (cla),
// добавляя найденные формы в список res.
// Временный буфер tmp используется для построения слова по мере обхода дерева.
// Аналог метода GetAllWordsByClass в C# (Pullenti).
func (me *MorphEngine) GetAllWordsByClass(cla MorphForm.MorphClass, res *[]*MorphForm.MorphWordForm, tmp *strings.Builder) {
	me._getAllWordsByClass(cla, me.MRoot, res, tmp)
}

// ProcessResult обрабатывает список вариантов правила (mvs),
// создавая на их основе MorphWordForm и добавляя их в результат res,
// если такая форма ещё не присутствует.
// Если mv.NormalTail или mv.FullNormalTail указаны и не начинаются с '-',
// они добавляются к wordBegin для формирования нормальных форм.
func (me *MorphEngine) ProcessResult(res *[]*MorphForm.MorphWordForm, wordBegin string, mvs []*MorphForm.MorphRuleVariant) {
	for _, mv := range mvs {
		r := MorphForm.NewMorphWordFormFromVariant(mv, "", me.GetMiscInfo(int(mv.MiscInfoId)))

		// Формируем NormalCase — нормальную форму слова
		if mv.NormalTail != "" && mv.NormalTail[0] != '-' {
			r.NormalCase = wordBegin + mv.NormalTail
		} else {
			r.NormalCase = wordBegin
		}

		// Формируем NormalFull — полную нормальную форму слова
		if mv.FullNormalTail != "" && mv.FullNormalTail[0] != '-' {
			r.NormalFull = wordBegin + mv.FullNormalTail
		} else {
			r.NormalFull = wordBegin
		}

		// Проверка на уникальность и добавление в результат
		if !r.HasMorphEquals(*res) {
			r.UndefCoef = 0
			*res = append(*res, r)
		}
	}
}

// Метод Process(string word, bool ignoreNoDict = false) — это основной алгоритм морфологического анализа слова в Pullenti. Он работает по следующему сценарию:
// 🔧 Этапы разбора слова:
// 1. Проверка на пустоту и отсутствие гласных:
// Если слово пустое или не содержит гласных — сразу return null.
// 2. Обход основного дерева (m_Root):
//   - Поочерёдно разделяет слово на wordBegin и wordEnd.
//   - Находит в дереве tn.RuleIds — идентификаторы правил.
//   - Из каждого правила берутся варианты (GetVars(wordEnd)).
//   - Каждому варианту создаётся MorphWordForm, склеиваются нормальные формы и добавляются в результат.
//
// 3. Обработка дублирующихся форм:
// Сливаются формы с одинаковым Class, Gender, Number, NormalCase — объединяются падежи и род.
// 4. Флаг needTestUnknownVars:
// - Устанавливается в true, если все формы — малозначимые, например, только глаголы или местоимения.
// - Если есть нормальные существительные, прилагательные и т.п. — флаг сбрасывается.
// 5. Обход обратного дерева (m_RootReverce):
//   - Ищутся ReverceVariants (варианты, применимые к концовке слова).
//   - Если ни одна форма из словаря не соответствует, но есть вариант с аналогичной морфологией — добавляется как UndefCoef.
//
// 6. Удаление неправильных гео-форм для "ПРИ":
// Если слово — "ПРИ", исключаются формы с IsProperGeo.
// 7. Финальная сортировка (Sort):
//   - Формы сортируются по критериям:
//   - Присутствие в словаре (IsInDictionary)
//   - Морфологическая определённость (_calcCoef)
//   - Предпочтение существительных перед глаголами
//
// 8. Установка языка и нормализация предлогов.
func (me *MorphEngine) Process(word string, ignoreNoDict bool) []*MorphForm.MorphWordForm {
	if word == "" {
		return nil
	}
	var res []*MorphForm.MorphWordForm
	if len(word) > 1 {
		found := false
		for _, ch := range word {
			if MorphForm.IsCyrillicVowel(ch) || MorphForm.IsLatinVowel(ch) {
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}

	var tn = me.MRoot
	for i := 0; i <= len(word); i++ {
		if tn.LazyPos > 0 {
			me.loadTreeNode(tn)
		}
		if tn.RuleIds != nil {
			var wordBegin, wordEnd string
			if i == 0 {
				wordEnd = word
			} else if i < len(word) {
				wordEnd = word[i:]
			}
			if res == nil {
				res = []*MorphForm.MorphWordForm{}
			}
			for _, rid := range tn.RuleIds {
				r := me.GetRule(rid)
				mvs := r.GetVars(wordEnd)
				if mvs == nil {
					continue
				}
				if wordBegin == "" {
					if i == len(word) {
						wordBegin = word
					} else if i > 0 {
						wordBegin = word[:i]
					}
				}
				me.ProcessResult(&res, wordBegin, mvs)
			}
		}
		if tn.Nodes == nil || i >= len(word) {
			break
		}
		ch := int16(word[i])
		next, ok := tn.Nodes[ch]
		if !ok {
			break
		}
		tn = next
	}

	// Анализ на необходимость теста с реверсным деревом
	needTestUnknownVars := true
	var toFirst *MorphForm.MorphWordForm
	for _, r := range res {
		if r.Class.IsPronoun() || r.Class.IsNoun() || r.Class.IsAdjective() ||
			(r.Class.IsMisc() && r.Class.IsConjunction()) || r.Class.IsPreposition() {
			needTestUnknownVars = false
		} else if r.Class.IsAdverb() && r.NormalCase != "" {
			if !MorphForm.EndsWithEx(r.NormalCase, "О", "А") || r.NormalCase == "МНОГО" {
				needTestUnknownVars = false
			}
		} else if r.Class.IsVerb() && len(res) > 1 {
			ok := false
			for _, rr := range res {
				if rr != r && !rr.Class.Equals(r.Class) {
					ok = true
					break
				}
			}
			if ok && !MorphForm.EndsWith(word, "ИМ") {
				needTestUnknownVars = false
			}
		}

		if len(res) > 1 && toFirst == nil {
			if morphinternal.InStrArr(r.NormalFull, "КОПИЯ", "ПОЛК", "СУД", "ПАРК", "БАНК", "ПОЛОЖЕНИЕ") ||
				morphinternal.InStrArr(r.NormalCase, "МОРЕ", "МАРИЯ", "ВЕТЕР", "КИЕВ") {
				toFirst = r
			}
		}
	}

	if toFirst != nil && res[0] != toFirst {
		res = removeFromSlice(res, toFirst)
		res = append([]*MorphForm.MorphWordForm{toFirst}, res...)
	}

	if needTestUnknownVars && MorphForm.IsCyrillicChar(rune(word[0])) {
		gl, sog := 0, 0
		for _, ch := range word {
			if MorphForm.IsCyrillicVowel(ch) {
				gl++
			} else {
				sog++
			}
		}
		if gl < 2 || sog < 2 {
			needTestUnknownVars = false
		}
	}

	if needTestUnknownVars && len(res) == 1 {
		r := res[0]
		if r.Class.IsVerb() {
			if (r.Misc.Contains("н.вр.") && r.Misc.Contains("нес.в.") && !r.Misc.Contains("страд.з.")) ||
				(r.Misc.Contains("б.вр.") && r.Misc.Contains("сов.в.")) ||
				(r.Misc.Contains("инф.") && r.Misc.Contains("сов.в.")) ||
				(r.NormalCase != "" && MorphForm.EndsWith(r.NormalCase, "СЯ")) {
				needTestUnknownVars = false
			}
		}
		if r.Class.IsUndefined() && r.Misc.Contains("прдктв.") {
			needTestUnknownVars = false
		}
	}

	if needTestUnknownVars {
		if me.MRootReverce == nil || ignoreNoDict {
			return res
		}
		tn = me.MRootReverce
		tn0 := me.MRootReverce
		i := len(word) - 1
		for ; i >= 0; i-- {
			if tn.LazyPos > 0 {
				me.loadTreeNode(tn)
			}
			ch := int16(word[i])
			if tn.Nodes == nil {
				break
			}
			next, ok := tn.Nodes[ch]
			if !ok {
				break
			}
			tn = next
			if tn.LazyPos > 0 {
				me.loadTreeNode(tn)
			}
			if tn.ReverceVariants != nil {
				tn0 = tn
				break
			}
		}
		if tn0 != me.MRootReverce {
			gl := i < 4
			for j := i; j >= 0; j-- {
				if MorphForm.IsCyrillicVowel(rune(word[j])) || MorphForm.IsLatinVowel(rune(word[j])) {
					gl = true
					break
				}
			}
			if gl {
				for _, mvref := range tn0.ReverceVariants {
					mv := me.GetRuleVar(mvref.RuleId, int(mvref.VariantId))
					if mv == nil {
						continue
					}
					if !mv.GetClass().IsVerb() && !mv.GetClass().IsAdjective() && !mv.GetClass().IsNoun() &&
						!mv.GetClass().IsProperSurname() && !mv.GetClass().IsProperGeo() && !mv.GetClass().IsProperSecname() {
						continue
					}
					ok := false
					for _, rr := range res {
						if rr.IsInDictionary() {
							if rr.Class.Equals(mv.GetClass()) || rr.Class.IsNoun() {
								ok = true
								break
							}
							if !mv.GetClass().IsAdjective() && rr.Class.IsVerb() {
								ok = true
								break
							}
						}
					}
					if ok {
						continue
					}
					if mv.Tail != "" && !MorphForm.EndsWith(word, mv.Tail) {
						continue
					}
					r := MorphForm.NewMorphWordFormFromVariant(mv, word, me.GetMiscInfo(int(mv.MiscInfoId)))
					r.UndefCoef = mvref.Coef
					if !r.HasMorphEquals(res) {
						res = append(res, r)
					}
				}
			}
		}
	}

	if word == "ПРИ" && res != nil {
		for i := len(res) - 1; i >= 0; i-- {
			if res[i].Class.IsProperGeo() {
				res = append(res[:i], res[i+1:]...)
			}
		}
	}

	if len(res) == 0 {
		return nil
	}

	me.Sort(&res, word)

	for _, v := range res {
		if v.NormalCase == "" {
			v.NormalCase = word
		}
		if v.Class.IsVerb() && v.NormalFull == "" && MorphForm.EndsWith(v.NormalCase, "ТЬСЯ") {
			v.NormalFull = v.NormalCase[:len(v.NormalCase)-2]
		}
		v.SetLanguage(me.Language)
		if v.Class.IsPreposition() {
			v.NormalCase = MorphForm.NormalizePreposition(v.NormalCase)
		}
	}

	var mc MorphForm.MorphClass
	for i := len(res) - 1; i >= 0; i-- {
		if !res[i].IsInDictionary() && res[i].Class.IsAdjective() && len(res) > 1 {
			if res[i].Misc.Contains("к.ф.") || res[i].Misc.Contains("неизм.") {
				res = append(res[:i], res[i+1:]...)
				continue
			}
		}
		if res[i].IsInDictionary() {
			mc.Value |= res[i].Class.Value
		}
	}

	if mc.Equals(MorphForm.MorphClassVerb) && len(res) > 1 {
		for _, r := range res {
			if r.UndefCoef > 100 && r.Class.Equals(MorphForm.MorphClassAdjective) {
				r.UndefCoef = 0
			}
		}
	}

	if len(res) == 0 {
		return nil
	}

	return res
}

func removeFromSlice[T comparable](slice []T, target T) []T {
	for i, v := range slice {
		if v == target {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
