package MorphForm

import "strings"

// MorphPerson описывает грамматическое лицо: 1-е, 2-е, 3-е.
// Может быть флаговым типом (например, если определено несколько вариантов).
type MorphPerson uint16

// Значения MorphPerson.
const (
	MorphPersonUndefined MorphPerson = 0
	MorphPersonFirst     MorphPerson = 1 << iota // первое лицо
	MorphPersonSecond                            // второе лицо
	MorphPersonThird                             // третье лицо
)

// And возвращает пересечение лиц.
func (p MorphPerson) And(other MorphPerson) MorphPerson {
	return p & other
}

// Or возвращает объединение лиц.
func (p MorphPerson) Or(other MorphPerson) MorphPerson {
	return p | other
}

// IsDefined возвращает true, если лицо определено (не Undefined).
func (p MorphPerson) IsDefined() bool {
	return p != MorphPersonUndefined
}

// IsUndefined возвращает true, если лицо определено (не Undefined).
func (p MorphPerson) IsUndefined() bool {
	return p == MorphPersonUndefined
}

// String возвращает строковое представление лица (например, "1лицо|3лицо").
func (p MorphPerson) String() string {
	var res []string
	if p.And(MorphPersonFirst).IsDefined() {
		res = append(res, "1лицо")
	}
	if p.And(MorphPersonSecond).IsDefined() {
		res = append(res, "2лицо")
	}
	if p.And(MorphPersonThird).IsDefined() {
		res = append(res, "3лицо")
	}
	return strings.Join(res, "|")
}
