package morph

// Лицо (1, 2, 3)
type MorphPerson int16

const (
	// Неопределено
	MorphPersonUndefined MorphPerson = 0
	// Первое
	MorphPersonFirst MorphPerson = 1
	// Второе
	MorphPersonSecond MorphPerson = 2
	// Третье
	MorphPersonThird MorphPerson = 4
)

func (p MorphPerson) String() string {
	switch p {
	case MorphPersonFirst:
		return "First"
	case MorphPersonSecond:
		return "Second"
	case MorphPersonThird:
		return "Third"
	default:
		return "Undefined"
	}
}
