package MorphForm

import (
	"strings"
)

// MorphLang представляет язык (флаговый тип).
type MorphLang struct {
	Value int16
}

// Список имён языков по позициям флагов.
var langNames = []string{"RU", "UA", "BY", "EN", "IT", "KZ"}

// Глобальные переменные (аналог статических в C#).
var (
	MorphLangUnknown = MorphLang{Value: 0}
	MorphLangRU      = MorphLang{Value: 1 << 0}
	MorphLangUA      = MorphLang{Value: 1 << 1}
	MorphLangBY      = MorphLang{Value: 1 << 2}
	MorphLangEN      = MorphLang{Value: 1 << 3}
	MorphLangIT      = MorphLang{Value: 1 << 4}
	MorphLangKZ      = MorphLang{Value: 1 << 5}
)

// Copy возвращает копию структуры.
func (l MorphLang) Copy() MorphLang {
	return MorphLang{Value: l.Value}
}

// Equals сравнивает два MorphLang.
func (l MorphLang) Equals(other MorphLang) bool {
	return l.Value == other.Value
}

// And реализует побитное пересечение.
func (l MorphLang) And(other MorphLang) MorphLang {
	return MorphLang{Value: l.Value & other.Value}
}

// Or реализует побитное объединение.
func (l MorphLang) Or(other MorphLang) MorphLang {
	return MorphLang{Value: l.Value | other.Value}
}

// IsUndefined возвращает true, если не установлен ни один язык.
func (l MorphLang) IsUndefined() bool {
	return l.Value == 0
}

// Геттеры языков.
func (l MorphLang) IsRu() bool { return (l.Value & (1 << 0)) != 0 }
func (l MorphLang) IsUa() bool { return (l.Value & (1 << 1)) != 0 }
func (l MorphLang) IsBy() bool { return (l.Value & (1 << 2)) != 0 }
func (l MorphLang) IsEn() bool { return (l.Value & (1 << 3)) != 0 }
func (l MorphLang) IsIt() bool { return (l.Value & (1 << 4)) != 0 }
func (l MorphLang) IsKz() bool { return (l.Value & (1 << 5)) != 0 }

// Set методы (при необходимости).
func (l *MorphLang) SetRu(v bool) { l.setValue(0, v) }
func (l *MorphLang) SetUa(v bool) { l.setValue(1, v) }
func (l *MorphLang) SetBy(v bool) { l.setValue(2, v) }
func (l *MorphLang) SetEn(v bool) { l.setValue(3, v) }
func (l *MorphLang) SetIt(v bool) { l.setValue(4, v) }
func (l *MorphLang) SetKz(v bool) { l.setValue(5, v) }

func (l *MorphLang) setValue(i int, v bool) {
	if v {
		l.Value |= 1 << i
	} else {
		l.Value &= ^(1 << i)
	}
}

// IsCyrillic возвращает true, если это кириллический язык.
func (l MorphLang) IsCyrillic() bool {
	return l.IsRu() || l.IsUa() || l.IsBy() || l.IsKz()
}

// ToString возвращает строковое представление: "RU;UA;EN"
func (l MorphLang) String() string {
	var res []string
	if l.IsRu() {
		res = append(res, "RU")
	}
	if l.IsUa() {
		res = append(res, "UA")
	}
	if l.IsBy() {
		res = append(res, "BY")
	}
	if l.IsEn() {
		res = append(res, "EN")
	}
	if l.IsIt() {
		res = append(res, "IT")
	}
	if l.IsKz() {
		res = append(res, "KZ")
	}
	return strings.Join(res, ";")
}

// TryParse создаёт MorphLang из строки "RU;UA".
func TryParse(s string) (MorphLang, bool) {
	var result MorphLang
	for _, part := range strings.Split(strings.ToUpper(s), ";") {
		switch part {
		case "RU":
			result.SetRu(true)
		case "UA":
			result.SetUa(true)
		case "BY":
			result.SetBy(true)
		case "EN":
			result.SetEn(true)
		case "IT":
			result.SetIt(true)
		case "KZ":
			result.SetKz(true)
		}
	}
	if result.IsUndefined() {
		return result, false
	}
	return result, true
}

// Without возвращает копию языка без указанных языков (можно передать несколько).
func (ml MorphLang) Without(others ...MorphLang) MorphLang {
	res := ml.Value
	for _, other := range others {
		res &^= other.Value // удаляем флаги other из res
	}
	return MorphLang{Value: res}
}
