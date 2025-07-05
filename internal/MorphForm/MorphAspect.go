package MorphForm

import "strings"

// MorphAspect описывает аспект глагола: совершенный или несовершенный.
// Может быть флаговым типом.
type MorphAspect uint16

// Значения аспектов (флаговые).
const (
	MorphAspectUndefined    MorphAspect = 0
	MorphAspectPerfective   MorphAspect = 1 << iota // совершенный вид
	MorphAspectImperfective                         // несовершенный вид
)

// And возвращает пересечение аспектов.
func (a MorphAspect) And(other MorphAspect) MorphAspect {
	return a & other
}

// Or возвращает объединение аспектов.
func (a MorphAspect) Or(other MorphAspect) MorphAspect {
	return a | other
}

// IsDefined проверяет, задан ли аспект (не Undefined).
func (a MorphAspect) IsDefined() bool {
	return a != MorphAspectUndefined
}

// String возвращает строковое представление аспектов через "|".
func (a MorphAspect) String() string {
	var res []string
	if a.And(MorphAspectPerfective).IsDefined() {
		res = append(res, "соверш.")
	}
	if a.And(MorphAspectImperfective).IsDefined() {
		res = append(res, "несоверш.")
	}
	return strings.Join(res, "|")
}
