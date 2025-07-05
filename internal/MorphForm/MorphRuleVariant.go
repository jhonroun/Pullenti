package MorphForm

import (
	"strings"

	"github.com/jhonroun/pullenti/internal/morphinternal"
)

// MorphRuleVariant описывает морфологический вариант, производный от MorphBaseInfo.
type MorphRuleVariant struct {
	MorphBaseInfo
	Tail           string      // Конечная часть (хвост)
	MiscInfoId     int16       // Идентификатор дополнительной информации
	RuleId         int16       // Идентификатор правила
	Id             int16       // Внутренний ID варианта
	NormalTail     string      // Нормализованный хвост
	FullNormalTail string      // Полный нормализованный хвост
	Tag            interface{} // Пользовательский тег (аналог object Tag в C#)
}

// CopyFromVariant копирует содержимое другого варианта.
func (v *MorphRuleVariant) CopyFromVariant(src *MorphRuleVariant) {
	if src == nil {
		return
	}
	v.Tail = src.Tail
	v.CopyFrom(&src.MorphBaseInfo)
	v.MiscInfoId = src.MiscInfoId
	v.NormalTail = src.NormalTail
	v.FullNormalTail = src.FullNormalTail
	v.RuleId = src.RuleId
}

// String возвращает строковое представление, включая хвосты.
func (v *MorphRuleVariant) String() string {
	return v.ToStringEx(false)
}

// ToStringEx возвращает строку с опцией скрыть хвосты.
func (v *MorphRuleVariant) ToStringEx(hideTails bool) string {
	var res strings.Builder
	if !hideTails {
		res.WriteString("-")
		res.WriteString(v.Tail)
		if v.NormalTail != "" {
			res.WriteString(" [-")
			res.WriteString(v.NormalTail)
			res.WriteString("]")
		}
		if v.FullNormalTail != "" && v.FullNormalTail != v.NormalTail {
			res.WriteString(" [-")
			res.WriteString(v.FullNormalTail)
			res.WriteString("]")
		}
	}
	res.WriteString(" ")
	res.WriteString(v.MorphBaseInfo.String())
	return strings.TrimSpace(res.String())
}

// Compare сравнивает текущий вариант с другим по ключевым полям.
func (v *MorphRuleVariant) Compare(other *MorphRuleVariant) bool {
	if !v.GetClass().Equals(other.GetClass()) ||
		v.Gender != other.Gender ||
		v.Number != other.Number ||
		!v.GetCase().Equals(other.GetCase()) ||
		v.MiscInfoId != other.MiscInfoId ||
		v.NormalTail != other.NormalTail {
		return false
	}
	return true
}

// Deserialize загружает поля из сериализованного байтового потока.
func (v *MorphRuleVariant) Deserialize(str *morphinternal.ByteArrayWrapper, pos *int) bool {
	id := str.DeserializeShort(pos)
	if id <= 0 {
		return false
	}
	v.MiscInfoId = int16(id)

	// Класс
	iii := str.DeserializeShort(pos)
	class := MorphClass{Value: int16(iii)}
	if class.IsMisc() && class.IsProper() {
		class.SetMisc(false)
	}
	v.SetClass(class)

	// Род
	bbb := str.DeserializeByte(pos)
	v.Gender = MorphGender(bbb)

	// Число
	bbb = str.DeserializeByte(pos)
	v.Number = MorphNumber(bbb)

	// Падеж
	bbb = str.DeserializeByte(pos)
	v.SetCase(MorphCase{Value: int16(bbb)})

	// Хвосты
	v.NormalTail = str.DeserializeString(pos)
	v.FullNormalTail = str.DeserializeString(pos)

	return true
}
