package MorphForm

import "strings"

// MorphMood описывает наклонение глагола (изъявительное, повелительное, условное).
type MorphMood uint16

// Значения MorphMood — флаговый тип (можно комбинировать).
const (
	MorphMoodUndefined   MorphMood = 0
	MorphMoodIndicative  MorphMood = 1 << iota // изъявительное
	MorphMoodSubjunctive                       // условное
	MorphMoodImperative                        // повелительное
)

// And возвращает пересечение наклонений.
func (m MorphMood) And(other MorphMood) MorphMood {
	return m & other
}

// Or возвращает объединение наклонений.
func (m MorphMood) Or(other MorphMood) MorphMood {
	return m | other
}

// IsDefined возвращает true, если наклонение задано.
func (m MorphMood) IsDefined() bool {
	return m != MorphMoodUndefined
}

// String возвращает строковое представление наклонения.
func (m MorphMood) String() string {
	var res []string
	if m.And(MorphMoodIndicative).IsDefined() {
		res = append(res, "изъявит.")
	}
	if m.And(MorphMoodSubjunctive).IsDefined() {
		res = append(res, "условн.")
	}
	if m.And(MorphMoodImperative).IsDefined() {
		res = append(res, "повелит.")
	}
	return strings.Join(res, "|")
}
