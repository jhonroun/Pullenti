package morph

// Аспект (для глаголов)
type MorphAspect int16

const (
	// Неопределено
	MorphAspectUndefined MorphAspect = 0
	// Совершенный
	MorphAspectPerfective MorphAspect = 1
	// Несовершенный
	MorphAspectImperfective MorphAspect = 2
)

func (a MorphAspect) String() string {
	switch a {
	case MorphAspectPerfective:
		return "Perfective"
	case MorphAspectImperfective:
		return "Imperfective"
	default:
		return "Undefined"
	}
}
