package morph

import (
	"errors"
	"strings"
	"sync"

	"github.com/jhonroun/pullenti/internal/MorphForm"
	"github.com/jhonroun/pullenti/internal/morphinternal"
)

type InnerMorphology struct {
	Engines     map[string]*MorphEngine
	Initialized map[string]bool
	LazyLoad    bool
	LoadAll     bool
	Progress    int
	ProgressMax int
	ProgressMsg string

	EngineUa *MorphEngine
	EngineBy *MorphEngine
	EngineKz *MorphEngine
	EngineRu *MorphEngine
	EngineEn *MorphEngine

	lock        sync.Mutex
	lastPercent int
}

func (im *InnerMorphology) Initialize() *InnerMorphology {
	im = &InnerMorphology{
		Engines:     make(map[string]*MorphEngine),
		Initialized: make(map[string]bool),
		LazyLoad:    true,
		LoadAll:     false,
		Progress:    0,
		ProgressMax: 0,
		ProgressMsg: "",

		EngineRu: NewMorphEngine(),
		EngineEn: NewMorphEngine(),
		EngineUa: NewMorphEngine(),
		EngineBy: NewMorphEngine(),
		EngineKz: NewMorphEngine(),
	}

	im.Engines["ru"] = im.EngineRu
	im.Engines["en"] = im.EngineEn
	im.Engines["ua"] = im.EngineUa
	im.Engines["by"] = im.EngineBy
	im.Engines["kz"] = im.EngineKz

	im.Initialized["ru"] = false
	im.Initialized["en"] = false
	im.Initialized["ua"] = false
	im.Initialized["by"] = false
	im.Initialized["kz"] = false

	return im
}

// LoadedLanguages возвращает объединённый язык всех инициализированных движков
func (im *InnerMorphology) LoadedLanguages() MorphForm.MorphLang {
	lang := im.EngineRu.Language.Or(im.EngineEn.Language)
	lang = lang.Or(im.EngineUa.Language)
	lang = lang.Or(im.EngineBy.Language)
	lang = lang.Or(im.EngineKz.Language)
	return lang
}

// LoadLanguages загружает языки при необходимости (в потокобезопасном режиме)
func (im *InnerMorphology) LoadLanguages(langs MorphForm.MorphLang, lazyLoad bool) error {
	if langs.IsRu() && !im.EngineRu.Language.IsRu() {
		im.lock.Lock()
		if !im.EngineRu.Language.IsRu() {
			if !im.EngineRu.Initialize(MorphForm.MorphLangRU, lazyLoad) {
				im.lock.Unlock()
				return errors.New("not found resource file m_ru.dat in Morphology")
			}
		}
		im.lock.Unlock()
	}
	if langs.IsEn() && !im.EngineEn.Language.IsEn() {
		im.lock.Lock()
		if !im.EngineEn.Language.IsEn() {
			if !im.EngineEn.Initialize(MorphForm.MorphLangEN, lazyLoad) {
				im.lock.Unlock()
				return errors.New("not found resource file m_en.dat in Morphology")
			}
		}
		im.lock.Unlock()
	}
	if langs.IsUa() && !im.EngineUa.Language.IsUa() {
		im.lock.Lock()
		if !im.EngineUa.Language.IsUa() {
			im.EngineUa.Initialize(MorphForm.MorphLangUA, lazyLoad)
		}
		im.lock.Unlock()
	}
	if langs.IsBy() && !im.EngineBy.Language.IsBy() {
		im.lock.Lock()
		if !im.EngineBy.Language.IsBy() {
			im.EngineBy.Initialize(MorphForm.MorphLangBY, lazyLoad)
		}
		im.lock.Unlock()
	}
	if langs.IsKz() && !im.EngineKz.Language.IsKz() {
		im.lock.Lock()
		if !im.EngineKz.Language.IsKz() {
			im.EngineKz.Initialize(MorphForm.MorphLangKZ, lazyLoad)
		}
		im.lock.Unlock()
	}
	return nil
}

