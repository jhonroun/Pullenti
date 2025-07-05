package MorphForm

import "strings"

// MorphNumber описывает грамматическое число: единственное или множественное.
// Может быть флаговым типом (например, и то и другое при неоднозначности).
type MorphNumber uint16

// Значения MorphNumber.
const (
	MorphNumberUndefined MorphNumber = 0
	MorphNumberSingular  MorphNumber = 1 << iota // единственное число
	MorphNumberPlural                            // множественное число
)

// And возвращает пересечение чисел.
func (n MorphNumber) And(other MorphNumber) MorphNumber {
	return n & other
}

// Or возвращает объединение чисел.
func (n MorphNumber) Or(other MorphNumber) MorphNumber {
	return n | other
}

// IsDefined возвращает true, если значение задано.
func (n MorphNumber) IsDefined() bool {
	return n != MorphNumberUndefined
}

// String возвращает строковое представление числа (например, "единств.|множеств.")
func (n MorphNumber) String() string {
	var res []string
	if n.And(MorphNumberSingular).IsDefined() {
		res = append(res, "единств.")
	}
	if n.And(MorphNumberPlural).IsDefined() {
		res = append(res, "множеств.")
	}
	return strings.Join(res, "|")
}
