package MorphForm

import "strings"

// MorphFinite описывает глагольную форму для английских глаголов (личная, инфинитив, причастие, герундий).
type MorphFinite uint16

// Значения MorphFinite представляют собой флаги (можно комбинировать побитово).
const (
	MorphFiniteUndefined  MorphFinite = 0
	MorphFiniteFinite     MorphFinite = 1 << iota // личная форма (например: he walks)
	MorphFiniteInfinitive                         // инфинитив (to walk)
	MorphFiniteParticiple                         // причастие (walking, walked)
	MorphFiniteGerund                             // герундий (walking как существительное)
)

// And возвращает пересечение двух множеств признаков.
func (f MorphFinite) And(other MorphFinite) MorphFinite {
	return f & other
}

// Or возвращает объединение двух множеств признаков.
func (f MorphFinite) Or(other MorphFinite) MorphFinite {
	return f | other
}

// IsDefined возвращает true, если хотя бы один флаг установлен.
func (f MorphFinite) IsDefined() bool {
	return f != MorphFiniteUndefined
}

// String возвращает строковое представление установленных флагов через "|".
func (f MorphFinite) String() string {
	var res []string
	if f.And(MorphFiniteFinite).IsDefined() {
		res = append(res, "finite")
	}
	if f.And(MorphFiniteGerund).IsDefined() {
		res = append(res, "gerund")
	}
	if f.And(MorphFiniteInfinitive).IsDefined() {
		res = append(res, "infinitive")
	}
	if f.And(MorphFiniteParticiple).IsDefined() {
		res = append(res, "participle")
	}
	return strings.Join(res, "|")
}
