package morph

// Число (единственное-множественное)
type MorphNumber int16

const (
	// Неопределено
	MorphNumberUndefined MorphNumber = 0
	// Единственное
	MorphNumberSingular MorphNumber = 1
	// Множественное
	MorphNumberPlural MorphNumber = 2
)

func (n MorphNumber) String() string {
	switch n {
	case MorphNumberSingular:
		return "Singular"
	case MorphNumberPlural:
		return "Plural"
	default:
		return "Undefined"
	}
}
