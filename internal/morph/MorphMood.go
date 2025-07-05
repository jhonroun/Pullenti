package morph

// Наклонение (для глаголов)
type MorphMood int16

const (
	// Неопределено
	MorphMoodUndefined MorphMood = 0
	// Изъявительное
	MorphMoodIndicative MorphMood = 1
	// Условное
	MorphMoodSubjunctive MorphMood = 2
	// Повелительное
	MorphMoodImperative MorphMood = 4
)

func (m MorphMood) String() string {
	switch m {
	case MorphMoodIndicative:
		return "Indicative"
	case MorphMoodSubjunctive:
		return "Subjunctive"
	case MorphMoodImperative:
		return "Imperative"
	default:
		return "Undefined"
	}
}
