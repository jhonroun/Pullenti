package MorphForm

import (
	"strings"
)

// MorphWordForm представляет собой одну словоформу с грамматической информацией.
type MorphWordForm struct {
	MorphBaseInfo
	Tail           string
	NormalTail     string
	FullNormalTail string
	Misc           *MorphMiscInfo
	RuleId         int16
	Id             int16
	Class          MorphClass
	Case           MorphCase
	Gender         MorphGender
	Number         MorphNumber
	Person         MorphPerson
	Voice          MorphVoice
	Tense          MorphTense
	Aspect         MorphAspect
	Mood           MorphMood
	UndefCoef      int16
	NormalCase     string
	NormalFull     string
	Tag            any
}

// CopyFromWordForm копирует все поля из другой словоформы.
func (m *MorphWordForm) CopyFromWordForm(src *MorphWordForm) {
	m.Tail = src.Tail
	m.NormalTail = src.NormalTail
	m.FullNormalTail = src.FullNormalTail
	m.Misc = src.Misc
	m.RuleId = src.RuleId
	m.Id = src.Id
	m.Class = src.Class
	m.Case = src.Case
	m.Gender = src.Gender
	m.Number = src.Number
	m.Person = src.Person
	m.Voice = src.Voice
	m.Tense = src.Tense
	m.Aspect = src.Aspect
	m.Mood = src.Mood
	m.UndefCoef = src.UndefCoef
	m.NormalCase = src.NormalCase
	m.NormalFull = src.NormalFull
	m.Tag = src.Tag
}

// String возвращает строковое представление словоформы.
func (m *MorphWordForm) String() string {
	return m.ToStringEx(false)
}

// ToStringEx возвращает строку с деталями формы. Если ignoreNormals=true, не включаются нормальные формы.
func (m *MorphWordForm) ToStringEx(ignoreNormals bool) string {
	var res strings.Builder
	if !ignoreNormals {
		if m.NormalCase != "" {
			res.WriteString(m.NormalCase)
		} else if m.NormalFull != "" {
			res.WriteString(m.NormalFull)
		}
	}
	if !m.Class.IsUndefined() {
		res.WriteString(" ")
		res.WriteString(m.Class.String())
	}
	if !m.Case.IsUndefined() {
		res.WriteString(" ")
		res.WriteString(m.Case.String())
	}
	if m.Number != 0 {
		res.WriteString(" ")
		res.WriteString(m.Number.String())
	}
	if m.Gender != 0 {
		res.WriteString(" ")
		res.WriteString(m.Gender.String())
	}
	if m.Person != 0 {
		res.WriteString(" ")
		res.WriteString(m.Person.String())
	}
	if m.Tense != 0 {
		res.WriteString(" ")
		res.WriteString(m.Tense.String())
	}
	if m.Voice != 0 {
		res.WriteString(" ")
		res.WriteString(m.Voice.String())
	}
	if m.Aspect != 0 {
		res.WriteString(" ")
		res.WriteString(m.Aspect.String())
	}
	if m.Mood != 0 {
		res.WriteString(" ")
		res.WriteString(m.Mood.String())
	}
	if m.Misc != nil {
		res.WriteString(" ")
		res.WriteString(m.Misc.String())
	}
	return strings.TrimSpace(res.String())
}

// IsInDictionary возвращает true, если словоформа получена из словаря.
func (m *MorphWordForm) IsInDictionary() bool {
	return m.UndefCoef == 0
}

// HasMorphEquals проверяет морфологическое совпадение с одним из вариантов из списка.
func (m *MorphWordForm) HasMorphEquals(list []*MorphWordForm) bool {
	for _, wf := range list {
		if wf.Class.Equals(m.Class) &&
			wf.Number == m.Number &&
			wf.Case.Equals(m.Case) &&
			wf.Gender == m.Gender &&
			wf.Person == m.Person {
			return true
		}
	}
	return false
}

// ContainsAttr проверяет наличие морфологического признака в Misc и совпадение класса.
func (m *MorphWordForm) ContainsAttr(attr string, class *MorphClass) bool {
	if m.Misc != nil && m.Misc.ContainsAttr(attr) {
		if class == nil || class.IsUndefined() || m.Class.Equals(*class) || m.Class.IsSubclassOf(*class) {
			return true
		}
	}
	return false
}

// NewMorphWordFormFromVariant создаёт MorphWordForm на основе MorphRuleVariant и слова.
// Аналог конструктора MorphWordForm(C#).
func NewMorphWordFormFromVariant(v *MorphRuleVariant, word string, mi *MorphMiscInfo) *MorphWordForm {
	if v == nil {
		return nil
	}

	wf := &MorphWordForm{}
	wf.CopyFromVariant(v) // Копируем основные поля из варианта
	wf.Misc = mi

	// Построение NormalCase
	if v.NormalTail != "" && word != "" {
		wordBegin := word
		if LanguageHelperEndsWith(word, v.Tail) && len(v.Tail) <= len(word) {
			wordBegin = word[:len(word)-len(v.Tail)]
		}
		if v.NormalTail != "" {
			wf.NormalCase = wordBegin + v.NormalTail
		} else {
			wf.NormalCase = wordBegin
		}
	}

	// Построение NormalFull
	if v.FullNormalTail != "" && word != "" {
		wordBegin := word
		if LanguageHelperEndsWith(word, v.Tail) && len(v.Tail) <= len(word) {
			wordBegin = word[:len(word)-len(v.Tail)]
		}
		if v.FullNormalTail != "" {
			wf.NormalFull = wordBegin + v.FullNormalTail
		} else {
			wf.NormalFull = wordBegin
		}
	}

	return wf
}

// CopyFromVariant копирует все морфологические данные из MorphRuleVariant в MorphWordForm.
func (wf *MorphWordForm) CopyFromVariant(v *MorphRuleVariant) {
	if v == nil {
		return
	}
	wf.Tail = v.Tail
	wf.NormalTail = v.NormalTail
	wf.FullNormalTail = v.FullNormalTail
	wf.Class = v.GetClass() //v.Class
	wf.Gender = v.Gender
	wf.Number = v.Number
	wf.Case = v.GetCase() //v.Case
	wf.Misc = nil         // будет установлен позже
	wf.Id = v.MiscInfoId
	wf.RuleId = v.RuleId
	wf.UndefCoef = 0 // используется позже при сортировке
}

// LanguageHelperEndsWith проверяет, заканчивается ли строка word на suffix (без учёта регистра).
func LanguageHelperEndsWith(word, suffix string) bool {
	if len(suffix) > len(word) {
		return false
	}
	return word[len(word)-len(suffix):] == suffix
}
