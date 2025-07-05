package MorphForm

import "strings"

// MorphForm описывает грамматическую форму слова, например краткую или синонимичную.
// Может содержать один или несколько признаков.
type MorphForm uint16

// Значения MorphForm (флаговые).
const (
	MorphFormUndefined MorphForm = 0
	MorphFormShort     MorphForm = 1 << iota // краткая форма (прил./прич.)
	MorphFormSynonym                         // синонимичная форма (вариант)
)

// And возвращает пересечение флагов формы.
func (f MorphForm) And(other MorphForm) MorphForm {
	return f & other
}

// Or возвращает объединение флагов формы.
func (f MorphForm) Or(other MorphForm) MorphForm {
	return f | other
}

// IsDefined возвращает true, если установлен хотя бы один флаг.
func (f MorphForm) IsDefined() bool {
	return f != MorphFormUndefined
}

// String возвращает строковое представление формы (например, "кратк.|синонимич.")
func (f MorphForm) String() string {
	var res []string
	if f.And(MorphFormShort).IsDefined() {
		res = append(res, "кратк.")
	}
	if f.And(MorphFormSynonym).IsDefined() {
		res = append(res, "синонимич.")
	}
	return strings.Join(res, "|")
}
