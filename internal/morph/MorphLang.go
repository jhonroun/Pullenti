package morph

import (
	"strings"
	"unicode"
)

// Язык
type MorphLang struct {
	Value int16
}

var mNames = []string{"RU", "UA", "BY", "EN", "IT", "KZ"}

func (m *MorphLang) getValue(i int) bool {
	return (m.Value>>i)&1 != 0
}

func (m *MorphLang) setValue(i int, val bool) {
	if val {
		m.Value |= 1 << i
	} else {
		m.Value &= ^(1 << i)
	}
}

// Неопределённый язык
func (m *MorphLang) IsUndefined() bool {
	return m.Value == 0
}

func (m *MorphLang) SetUndefined(val bool) {
	if val {
		m.Value = 0
	}
}

// Русский язык
func (m *MorphLang) IsRu() bool {
	return m.getValue(0)
}
func (m *MorphLang) SetRu(val bool) {
	m.setValue(0, val)
}

// Украинский язык
func (m *MorphLang) IsUa() bool {
	return m.getValue(1)
}
func (m *MorphLang) SetUa(val bool) {
	m.setValue(1, val)
}

// Белорусский язык
func (m *MorphLang) IsBy() bool {
	return m.getValue(2)
}
func (m *MorphLang) SetBy(val bool) {
	m.setValue(2, val)
}

// Английский язык
func (m *MorphLang) IsEn() bool {
	return m.getValue(3)
}
func (m *MorphLang) SetEn(val bool) {
	m.setValue(3, val)
}

// Итальянский язык
func (m *MorphLang) IsIt() bool {
	return m.getValue(4)
}
func (m *MorphLang) SetIt(val bool) {
	m.setValue(4, val)
}

// Казахский язык
func (m *MorphLang) IsKz() bool {
	return m.getValue(5)
}
func (m *MorphLang) SetKz(val bool) {
	m.setValue(5, val)
}

// Русский, украинский, белорусский или казахский язык
func (m *MorphLang) IsCyrillic() bool {
	return m.IsRu() || m.IsUa() || m.IsBy() || m.IsKz()
}

// Преобразование в строку
func (m *MorphLang) String() string {
	var b strings.Builder
	if m.IsRu() {
		b.WriteString("RU;")
	}
	if m.IsUa() {
		b.WriteString("UA;")
	}
	if m.IsBy() {
		b.WriteString("BY;")
	}
	if m.IsEn() {
		b.WriteString("EN;")
	}
	if m.IsIt() {
		b.WriteString("IT;")
	}
	if m.IsKz() {
		b.WriteString("KZ;")
	}
	s := b.String()
	if len(s) > 0 && s[len(s)-1] == ';' {
		s = s[:len(s)-1]
	}
	return s
}

// Сравнение значение (полное совпадение)
func (m *MorphLang) Equals(other *MorphLang) bool {
	if other == nil {
		return false
	}
	return m.Value == other.Value
}

func (m *MorphLang) HashCode() int {
	return int(m.Value)
}

// Преобразовать из строки
func TryParseMorphLang(str string) (*MorphLang, bool) {
	lang := &MorphLang{}
	for len(str) > 0 {
		i := -1
		for idx, name := range mNames {
			if strings.HasPrefix(strings.ToUpper(str), name) {
				i = idx
				break
			}
		}
		if i < 0 {
			break
		}
		lang.Value |= 1 << i

		// Сдвигаем строку после текущего кода языка
		str = str[len(mNames[i]):]
		for j, r := range str {
			if unicode.IsLetter(r) {
				str = str[j:]
				break
			}
			if j == len(str)-1 {
				str = ""
			}
		}
	}
	if lang.IsUndefined() {
		return lang, false
	}
	return lang, true
}

// Моделирование побитного "AND"
func AndMorphLang(arg1, arg2 *MorphLang) *MorphLang {
	var val1, val2 int16
	if arg1 != nil {
		val1 = arg1.Value
	}
	if arg2 != nil {
		val2 = arg2.Value
	}
	return &MorphLang{Value: val1 & val2}
}

// Моделирование побитного "OR"
func OrMorphLang(arg1, arg2 *MorphLang) *MorphLang {
	var val1, val2 int16
	if arg1 != nil {
		val1 = arg1.Value
	}
	if arg2 != nil {
		val2 = arg2.Value
	}
	return &MorphLang{Value: val1 | val2}
}

// Неопределённое
var Unknown = &MorphLang{}

// Русский
var RU = &MorphLang{Value: 1 << 0}

// Украинский
var UA = &MorphLang{Value: 1 << 1}

// Белорусский
var BY = &MorphLang{Value: 1 << 2}

// Английский
var EN = &MorphLang{Value: 1 << 3}

// Итальянский
var IT = &MorphLang{Value: 1 << 4}

// Казахский
var KZ = &MorphLang{Value: 1 << 5}
