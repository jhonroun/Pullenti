package MorphForm

import "strings"

// MorphVoice описывает залог глагола: действительный, страдательный, средний.
type MorphVoice uint16

// Значения MorphVoice (флаговый тип).
const (
	MorphVoiceUndefined MorphVoice = 0
	MorphVoiceActive    MorphVoice = 1 << iota // действительный залог
	MorphVoicePassive                          // страдательный залог
	MorphVoiceMiddle                           // средний залог
)

// And возвращает пересечение залогов.
func (v MorphVoice) And(other MorphVoice) MorphVoice {
	return v & other
}

// Or возвращает объединение залогов.
func (v MorphVoice) Or(other MorphVoice) MorphVoice {
	return v | other
}

// IsDefined возвращает true, если хотя бы один залог задан.
func (v MorphVoice) IsDefined() bool {
	return v != MorphVoiceUndefined
}

// String возвращает строковое представление залогов, например "действит.|страдат."
func (v MorphVoice) String() string {
	var res []string
	if v.And(MorphVoiceActive).IsDefined() {
		res = append(res, "действит.")
	}
	if v.And(MorphVoicePassive).IsDefined() {
		res = append(res, "страдат.")
	}
	if v.And(MorphVoiceMiddle).IsDefined() {
		res = append(res, "средн.")
	}
	return strings.Join(res, "|")
}
