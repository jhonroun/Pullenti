package morph

// Залог (для глаголов)
type MorphVoice int16

const (
	// Неопределено
	MorphVoiceUndefined MorphVoice = 0
	// Действительный
	MorphVoiceActive MorphVoice = 1
	// Страдательный
	MorphVoicePassive MorphVoice = 2
	// Средний
	MorphVoiceMiddle MorphVoice = 4
)

func (v MorphVoice) String() string {
	switch v {
	case MorphVoiceActive:
		return "Active"
	case MorphVoicePassive:
		return "Passive"
	case MorphVoiceMiddle:
		return "Middle"
	default:
		return "Undefined"
	}
}
