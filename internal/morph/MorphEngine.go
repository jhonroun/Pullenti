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
						if (cla.Value&v.Class.Value) != 0 && v.NormalTail != "" {
							if cas.IsUndefined() {
								if !v.Case.IsNominative() && !v.Case.IsUndefined() {
									continue
								}
							} else if v.Case.And(cas).IsUndefined() {
								continue
							}

							sur := cla.IsProperSurname()
							sur0 := v.Class.IsProperSurname()
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
		next, ok := tn.Nodes[ch]
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
		if (mv.Class.Value & cla.Value) != 0 {
			if mv.Tail != "" && !MorphForm.EndsWith(word, mv.Tail) {
				continue
			}
			wordBegin := word[:len(word)-len(mv.Tail)]
			for _, liv := range rule.MorphVars {
				for _, v := range liv {
					if (v.Class.Value & cla.Value) != 0 {
						sur := cla.IsProperSurname()
						sur0 := v.Class.IsProperSurname()
						if sur || sur0 {
							if sur != sur0 {
								continue
							}
						}
						if !cas.IsUndefined() {
							if v.Case.And(cas).IsUndefined() && !v.Case.IsUndefined() {
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
func (me *MorphEngine) _calcEqCoef(v *MorphForm.MorphRuleVariant, wf *MorphForm.MorphWordForm) int {
	if wf.Class.Value != 0 {
		if v.Class.Value&wf.Class.Value == 0 {
			return -1
		}
	}

	if wf.Misc != nil && v.MiscInfoId != int16(wf.Misc.Id) {
		vi := me.GetMiscInfo(int(v.MiscInfoId))
		if vi == nil {
			return -1
		}

		if vi.Mood.IsDefined() && wf.Misc.Mood.IsDefined() {
			if vi.Mood != wf.Misc.Mood {
				return -1
			}
		}

		if vi.Tense.IsDefined() && wf.Misc.Tense.IsDefined() {
			if (vi.Tense & wf.Misc.Tense).IsUndefined() {
				return -1
			}
		}

		if vi.Voice.IsDefined() && wf.Misc.Voice.IsDefined() {
			if vi.Voice != wf.Misc.Voice {
				return -1
			}
		}

		if vi.Person.IsDefined() && wf.Misc.Person.IsDefined() {
			if (vi.Person & wf.Misc.Person).IsUndefined() {
				return -1
			}
		}

		return 0
	}

	if !v.CheckAccord(wf, false, false) {
		return -1
	}

	return 1
}