// UnloadLanguages выгружает заданные языки
func (im *InnerMorphology) UnloadLanguages(langs MorphForm.MorphLang) {
	if langs.IsRu() && im.EngineRu.Language.IsRu() {
		im.EngineRu = NewMorphEngine()
	}
	if langs.IsEn() && im.EngineEn.Language.IsEn() {
		im.EngineEn = NewMorphEngine()
	}
	if langs.IsUa() && im.EngineUa.Language.IsUa() {
		im.EngineUa = NewMorphEngine()
	}
	if langs.IsBy() && im.EngineBy.Language.IsBy() {
		im.EngineBy = NewMorphEngine()
	}
	if langs.IsKz() && im.EngineKz.Language.IsKz() {
		im.EngineKz = NewMorphEngine()
	}
	// Принудительная сборка мусора
	// runtime.GC() — по желанию
}

// SetEngines задаёт общий MorphEngine на все языки (обычно для отладки)
func (im *InnerMorphology) SetEngines(engine *MorphEngine) {
	if engine != nil {
		im.EngineRu = engine
		im.EngineEn = engine
		im.EngineUa = engine
		im.EngineBy = engine
	}
}

// OnProgress вызывает callback с текущим прогрессом
func (im *InnerMorphology) OnProgress(val, max int) {
	p := val
	if max > 0xFFFF {
		p = p / (max / 100)
	} else {
		p = (p * 100) / max
	}

	im.lastPercent = p
}

