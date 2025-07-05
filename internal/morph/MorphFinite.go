package morph

// Для английских глаголов
type MorphFinite int16

const (
	// Неопределено
	MorphFiniteUndefined MorphFinite = 0
	// Финитная форма
	MorphFiniteFinite MorphFinite = 1
	// Инфинитив
	MorphFiniteInfinitive MorphFinite = 2
	// Причастие
	MorphFiniteParticiple MorphFinite = 4
	// Герундий
	MorphFiniteGerund MorphFinite = 8
)

func (m MorphFinite) String() string {
	switch m {
	case MorphFiniteFinite:
		return "Finite"
	case MorphFiniteInfinitive:
		return "Infinitive"
	case MorphFiniteParticiple:
		return "Participle"
	case MorphFiniteGerund:
		return "Gerund"
	default:
		return "Undefined"
	}
}
