package MorphForm

import (
	"strings"
)

// MorphClass представляет собой часть речи (существ., глагол, прилагат. и т.д.).
type MorphClass struct {
	Value int16
}

// Битовые позиции для каждого типа.
const (
	mcNoun            = 0
	mcAdjective       = 1
	mcVerb            = 2
	mcAdverb          = 3
	mcPronoun         = 4
	mcMisc            = 5
	mcPreposition     = 6
	mcConjunction     = 7
	mcProper          = 8
	mcProperSurname   = 9
	mcProperName      = 10
	mcProperSecname   = 11
	mcProperGeo       = 12
	mcPersonalPronoun = 13
)

// get проверяет установлен ли флаг.
func (m MorphClass) get(i int) bool {
	return (m.Value>>i)&1 != 0
}

// set устанавливает или сбрасывает флаг.
func (m *MorphClass) set(i int, val bool) {
	if val {
		m.Value |= 1 << i
	} else {
		m.Value &= ^(1 << i)
	}
}

// Copy создает копию.
func (m MorphClass) Copy() MorphClass {
	return MorphClass{Value: m.Value}
}

// Equals сравнивает два MorphClass.
func (m MorphClass) Equals(other MorphClass) bool {
	return m.Value == other.Value
}

// Hash возвращает хеш-код.
func (m MorphClass) Hash() int {
	return int(m.Value)
}

// And — побитное И.
func (m MorphClass) And(other MorphClass) MorphClass {
	return MorphClass{Value: m.Value & other.Value}
}

// Or — побитное ИЛИ.
func (m MorphClass) Or(other MorphClass) MorphClass {
	return MorphClass{Value: m.Value | other.Value}
}

// Xor — побитное XOR.
func (m MorphClass) Xor(other MorphClass) MorphClass {
	return MorphClass{Value: m.Value ^ other.Value}
}

// Статические эквиваленты
var (
	MorphClassUndefined       = MorphClass{Value: 0}
	MorphClassNoun            = MorphClassFrom(mcNoun)
	MorphClassPronoun         = MorphClassFrom(mcPronoun)
	MorphClassPersonalPronoun = MorphClassFrom(mcPersonalPronoun)
	MorphClassVerb            = MorphClassFrom(mcVerb)
	MorphClassAdjective       = MorphClassFrom(mcAdjective)
	MorphClassAdverb          = MorphClassFrom(mcAdverb)
	MorphClassPreposition     = MorphClassFrom(mcPreposition)
	MorphClassConjunction     = MorphClassFrom(mcConjunction)
)

func MorphClassFrom(bit int) MorphClass {
	return MorphClass{Value: 1 << bit}
}

// Свойства

func (m MorphClass) IsUndefined() bool { return m.Value == 0 }
func (m *MorphClass) SetUndefined()    { m.Value = 0 }
func (m MorphClass) IsNoun() bool      { return m.get(mcNoun) }
func (m *MorphClass) SetNoun(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcNoun, v)
}
func (m MorphClass) IsAdjective() bool { return m.get(mcAdjective) }
func (m *MorphClass) SetAdjective(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcAdjective, v)
}
func (m MorphClass) IsVerb() bool { return m.get(mcVerb) }
func (m *MorphClass) SetVerb(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcVerb, v)
}
func (m MorphClass) IsAdverb() bool { return m.get(mcAdverb) }
func (m *MorphClass) SetAdverb(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcAdverb, v)
}
func (m MorphClass) IsPronoun() bool { return m.get(mcPronoun) }
func (m *MorphClass) SetPronoun(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcPronoun, v)
}
func (m MorphClass) IsMisc() bool { return m.get(mcMisc) }
func (m *MorphClass) SetMisc(v bool) {
	if v {
		m.Value = 0
	}
	m.set(mcMisc, v)
}
func (m MorphClass) IsPreposition() bool    { return m.get(mcPreposition) }
func (m *MorphClass) SetPreposition(v bool) { m.set(mcPreposition, v) }
func (m MorphClass) IsConjunction() bool    { return m.get(mcConjunction) }
func (m *MorphClass) SetConjunction(v bool) { m.set(mcConjunction, v) }
func (m MorphClass) IsProper() bool         { return m.get(mcProper) }
func (m *MorphClass) SetProper(v bool)      { m.set(mcProper, v) }
func (m MorphClass) IsProperSurname() bool  { return m.get(mcProperSurname) }
func (m *MorphClass) SetProperSurname(v bool) {
	if v {
		m.set(mcProper, true)
	}
	m.set(mcProperSurname, v)
}
func (m MorphClass) IsProperName() bool { return m.get(mcProperName) }
func (m *MorphClass) SetProperName(v bool) {
	if v {
		m.set(mcProper, true)
	}
	m.set(mcProperName, v)
}
func (m MorphClass) IsProperSecname() bool { return m.get(mcProperSecname) }
func (m *MorphClass) SetProperSecname(v bool) {
	if v {
		m.set(mcProper, true)
	}
	m.set(mcProperSecname, v)
}
func (m MorphClass) IsProperGeo() bool { return m.get(mcProperGeo) }
func (m *MorphClass) SetProperGeo(v bool) {
	if v {
		m.set(mcProper, true)
	}
	m.set(mcProperGeo, v)
}
func (m MorphClass) IsPersonalPronoun() bool { return m.get(mcPersonalPronoun) }
func (m *MorphClass) SetPersonalPronoun(v bool) {
	m.set(mcPersonalPronoun, v)
}

// String возвращает строковое представление.
func (m MorphClass) String() string {
	var sb strings.Builder
	if m.IsNoun() {
		sb.WriteString("существ.|")
	}
	if m.IsAdjective() {
		sb.WriteString("прилаг.|")
	}
	if m.IsVerb() {
		sb.WriteString("глагол|")
	}
	if m.IsAdverb() {
		sb.WriteString("наречие|")
	}
	if m.IsPronoun() {
		sb.WriteString("местоим.|")
	}
	if m.IsMisc() && !m.IsConjunction() && !m.IsPreposition() && !m.IsProper() {
		sb.WriteString("разное|")
	}
	if m.IsPreposition() {
		sb.WriteString("предлог|")
	}
	if m.IsConjunction() {
		sb.WriteString("союз|")
	}
	if m.IsProper() {
		sb.WriteString("собств.|")
	}
	if m.IsProperSurname() {
		sb.WriteString("фамилия|")
	}
	if m.IsProperName() {
		sb.WriteString("имя|")
	}
	if m.IsProperSecname() {
		sb.WriteString("отч.|")
	}
	if m.IsProperGeo() {
		sb.WriteString("геогр.|")
	}
	if m.IsPersonalPronoun() {
		sb.WriteString("личн.местоим.|")
	}
	res := sb.String()
	if len(res) > 0 {
		res = res[:len(res)-1] // remove last '|'
	}
	return res
}

// ClearMisc сбрасывает флаги подтипов: фамилия, имя, отчество, геогр., личн. местоим.
func (m *MorphClass) ClearMisc() {
	m.Value &= ^(MorphClassFrom(mcProperSurname).Value |
		MorphClassFrom(mcProperName).Value |
		MorphClassFrom(mcProperSecname).Value |
		MorphClassFrom(mcProperGeo).Value |
		MorphClassFrom(mcPersonalPronoun).Value)
}
