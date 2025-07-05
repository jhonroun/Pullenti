package MorphForm

import (
	"strings"

	"github.com/jhonroun/pullenti/internal/morphinternal"
)

// MorphMiscInfo содержит дополнительную морфологическую информацию для слова (время, лицо, аспект и т.п.).
type MorphMiscInfo struct {
	Value int16
	Attrs []string
	Id    int

	IsAdjective         bool
	IsProperSurname     bool
	IsProperName        bool
	IsProperSecname     bool
	IsShortForm         bool
	IsPassiveParticiple bool
	IsActiveParticiple  bool
	IsTransitive        bool
	IsIntransitive      bool
	IsReflexive         bool
	IsPossessive        bool
}

// AddAttr добавляет новый атрибут, если его ещё нет.
func (m *MorphMiscInfo) AddAttr(a string) {
	for _, s := range m.Attrs {
		if s == a {
			return
		}
	}
	m.Attrs = append(m.Attrs, a)
}

// CopyFrom копирует данные из другого MorphMiscInfo.
func (m *MorphMiscInfo) CopyFrom(src *MorphMiscInfo) {
	m.Value = src.Value
	for _, a := range src.Attrs {
		m.AddAttr(a)
	}
}

// GetBoolValue возвращает булевый бит-флаг по индексу.
func (m *MorphMiscInfo) GetBoolValue(i int) bool {
	return (m.Value>>i)&1 != 0
}

// SetBoolValue устанавливает бит-флаг по индексу.
func (m *MorphMiscInfo) SetBoolValue(i int, val bool) {
	if val {
		m.Value |= 1 << i
	} else {
		m.Value &= ^(1 << i)
	}
}

// Person возвращает грамматическое лицо по атрибутам.
func (m *MorphMiscInfo) Person() MorphPerson {
	var res MorphPerson
	for _, a := range m.Attrs {
		switch a {
		case "1 л.":
			res |= MorphPersonFirst
		case "2 л.":
			res |= MorphPersonSecond
		case "3 л.":
			res |= MorphPersonThird
		}
	}
	return res
}

// SetPerson задаёт грамматическое лицо.
func (m *MorphMiscInfo) SetPerson(p MorphPerson) {
	if p.And(MorphPersonFirst).IsDefined() {
		m.AddAttr("1 л.")
	}
	if p.And(MorphPersonSecond).IsDefined() {
		m.AddAttr("2 л.")
	}
	if p.And(MorphPersonThird).IsDefined() {
		m.AddAttr("3 л.")
	}
}

// Tense возвращает грамматическое время по атрибутам.
func (m *MorphMiscInfo) Tense() MorphTense {
	for _, a := range m.Attrs {
		switch a {
		case "п.вр.":
			return MorphTensePast
		case "н.вр.":
			return MorphTensePresent
		case "б.вр.":
			return MorphTenseFuture
		}
	}
	return MorphTenseUndefined
}

// SetTense задаёт грамматическое время.
func (m *MorphMiscInfo) SetTense(t MorphTense) {
	if t == MorphTensePast {
		m.AddAttr("п.вр.")
	}
	if t == MorphTensePresent {
		m.AddAttr("н.вр.")
	}
	if t == MorphTenseFuture {
		m.AddAttr("б.вр.")
	}
}

// Aspect возвращает аспект (соверш./несоверш.).
func (m *MorphMiscInfo) Aspect() MorphAspect {
	for _, a := range m.Attrs {
		switch a {
		case "нес.в.":
			return MorphAspectImperfective
		case "сов.в.":
			return MorphAspectPerfective
		}
	}
	return MorphAspectUndefined
}

// SetAspect задаёт аспект.
func (m *MorphMiscInfo) SetAspect(a MorphAspect) {
	if a == MorphAspectImperfective {
		m.AddAttr("нес.в.")
	}
	if a == MorphAspectPerfective {
		m.AddAttr("сов.в.")
	}
}

// Mood возвращает наклонение (повелительное).
func (m *MorphMiscInfo) Mood() MorphMood {
	for _, a := range m.Attrs {
		if a == "пов.накл." {
			return MorphMoodImperative
		}
	}
	return MorphMoodUndefined
}

// SetMood задаёт наклонение.
func (m *MorphMiscInfo) SetMood(mood MorphMood) {
	if mood == MorphMoodImperative {
		m.AddAttr("пов.накл.")
	}
}

// Voice возвращает залог (действит./страдат.).
func (m *MorphMiscInfo) Voice() MorphVoice {
	for _, a := range m.Attrs {
		switch a {
		case "дейст.з.":
			return MorphVoiceActive
		case "страд.з.":
			return MorphVoicePassive
		}
	}
	return MorphVoiceUndefined
}

// SetVoice задаёт залог.
func (m *MorphMiscInfo) SetVoice(v MorphVoice) {
	if v == MorphVoiceActive {
		m.AddAttr("дейст.з.")
	}
	if v == MorphVoicePassive {
		m.AddAttr("страд.з.")
	}
}

// Form возвращает форму (краткая, синонимичная).
func (m *MorphMiscInfo) Form() MorphForm {
	for _, a := range m.Attrs {
		switch a {
		case "к.ф.":
			return MorphFormShort
		case "синоним.форма":
			return MorphFormSynonym
		}
	}
	if m.IsSynonymForm() {
		return MorphFormSynonym
	}
	return MorphFormUndefined
}

// IsSynonymForm возвращает true, если форма синонимичная.
func (m *MorphMiscInfo) IsSynonymForm() bool {
	return m.GetBoolValue(0)
}

// SetIsSynonymForm задаёт признак синонимичности формы.
func (m *MorphMiscInfo) SetIsSynonymForm(v bool) {
	m.SetBoolValue(0, v)
}

// String возвращает строковое представление всех морфологических атрибутов из списка Attrs.
func (m *MorphMiscInfo) String() string {
	if m == nil || (len(m.Attrs) == 0 && m.Value == 0) {
		return ""
	}
	var res strings.Builder
	if m.IsSynonymForm() {
		res.WriteString("синоним.форма ")
	}
	for _, attr := range m.Attrs {
		res.WriteString(attr)
		res.WriteByte(' ')
	}
	return strings.TrimSpace(res.String())
}

// ContainsAttr проверяет, содержится ли атрибут attr в списке Attrs.
func (m *MorphMiscInfo) ContainsAttr(attr string) bool {
	if m == nil {
		return false
	}
	for _, a := range m.Attrs {
		if a == attr {
			return true
		}
	}
	return false
}

// IsSubclassOf проверяет, входит ли текущий MorphClass в переданный (по флагам).
func (m MorphClass) IsSubclassOf(other MorphClass) bool {
	if m.Value == other.Value {
		return true
	}
	return (other.Value & m.Value) == m.Value
}

// Deserialize загружает объект MorphMiscInfo из сериализованного потока.
func (m *MorphMiscInfo) Deserialize(stream *morphinternal.ByteArrayWrapper, pos *int) {
	m.Value = int16(stream.DeserializeShort(pos))
	count := int(stream.DeserializeByte(pos))
	if count > 0 {
		m.Attrs = make([]string, count)
		for i := 0; i < count; i++ {
			m.Attrs[i] = stream.DeserializeString(pos)
		}
	}
}