// Run выполняет морфологический разбор текста с автоматическим определением языка
func (im *InnerMorphology) Run(text string, onlyTokenizing bool, dlang *MorphForm.MorphLang, goodText bool, progress func(done, total int)) []*MorphForm.MorphToken {
	if len(text) == 0 {
		return nil
	}

	twr := morphinternal.NewTextWrapper(text, goodText)
	twrch := twr.Chars
	res := make([]*MorphForm.MorphToken, 0, len(text)/6)
	uniLex := make(map[string]*MorphForm.UniLexWrap)

	var term0 string
	var pureRusWords, pureUkrWords, pureByWords, pureKzWords int
	var totRusWords, totUkrWords, totByWords, totKzWords int

	for i := 0; i < twr.Length; i++ {
		ty := im.GetCharTyp(twrch.At(i))
		if ty == 0 {
			continue
		}

		var j int
		if ty > 2 {
			j = i + 1
		} else {
			j = i + 1
			for ; j < twr.Length; j++ {
				if im.GetCharTyp(twrch.At(j)) != ty {
					break
				}
			}
		}

		wstr := text[i:j]
		var term string
		if goodText {
			term = wstr
		} else {
			trstr := MorphForm.TransliteralCorrection(wstr, term0, false)
			term = MorphForm.CorrectWord(trstr)
		}

		if term == "" {
			i = j - 1
			continue
		}

		if strings.HasPrefix(term, "ЗАВ") {
			term1 := term[3:]
			if strings.HasPrefix(term1, "ОТДЕЛ") || strings.HasPrefix(term1, "ЛАБОРАТ") || strings.HasPrefix(term1, "КАФЕДР") {
				term = "ЗАВ"
				j = i + 3
			}
		}

		lang := MorphForm.GetWordLang(term)
		if len(term) > 2 {
			switch {
			case lang.Equals(MorphForm.MorphLangUA):
				pureUkrWords++
			case lang.Equals(MorphForm.MorphLangRU):
				pureRusWords++
			case lang.Equals(MorphForm.MorphLangBY):
				pureByWords++
			case lang.Equals(MorphForm.MorphLangKZ):
				pureKzWords++
			}
		}
		if lang.IsRu() {
			totRusWords++
		}
		if lang.IsUa() {
			totUkrWords++
		}
		if lang.IsBy() {
			totByWords++
		}
		if lang.IsKz() {
			totKzWords++
		}
		if ty == 1 {
			term0 = term
		}

		var lemmas *MorphForm.UniLexWrap
		if ty == 1 && !onlyTokenizing {
			if u, ok := uniLex[term]; ok {
				lemmas = u
			} else {
				nuni := MorphForm.NewUniLexWrap(*lang)
				uniLex[term] = nuni
				lemmas = nuni
			}
		}

		tok := &MorphForm.MorphToken{
			Term:      term,
			BeginChar: i,
			EndChar:   j - 1,
			Tag:       lemmas,
		}
		res = append(res, tok)
		i = j - 1
	}

	// Определение языка по частоте слов
	defLang := MorphForm.NewMorphLang()
	if dlang != nil {
		defLang.Value = dlang.Value
	}

	switch {
	case pureRusWords > pureUkrWords && pureRusWords > pureByWords && pureRusWords > pureKzWords:
		defLang = MorphForm.MorphLangRU
	case totRusWords > totUkrWords && (totRusWords > totByWords || (totRusWords == totByWords && pureByWords == 0)) && totRusWords > totKzWords:
		defLang = MorphForm.MorphLangRU
	case pureUkrWords > pureRusWords && pureUkrWords > pureByWords && pureUkrWords > pureKzWords:
		defLang = MorphForm.MorphLangUA
	case totUkrWords > totRusWords && totUkrWords > totByWords && totUkrWords > totKzWords:
		defLang = MorphForm.MorphLangUA
	case pureKzWords > pureRusWords && (totKzWords+pureKzWords) > totRusWords && pureKzWords > pureUkrWords && pureKzWords > pureByWords:
		defLang = MorphForm.MorphLangKZ
	case totKzWords > totRusWords && totKzWords > totUkrWords && totKzWords > totByWords:
		defLang = MorphForm.MorphLangKZ
	case pureByWords > pureRusWords && pureByWords > pureUkrWords && pureByWords > pureKzWords:
		if pureByWords < 10 && dlang != nil && dlang.IsRu() {
			// ничего
		} else if pureByWords > 5 {
			defLang = MorphForm.MorphLangBY
		}
	case totByWords > totRusWords && totByWords > totUkrWords && totByWords > totKzWords:
		if totRusWords > 10 && totByWords > (totRusWords+20) {
			defLang = MorphForm.MorphLangBY
		} else if totRusWords == 0 || totByWords >= (totRusWords*2) {
			defLang = MorphForm.MorphLangBY
		}
	}

	// Повторный анализ, если язык не определён или сомнительный
	if (defLang.IsUndefined() || defLang.IsUa()) && totRusWords > 0 {
		if (totUkrWords > totRusWords && im.EngineUa.Language.IsUa()) ||
			(totByWords > totRusWords && im.EngineBy.Language.IsBy()) ||
			(totKzWords > totRusWords && im.EngineKz.Language.IsKz()) {
			cou0 := 0
			totRusWords = 0
			totByWords = 0
			totUkrWords = 0
			totKzWords = 0
			for term, kp := range uniLex {
				lang := MorphForm.NewMorphLang()
				kp.WordForms = im.ProcessOneWord(term, &lang)
				if kp.WordForms != nil {
					for _, wf := range kp.WordForms {
						lang = lang.Or(wf.GetLanguage())
					}
				}
				kp.Lang = lang
				if lang.IsRu() {
					totRusWords++
				}
				if lang.IsUa() {
					totUkrWords++
				}
				if lang.IsBy() {
					totByWords++
				}
				if lang.IsKz() {
					totKzWords++
				}
				if lang.IsCyrillic() {
					cou0++
				}
				if cou0 >= 100 {
					break
				}
			}

			if totRusWords > (totByWords/2) && totRusWords > (totUkrWords/2) {
				defLang = MorphForm.MorphLangRU
			} else if totUkrWords > (totRusWords/2) && totUkrWords > (totByWords/2) {
				defLang = MorphForm.MorphLangUA
			} else if totByWords > (totRusWords/2) && totByWords > (totUkrWords/2) {
				defLang = MorphForm.MorphLangBY
			}
		} else if defLang.IsUndefined() {
			defLang = MorphForm.MorphLangRU
		}
	}

	// Вызов ProcessOneWord для каждого слова
	cou := 0
	totRusWords, totByWords, totUkrWords, totKzWords = 0, 0, 0, 0
	for term, kp := range uniLex {
		lang := defLang.Clone()
		if lang.IsUndefined() {
			if totRusWords > totByWords && totRusWords > totUkrWords && totRusWords > totKzWords {
				lang = MorphForm.MorphLangRU
			} else if totUkrWords > totRusWords && totUkrWords > totByWords && totUkrWords > totKzWords {
				lang = MorphForm.MorphLangUA
			} else if totByWords > totRusWords && totByWords > totUkrWords && totByWords > totKzWords {
				lang = MorphForm.MorphLangBY
			} else if totKzWords > totRusWords && totKzWords > totUkrWords && totKzWords > totByWords {
				lang = MorphForm.MorphLangKZ
			}
		}
		kp.WordForms = im.ProcessOneWord(term, &lang)
		kp.Lang = lang
		if lang.IsRu() {
			totRusWords++
		}
		if lang.IsUa() {
			totUkrWords++
		}
		if lang.IsBy() {
			totByWords++
		}
		if lang.IsKz() {
			totKzWords++
		}
		if progress != nil {
			progress(cou, len(uniLex))
		}
		cou++
	}
	// Назначение WordForms каждому MorphToken
	var emptyList []*MorphForm.MorphWordForm
	for _, r := range res {
		uni, ok := r.Tag.(*MorphForm.UniLexWrap) // безопасный type assertion
		r.Tag = nil                              // обнуляем ссылку

		if !ok || uni == nil || uni.WordForms == nil || len(uni.WordForms) == 0 {
			if emptyList == nil {
				emptyList = []*MorphForm.MorphWordForm{}
			}
			r.WordForms = emptyList
			if ok && uni != nil {
				r.Language = uni.Lang
			}
		} else {
			r.WordForms = uni.WordForms
		}
	}

	// Установка CharInfo и языка для каждого MorphToken
	for i := 0; i < len(res); i++ {
		mt := res[i]
		mt.CharInfo = MorphForm.CharsInfo{}
		ui0 := twrch.At(mt.BeginChar)
		ui00 := morphinternal.GetChar(rune(mt.Term[0]))
		for j := mt.BeginChar + 1; j <= mt.EndChar; j++ {
			if ui0.IsLetter() {
				break
			}
			ui0 = twrch.At(j)
		}
		if ui0.IsLetter() {
			mt.CharInfo.SetLetter(true)
			if ui00.IsLatin() {
				mt.CharInfo.SetLatinLetter(true)
			} else if ui00.IsCyrillic() {
				mt.CharInfo.SetCyrillicLetter(true)
			}
			if mt.Language.IsUndefined() && MorphForm.IsCyrillic(mt.Term) {
				if defLang.IsUndefined() {
					mt.Language = MorphForm.MorphLangRU
				} else {
					mt.Language = defLang
				}
			}
			if goodText {
				continue
			}
			allUp := true
			allLo := true
			for j := mt.BeginChar; j <= mt.EndChar; j++ {
				if twrch.At(j).IsUpper() || twrch.At(j).IsDigit() {
					allLo = false
				} else {
					allUp = false
				}
			}
			if allUp {
				mt.CharInfo.SetAllUpper(true)
			} else if allLo {
				mt.CharInfo.SetAllLower(true)
			} else if (ui0.IsUpper() || twrch.At(mt.BeginChar).IsDigit()) && mt.EndChar > mt.BeginChar {
				allLo = true
				for j := mt.BeginChar + 1; j <= mt.EndChar; j++ {
					if twrch.At(j).IsUpper() || twrch.At(j).IsDigit() {
						allLo = false
						break
					}
				}
				if allLo {
					mt.CharInfo.SetCapitalUpper(true)
				} else if twrch.At(mt.EndChar).IsLower() && (mt.EndChar-mt.BeginChar) > 1 {
					allUp = true
					for j := mt.BeginChar; j < mt.EndChar; j++ {
						if twrch.At(j).IsLower() {
							allUp = false
							break
						}
					}
					if allUp {
						mt.CharInfo.SetLastLower(true)
					}
				}
			}
		}
		// Heuristic вставки недостающих форм при LastLower
		if mt.CharInfo.IsLastLower() && mt.Length() > 2 && mt.CharInfo.IsCyrillicLetter() {
			pref := text[mt.BeginChar : mt.EndChar+1]
			ok := false
			for _, wf := range mt.WordForms {
				if wf.NormalCase == pref || wf.NormalFull == pref {
					ok = true
					break
				}
			}
			if !ok {
				wf0 := &MorphForm.MorphWordForm{
					NormalCase: pref,
					Class:      MorphForm.MorphClassNoun,
					UndefCoef:  1,
				}
				mt.WordForms = append([]*MorphForm.MorphWordForm{wf0}, mt.WordForms...)
			}
		}
	}

	// Коррекция латинских букв A, C, P, если окружены кириллическими
	if !goodText && !onlyTokenizing {
		for i := 0; i < len(res); i++ {
			if res[i].Length() == 1 && res[i].CharInfo.IsLatinLetter() {
				ch := res[i].Term[0]
				if ch != 'C' && ch != 'A' && ch != 'P' {
					continue
				}
				isRus := false
				for ii := i - 1; ii >= 0; ii-- {
					if res[ii].EndChar+1 != res[ii+1].BeginChar {
						break
					} else if res[ii].CharInfo.IsLetter() {
						isRus = res[ii].CharInfo.IsCyrillicLetter()
						break
					}
				}
				if !isRus {
					for ii := i + 1; ii < len(res); ii++ {
						if res[ii-1].EndChar+1 != res[ii].BeginChar {
							break
						} else if res[ii].CharInfo.IsLetter() {
							isRus = res[ii].CharInfo.IsCyrillicLetter()
							break
						}
					}
				}
				if isRus {
					res[i].Term = MorphForm.TransliteralCorrection(res[i].Term, "", true)
					res[i].CharInfo.SetLatinLetter(true)
					res[i].CharInfo.SetCyrillicLetter(true)

				}
			}
		}
	}

	// Попытка добавить варианты фамилий к словам в верхнем регистре
	for _, r := range res {
		if (r.CharInfo.IsAllUpper() || r.CharInfo.IsCapitalUpper()) && r.Language.IsCyrillic() {
			ok := false
			for _, wf := range r.WordForms {
				if wf.Class.IsProperSurname() {
					ok = true
					break
				}
			}
			if !ok {
				r.WordForms = append([]*MorphForm.MorphWordForm{}, r.WordForms...)
				im.EngineRu.ProcessSurnameVariants(r.Term, &r.WordForms)
			}
		}
	}

	// Установка NormalCase, если отсутствует
	for _, r := range res {
		for _, wf := range r.WordForms {
			if wf.NormalCase == "" {
				wf.NormalCase = r.Term
			}
		}
	}

	for i := 0; i < len(res)-2; i++ {
		// Пример: A"Бакумов -> AБакумов
		if res[i].CharInfo.IsLatinLetter() && res[i].CharInfo.IsAllUpper() && res[i].Length() == 1 {
			ui1 := twrch.At(res[i+1].BeginChar)
			if ui1.IsQuot() && res[i+2].CharInfo.IsLatinLetter() && res[i+2].Length() > 2 {
				if res[i].EndChar+1 == res[i+1].BeginChar && res[i+1].EndChar+1 == res[i+2].BeginChar {
					wstr := res[i].Term + res[i+2].Term
					li := im.ProcessOneWord0(wstr)
					if li != nil {
						res[i].WordForms = li
					}
					res[i].EndChar = res[i+2].EndChar
					res[i].Term = wstr
					if res[i+2].CharInfo.IsAllLower() {
						res[i].CharInfo.SetAllUpper(false)
						res[i].CharInfo.SetCapitalUpper(true)
					} else if !res[i+2].CharInfo.IsAllUpper() {
						res[i].CharInfo.SetAllUpper(false)
					}
					res = append(res[:i+1], res[i+3:]...)
				}
			}
		}
	}

	// Объединение двойных дефисов типа "--"
	for i := 0; i < len(res)-1; i++ {
		if !res[i].CharInfo.IsLetter() && !res[i+1].CharInfo.IsLetter() && res[i].EndChar+1 == res[i+1].BeginChar {
			if twrch.At(res[i].BeginChar).IsHiphen() && twrch.At(res[i+1].BeginChar).IsHiphen() {
				if i == 0 || !twrch.At(res[i-1].BeginChar).IsHiphen() {
					if i+2 == len(res) || !twrch.At(res[i+2].BeginChar).IsHiphen() {
						res[i].EndChar = res[i+1].EndChar
						res = append(res[:i+1], res[i+2:]...)
					}
				}
			}
		}
	}

	return res // временно, для завершённого фрагмента
}

