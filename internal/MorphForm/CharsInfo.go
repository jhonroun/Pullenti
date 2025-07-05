package MorphForm

import (
	"strings"
	"unicode"
)

// Информация о символах токена
type CharsInfo struct {
	Value int16
}

func (c *CharsInfo) getValue(i int) bool {
	return ((c.Value >> i) & 1) != 0
}

func (c *CharsInfo) setValue(i int, val bool) {
	if val {
		c.Value |= 1 << i
	} else {
		c.Value &= ^(1 << i)
	}
}

// Все символы в верхнем регистре
func (c *CharsInfo) IsAllUpper() bool {
	return c.getValue(0)
}

func (c *CharsInfo) SetAllUpper(val bool) {
	c.setValue(0, val)
}

// Все символы в нижнем регистре
func (c *CharsInfo) IsAllLower() bool {
	return c.getValue(1)
}

func (c *CharsInfo) SetAllLower(val bool) {
	c.setValue(1, val)
}

// Первый символ в верхнем регистре, остальные в нижнем.
// Для однобуквенной комбинации false.
func (c *CharsInfo) IsCapitalUpper() bool {
	return c.getValue(2)
}

func (c *CharsInfo) SetCapitalUpper(val bool) {
	c.setValue(2, val)
}

// Все символы в верхнем регистре, кроме последнего (длина >= 3)
func (c *CharsInfo) IsLastLower() bool {
	return c.getValue(3)
}

func (c *CharsInfo) SetLastLower(val bool) {
	c.setValue(3, val)
}

// Это буквы
func (c *CharsInfo) IsLetter() bool {
	return c.getValue(4)
}

func (c *CharsInfo) SetLetter(val bool) {
	c.setValue(4, val)
}

// Это латиница
func (c *CharsInfo) IsLatinLetter() bool {
	return c.getValue(5)
}

func (c *CharsInfo) SetLatinLetter(val bool) {
	c.setValue(5, val)
}

// Это кириллица
func (c *CharsInfo) IsCyrillicLetter() bool {
	return c.getValue(6)
}

func (c *CharsInfo) SetCyrillicLetter(val bool) {
	c.setValue(6, val)
}

func (c *CharsInfo) String() string {
	if !c.IsLetter() {
		return "Nonletter"
	}
	var tmp strings.Builder
	if c.IsAllUpper() {
		tmp.WriteString("AllUpper")
	} else if c.IsAllLower() {
		tmp.WriteString("AllLower")
	} else if c.IsCapitalUpper() {
		tmp.WriteString("CapitalUpper")
	} else if c.IsLastLower() {
		tmp.WriteString("LastLower")
	} else {
		tmp.WriteString("Nonstandard")
	}
	if c.IsLatinLetter() {
		tmp.WriteString(" Latin")
	} else if c.IsCyrillicLetter() {
		tmp.WriteString(" Cyrillic")
	} else if c.IsLetter() {
		tmp.WriteString(" Letter")
	}
	return tmp.String()
}

// Сравнение на совпадение значений всех полей
// obj - сравниваемый объект
func (c *CharsInfo) Equals(other *CharsInfo) bool {
	if other == nil {
		return false
	}
	return c.Value == other.Value
}

func (c *CharsInfo) ConvertWord(word string) string {
	if word == "" {
		return word
	}
	if c.IsAllLower() {
		return strings.ToLower(word)
	}
	if c.IsAllUpper() {
		return strings.ToUpper(word)
	}
	if c.IsCapitalUpper() && len(word) > 0 {
		runes := []rune(word)
		for i := 0; i < len(runes); i++ {
			if i == 0 {
				runes[0] = unicode.ToUpper(runes[0])
			} else if runes[i-1] == '-' || runes[i-1] == ' ' {
				runes[i] = unicode.ToUpper(runes[i])
			} else {
				runes[i] = unicode.ToLower(runes[i])
			}
		}
		return string(runes)
	}
	return word
}
