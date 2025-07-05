package morph

// Род (мужской-средний-женский)
type MorphGender int16

const (
	// Неопределён
	MorphGenderUndefined MorphGender = 0
	// Мужской
	MorphGenderMasculine MorphGender = 1
	// Женский
	MorphGenderFeminine MorphGender = 2
	// Средний
	MorphGenderNeuter MorphGender = 4
)

func (g MorphGender) String() string {
	switch g {
	case MorphGenderMasculine:
		return "Masculine"
	case MorphGenderFeminine:
		return "Feminine"
	case MorphGenderNeuter:
		return "Neuter"
	default:
		return "Undefined"
	}
}