// GetAllWordforms возвращает все морфологические формы слова для заданного языка.
func (im *InnerMorphology) GetAllWordforms(word string, lang MorphForm.MorphLang) []*MorphForm.MorphWordForm {
	if MorphForm.IsCyrillicChar(rune(word[0])) {
		if !lang.IsUndefined() {
			if im.Engines["RU"].Language.IsRu() && lang.IsRu() {
				return im.Engines["RU"].GetAllWordforms(word)
			}
			if im.Engines["UA"].Language.IsUa() && lang.IsUa() {
				return im.Engines["UA"].GetAllWordforms(word)
			}
			if im.Engines["BY"].Language.IsBy() && lang.IsBy() {
				return im.Engines["BY"].GetAllWordforms(word)
			}
			if im.Engines["KZ"].Language.IsKz() && lang.IsKz() {
				return im.Engines["KZ"].GetAllWordforms(word)
			}
		}
		return im.Engines["RU"].GetAllWordforms(word)
	}
	return im.Engines["EN"].GetAllWordforms(word)
}

// GetAllWordsByClass возвращает все слова заданного класса (например, глаголы, существительные) для указанного языка.
func (im *InnerMorphology) GetAllWordsByClass(class MorphForm.MorphClass, lang MorphForm.MorphLang) []*MorphForm.MorphWordForm {
	var tmp strings.Builder
	var res []*MorphForm.MorphWordForm

	if im.Engines["RU"].Language.IsRu() && (lang.IsUndefined() || lang.IsRu()) {
		im.Engines["RU"].GetAllWordsByClass(class, &res, &tmp)
	}
	if im.Engines["EN"].Language.IsEn() && (lang.IsUndefined() || lang.IsEn()) {
		im.Engines["EN"].GetAllWordsByClass(class, &res, &tmp)
	}
	if im.Engines["UA"].Language.IsRu() && (lang.IsUndefined() || lang.IsUa()) {
		im.Engines["UA"].GetAllWordsByClass(class, &res, &tmp)
	}
	return res
}

