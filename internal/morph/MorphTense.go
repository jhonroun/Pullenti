package morph

// Время (для глаголов)
type MorphTense int16

const (
	// Неопределено
	MorphTenseUndefined MorphTense = 0
	// Прошлое
	MorphTensePast MorphTense = 1
	// Настоящее
	MorphTensePresent MorphTense = 2
	// Будущее
	MorphTenseFuture MorphTense = 4
)

func (t MorphTense) String() string {
	switch t {
	case MorphTensePast:
		return "Past"
	case MorphTensePresent:
		return "Present"
	case MorphTenseFuture:
		return "Future"
	default:
		return "Undefined"
	}
}
