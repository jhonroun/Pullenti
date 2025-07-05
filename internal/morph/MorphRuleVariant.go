package morph

import (
	"fmt"
	"strings"
)

// MorphRuleVariant соответствует MorphRuleVariant из C#
type MorphRuleVariant struct {
	MorphBaseInfo

	Tail           string
	MiscInfoId     int16
	RuleId         int16
	Id             int16
	NormalTail     string
	FullNormalTail string
	Tag            any
}

// CopyFromVariant копирует значения из другого варианта
func (m *MorphRuleVariant) CopyFromVariant(src *MorphRuleVariant) {
	if src == nil {
		return
	}
	m.Tail = src.Tail
	m.MorphBaseInfo.CopyFrom(&src.MorphBaseInfo)
	m.MiscInfoId = src.MiscInfoId
	m.NormalTail = src.NormalTail
	m.FullNormalTail = src.FullNormalTail
	m.RuleId = src.RuleId
}

// ToString возвращает строковое представление без скрытия хвостов
func (m *MorphRuleVariant) String() string {
	return m.ToStringEx(false)
}

// ToStringEx с возможностью скрытия хвостов
func (m *MorphRuleVariant) ToStringEx(hideTails bool) string {
	var res strings.Builder
	if !hideTails {
		res.WriteString(fmt.Sprintf("-%s", m.Tail))
		if m.NormalTail != "" {
			res.WriteString(fmt.Sprintf(" [-%s]", m.NormalTail))
		}
		if m.FullNormalTail != "" && m.FullNormalTail != m.NormalTail {
			res.WriteString(fmt.Sprintf(" [-%s]", m.FullNormalTail))
		}
	}
	res.WriteString(" " + m.MorphBaseInfo.String())
	return strings.TrimSpace(res.String())
}

// Compare проверяет совпадение по морфологическим признакам и хвостам
func (m *MorphRuleVariant) Compare(other *MorphRuleVariant) bool {
	if !other.Class.Equal(m.Class) ||
		other.Gender != m.Gender ||
		other.Number != m.Number ||
		!other.Case.Equal(m.Case) {
		return false
	}
	if other.MiscInfoId != m.MiscInfoId {
		return false
	}
	if other.NormalTail != m.NormalTail {
		return false
	}
	return true
}

// Deserialize читает поля из ByteArrayWrapper
func (m *MorphRuleVariant) Deserialize(str *ByteArrayWrapper, pos *int) bool {
	id := str.DeserializeShort(pos)
	if id <= 0 {
		return false
	}
	m.MiscInfoId = int16(id)

	iii := str.DeserializeShort(pos)
	mc := morph.MorphClass{}
	mc.Value = int16(iii)
	if mc.IsMisc() && mc.IsProper() {
		mc.ClearMisc()
	}
	m.Class = mc

	bbb := str.DeserializeByte(pos)
	m.Gender = morph.MorphGender(bbb)

	bbb = str.DeserializeByte(pos)
	m.Number = morph.MorphNumber(bbb)

	bbb = str.DeserializeByte(pos)
	mca := morph.MorphCase{}
	mca.Value = int16(bbb)
	m.Case = mca

	m.NormalTail = str.DeserializeString(pos)
	m.FullNormalTail = str.DeserializeString(pos)
	return true
}