// GetCharTyp возвращает тип символа: 1 — буква, 2 — цифра, 0 — пробел или код символа в остальных случаях.
func GetCharTyp(ui *morphinternal.UnicodeInfo) int {
	if ui.IsLetter() {
		return 1
	}
	if ui.IsDigit() {
		return 2
	}
	if ui.IsWhitespace() {
		return 0
	}
	if ui.IsUdaren() {
		return 1
	}
	return int(ui.Code)
}

// GetWordform возвращает слово в заданной морфологической форме (падеж, род, число и т.п.).
// Выбор морфологического движка осуществляется на основе языка слова и указанного языка.
func (im *InnerMorphology) GetWordform(
	word string,
	class MorphForm.MorphClass,
	gender MorphForm.MorphGender,
	cas MorphForm.MorphCase,
	num MorphForm.MorphNumber,
	lang MorphForm.MorphLang,
	addInfo *MorphForm.MorphWordForm,
) string {
	if MorphForm.IsCyrillicChar(rune(word[0])) {
		if im.Engines["RU"].Language.IsRu() && lang.IsRu() {
			return im.Engines["RU"].GetWordform(word, class, gender, cas, num, addInfo)
		}
		if im.Engines["UA"].Language.IsUa() && lang.IsUa() {
			return im.Engines["UA"].GetWordform(word, class, gender, cas, num, addInfo)
		}
		if im.Engines["BY"].Language.IsBy() && lang.IsBy() {
			return im.Engines["BY"].GetWordform(word, class, gender, cas, num, addInfo)
		}
		if im.Engines["KZ"].Language.IsKz() && lang.IsKz() {
			return im.Engines["KZ"].GetWordform(word, class, gender, cas, num, addInfo)
		}
		return im.Engines["RU"].GetWordform(word, class, gender, cas, num, addInfo)
	}
	return im.Engines["EN"].GetWordform(word, class, gender, cas, num, addInfo)
}

