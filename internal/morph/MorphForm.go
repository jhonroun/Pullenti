package morph

// Форма
type MorphForm int16

const (
	// Не определена
	MorphFormUndefined MorphForm = 0
	// Краткая форма
	MorphFormShort MorphForm = 1
	// Синонимичная форма
	MorphFormSynonym MorphForm = 2
)

func (f MorphForm) String() string {
	switch f {
	case MorphFormShort:
		return "Short"
	case MorphFormSynonym:
		return "Synonym"
	default:
		return "Undefined"
	}
}
