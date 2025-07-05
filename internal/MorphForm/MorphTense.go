package MorphForm

import "strings"

// MorphTense описывает грамматическое время глагола: прошедшее, настоящее, будущее.
type MorphTense uint16

// Значения MorphTense (флаговые).
const (
	MorphTenseUndefined MorphTense = 0
	MorphTensePast      MorphTense = 1 << iota // прошедшее время
	MorphTensePresent                          // настоящее время
	MorphTenseFuture                           // будущее время
)

// And возвращает пересечение времён.
func (t MorphTense) And(other MorphTense) MorphTense {
	return t & other
}

// Or возвращает объединение времён.
func (t MorphTense) Or(other MorphTense) MorphTense {
	return t | other
}

// IsDefined возвращает true, если установлено какое-либо время.
func (t MorphTense) IsDefined() bool {
	return t != MorphTenseUndefined
}

// IsUndefined возвращает true, если установлено какое-либо время.
func (t MorphTense) IsUndefined() bool {
	return t == MorphTenseUndefined
}

// String возвращает строковое представление времени (например, "прошедшее|настоящее").
func (t MorphTense) String() string {
	var res []string
	if t.And(MorphTensePast).IsDefined() {
		res = append(res, "прошедшее")
	}
	if t.And(MorphTensePresent).IsDefined() {
		res = append(res, "настоящее")
	}
	if t.And(MorphTenseFuture).IsDefined() {
		res = append(res, "будущее")
	}
	return strings.Join(res, "|")
}