// CorrectWordByMorph пытается исправить слово на основе морфологического анализа,
// используя язык, если он задан. Если язык не задан, используется русский или английский по умолчанию.
func (im *InnerMorphology) CorrectWordByMorph(word string, lang MorphForm.MorphLang, oneVar bool) []string {
	if MorphForm.IsCyrillicChar(rune(word[0])) {
		if !lang.IsUndefined() {
			if im.Engines["RU"].Language.IsRu() && lang.IsRu() {
				return im.Engines["RU"].CorrectWordByMorph(word, oneVar)
			}
			if im.Engines["UA"].Language.IsUa() && lang.IsUa() {
				return im.Engines["UA"].CorrectWordByMorph(word, oneVar)
			}
			if im.Engines["BY"].Language.IsBy() && lang.IsBy() {
				return im.Engines["BY"].CorrectWordByMorph(word, oneVar)
			}
			if im.Engines["KZ"].Language.IsKz() && lang.IsKz() {
				return im.Engines["KZ"].CorrectWordByMorph(word, oneVar)
			}
		}
		return im.Engines["RU"].CorrectWordByMorph(word, oneVar)
	} else {
		return im.Engines["EN"].CorrectWordByMorph(word, oneVar)
	}
}

// ProcessOneWord0 выполняет морфологический разбор слова без явного указания языка.
// Язык может быть определён автоматически в процессе разбора.
func (im *InnerMorphology) ProcessOneWord0(wstr string) []*MorphForm.MorphWordForm {
	var dl MorphForm.MorphLang
	return im.ProcessOneWord(wstr, &dl)
}

