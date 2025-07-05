package MorphForm

import "strings"

// MorphGender описывает грамматический род: мужской, женский, средний.
// Может быть флаговым (например: мужской|женский).
type MorphGender uint16

// Значения MorphGender (можно комбинировать).
const (
	MorphGenderUndefined MorphGender = 0
	MorphGenderMasculine MorphGender = 1 << iota // мужской род
	MorphGenderFeminine                          // женский род
	MorphGenderNeuter                            // средний род
)

// And возвращает пересечение родов.
func (g MorphGender) And(other MorphGender) MorphGender {
	return g & other
}

// Or возвращает объединение родов.
func (g MorphGender) Or(other MorphGender) MorphGender {
	return g | other
}

// IsDefined возвращает true, если род задан (не Undefined).
func (g MorphGender) IsDefined() bool {
	return g != MorphGenderUndefined
}

// String возвращает строковое представление родов (через "|").
func (g MorphGender) String() string {
	var res []string
	if g.And(MorphGenderMasculine).IsDefined() {
		res = append(res, "муж.")
	}
	if g.And(MorphGenderFeminine).IsDefined() {
		res = append(res, "жен.")
	}
	if g.And(MorphGenderNeuter).IsDefined() {
		res = append(res, "средн.")
	}
	return strings.Join(res, "|")
}
