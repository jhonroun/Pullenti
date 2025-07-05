package morph

import (
	"strings"
)

// Падеж
type MorphCase struct {
	Value int16
}

var caseNames = []string{"именит.", "родит.", "дател.", "винит.", "творит.", "предлож.", "зват.", "частич.", "общ.", "притяж."}

func (m *MorphCase) getValue(i int) bool {
	return (m.Value>>i)&1 != 0
}

func (m *MorphCase) setValue(i int, val bool) {
	if val {
		m.Value |= 1 << i
	} else {
		m.Value &= ^(1 << i)
	}
}

// Неопределённый падеж
func (m *MorphCase) IsUndefined() bool {
	return m.Value == 0
}

func (m *MorphCase) SetUndefined(val bool) {
	if val {
		m.Value = 0
	}
}

// Количество падежей
func (m *MorphCase) Count() int {
	if m.Value == 0 {
		return 0
	}
	count := 0
	for i := 0; i < 12; i++ {
		if (m.Value & (1 << i)) != 0 {
			count++
		}
	}
	return count
}

// Именительный
func (m *MorphCase) IsNominative() bool   { return m.getValue(0) }
func (m *MorphCase) SetNominative(v bool) { m.setValue(0, v) }

// Родительный
func (m *MorphCase) IsGenitive() bool   { return m.getValue(1) }
func (m *MorphCase) SetGenitive(v bool) { m.setValue(1, v) }

// Дательный
func (m *MorphCase) IsDative() bool   { return m.getValue(2) }
func (m *MorphCase) SetDative(v bool) { m.setValue(2, v) }

// Винительный
func (m *MorphCase) IsAccusative() bool   { return m.getValue(3) }
func (m *MorphCase) SetAccusative(v bool) { m.setValue(3, v) }

// Творительный
func (m *MorphCase) IsInstrumental() bool   { return m.getValue(4) }
func (m *MorphCase) SetInstrumental(v bool) { m.setValue(4, v) }

// Предложный
func (m *MorphCase) IsPrepositional() bool   { return m.getValue(5) }
func (m *MorphCase) SetPrepositional(v bool) { m.setValue(5, v) }

// Звательный
func (m *MorphCase) IsVocative() bool   { return m.getValue(6) }
func (m *MorphCase) SetVocative(v bool) { m.setValue(6, v) }

// Частичный
func (m *MorphCase) IsPartial() bool   { return m.getValue(7) }
func (m *MorphCase) SetPartial(v bool) { m.setValue(7, v) }

// Общий (английский)
func (m *MorphCase) IsCommon() bool   { return m.getValue(8) }
func (m *MorphCase) SetCommon(v bool) { m.setValue(8, v) }

// Притяжательный (английский)
func (m *MorphCase) IsPossessive() bool   { return m.getValue(9) }
func (m *MorphCase) SetPossessive(v bool) { m.setValue(9, v) }

// Строковое представление
func (m *MorphCase) String() string {
	var b strings.Builder
	if m.IsNominative() {
		b.WriteString("именит.|")
	}
	if m.IsGenitive() {
		b.WriteString("родит.|")
	}
	if m.IsDative() {
		b.WriteString("дател.|")
	}
	if m.IsAccusative() {
		b.WriteString("винит.|")
	}
	if m.IsInstrumental() {
		b.WriteString("творит.|")
	}
	if m.IsPrepositional() {
		b.WriteString("предлож.|")
	}
	if m.IsVocative() {
		b.WriteString("зват.|")
	}
	if m.IsPartial() {
		b.WriteString("частич.|")
	}
	if m.IsCommon() {
		b.WriteString("общ.|")
	}
	if m.IsPossessive() {
		b.WriteString("притяж.|")
	}
	s := b.String()
	if len(s) > 0 {
		s = s[:len(s)-1]
	}
	return s
}

// Восстановить падежи из строки, полученной ToString
func ParseMorphCase(str string) *MorphCase {
	res := &MorphCase{}
	if str == "" {
		return res
	}
	parts := strings.Split(str, "|")
	for _, part := range parts {
		for i, name := range caseNames {
			if part == name {
				res.setValue(i, true)
				break
			}
		}
	}
	return res
}

// Проверка на полное совпадение значений
func (m *MorphCase) Equals(other *MorphCase) bool {
	if other == nil {
		return false
	}
	return m.Value == other.Value
}

func (m *MorphCase) HashCode() int {
	return int(m.Value)
}

// Побитный AND
func AndMorphCase(a, b *MorphCase) *MorphCase {
	var v1, v2 int16
	if a != nil {
		v1 = a.Value
	}
	if b != nil {
		v2 = b.Value
	}
	return &MorphCase{Value: v1 & v2}
}

// Побитный OR
func OrMorphCase(a, b *MorphCase) *MorphCase {
	var v1, v2 int16
	if a != nil {
		v1 = a.Value
	}
	if b != nil {
		v2 = b.Value
	}
	return &MorphCase{Value: v1 | v2}
}

// Побитный XOR
func XorMorphCase(a, b *MorphCase) *MorphCase {
	var v1, v2 int16
	if a != nil {
		v1 = a.Value
	}
	if b != nil {
		v2 = b.Value
	}
	return &MorphCase{Value: v1 ^ v2}
}

// Статические экземпляры
var (
	MorphCaseUndefined     = &MorphCase{Value: 0}
	MorphCaseNominative    = &MorphCase{Value: 1}
	MorphCaseGenitive      = &MorphCase{Value: 2}
	MorphCaseDative        = &MorphCase{Value: 4}
	MorphCaseAccusative    = &MorphCase{Value: 8}
	MorphCaseInstrumental  = &MorphCase{Value: 0x10}
	MorphCasePrepositional = &MorphCase{Value: 0x20}
	MorphCaseVocative      = &MorphCase{Value: 0x40}
	MorphCasePartial       = &MorphCase{Value: 0x80}
	MorphCaseCommon        = &MorphCase{Value: 0x100}
	MorphCasePossessive    = &MorphCase{Value: 0x200}
	MorphCaseAllCases      = &MorphCase{Value: 0x3FF}
)