// ProcessOneWord выполняет морфологический разбор слова с попыткой определить язык.
// Если язык не определён, defLang сбрасывается в неопределённое состояние и возвращается nil.
// Если язык точно установлен (RU, UA, BY, KZ, EN) — используется соответствующий движок.
// Если язык не определён, пробуются все доступные движки (RU, UA, BY, KZ), и результат
// выбирается в зависимости от наличия словарных форм и приоритетов.
// defLang может быть скорректирован в ходе выполнения.
func (im *InnerMorphology) ProcessOneWord(wstr string, defLang *MorphForm.MorphLang) []*MorphForm.MorphWordForm {
	lang := MorphForm.GetWordLang(wstr)
	if lang.IsUndefined() {
		*defLang = MorphForm.MorphLang{}
		return nil
	}

	if lang.Equals(MorphForm.MorphLangEN) {
		return im.Engines["EN"].Process(wstr, false)
	}

	if defLang.Equals(MorphForm.MorphLangRU) {
		if lang.IsRu() {
			return im.Engines["RU"].Process(wstr, false)
		}
	}

	if lang.Equals(MorphForm.MorphLangRU) {
		*defLang = *lang
		return im.Engines["RU"].Process(wstr, false)
	}

	if defLang.Equals(MorphForm.MorphLangUA) {
		if lang.IsUa() {
			return im.Engines["UA"].Process(wstr, false)
		}
	}

	if lang.Equals(MorphForm.MorphLangUA) {
		*defLang = *lang
		return im.Engines["UA"].Process(wstr, false)
	}

	if defLang.Equals(MorphForm.MorphLangBY) {
		if lang.IsBy() {
			return im.Engines["BY"].Process(wstr, false)
		}
	}

	if lang.Equals(MorphForm.MorphLangBY) {
		*defLang = *lang
		return im.Engines["BY"].Process(wstr, false)
	}

	if defLang.Equals(MorphForm.MorphLangKZ) {
		if lang.IsKz() {
			return im.Engines["KZ"].Process(wstr, false)
		}
	}

	if lang.Equals(MorphForm.MorphLangKZ) {
		*defLang = *lang
		return im.Engines["KZ"].Process(wstr, false)
	}

	var ru, ua, by, kz []*MorphForm.MorphWordForm

	if lang.IsRu() {
		ru = im.Engines["RU"].Process(wstr, false)
	}
	if lang.IsUa() {
		ua = im.Engines["UA"].Process(wstr, false)
	}
	if lang.IsBy() {
		by = im.Engines["BY"].Process(wstr, false)
	}
	if lang.IsKz() {
		kz = im.Engines["KZ"].Process(wstr, false)
	}

	hasRu := hasInDict(ru)
	hasUa := hasInDict(ua)
	hasBy := hasInDict(by)
	hasKz := hasInDict(kz)

	switch {
	case hasRu && !hasUa && !hasBy && !hasKz:
		*defLang = MorphForm.MorphLangRU
		return ru
	case hasUa && !hasRu && !hasBy && !hasKz:
		*defLang = MorphForm.MorphLangUA
		return ua
	case hasBy && !hasRu && !hasUa && !hasKz:
		*defLang = MorphForm.MorphLangBY
		return by
	case hasKz && !hasRu && !hasUa && !hasBy:
		*defLang = MorphForm.MorphLangKZ
		return kz
	}

	if ru == nil && ua == nil && by == nil && kz == nil {
		return nil
	}

	if ru != nil && ua == nil && by == nil && kz == nil {
		return ru
	}
	if ua != nil && ru == nil && by == nil && kz == nil {
		return ua
	}
	if by != nil && ru == nil && ua == nil && kz == nil {
		return by
	}
	if kz != nil && ru == nil && ua == nil && by == nil {
		return kz
	}

	var res []*MorphForm.MorphWordForm
	if ru != nil {
		langVal := *lang
		langVal = langVal.Or(MorphForm.MorphLangRU)
		*lang = langVal
		res = append(res, ru...)
	}
	if ua != nil {
		langVal := *lang
		langVal = langVal.Or(MorphForm.MorphLangUA)
		*lang = langVal
		res = append(res, ua...)
	}
	if by != nil {
		langVal := *lang
		langVal = langVal.Or(MorphForm.MorphLangBY)
		*lang = langVal
		res = append(res, by...)
	}
	if kz != nil {
		langVal := *lang
		langVal = langVal.Or(MorphForm.MorphLangKZ)
		*lang = langVal
		res = append(res, kz...)
	}
	return res
}

// hasInDict проверяет, содержит ли список морфформ хотя бы одну словарную форму.
func hasInDict(forms []*MorphForm.MorphWordForm) bool {
	for _, f := range forms {
		if f.IsInDictionary() {
			return true
		}
	}
	return false
}

// GetCharTyp определяет тип символа:
// 0 — пробел, 1 — буква (или ударение), 2 — цифра, иначе — Unicode-код символа.
func (m *InnerMorphology) GetCharTyp(ui *morphinternal.UnicodeInfo) int {
	if ui.IsLetter() {
		return 1
	}
	if ui.IsDigit() {
		return 2
	}
	if ui.IsWhitespace() {
		return 0
	}
	if ui.IsUdaren() {
		return 1
	}
	return int(ui.Code)
}
