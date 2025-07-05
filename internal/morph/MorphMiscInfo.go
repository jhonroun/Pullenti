package morph

import (
	"strings"

	"github.com/jhonroun/pullenti/internal/morphinternal"
)

type MorphMiscInfo struct {
	Attrs []string
	Value int16
	Id    int
}

func (mi *MorphMiscInfo) AddAttr(a string) {
	for _, x := range mi.Attrs {
		if x == a {
			return
		}
	}
	mi.Attrs = append(mi.Attrs, a)
}

func (mi *MorphMiscInfo) getBoolValue(i int) bool {
	return ((mi.Value >> i) & 1) != 0
}

func (mi *MorphMiscInfo) setBoolValue(i int, val bool) {
	if val {
		mi.Value |= (1 << i)
	} else {
		mi.Value &= ^(1 << i)
	}
}

func (mi *MorphMiscInfo) CopyFrom(src *MorphMiscInfo) {
	mi.Value = src.Value
	mi.Attrs = append(mi.Attrs, src.Attrs...)
}

func (mi *MorphMiscInfo) Person() MorphPerson {
	res := MorphPersonUndefined
	for _, a := range mi.Attrs {
		switch a {
		case "1 л.":
			res |= MorphPersonFirst
		case "2 л.":
			res |= MorphPersonSecond
		case "3 л.":
			res |= MorphPersonThird
		}
	}
	return res
}

func (mi *MorphMiscInfo) SetPerson(val MorphPerson) {
	if val&MorphPersonFirst != 0 {
		mi.AddAttr("1 л.")
	}
	if val&MorphPersonSecond != 0 {
		mi.AddAttr("2 л.")
	}
	if val&MorphPersonThird != 0 {
		mi.AddAttr("3 л.")
	}
}

func (mi *MorphMiscInfo) Tense() MorphTense {
	for _, a := range mi.Attrs {
		switch a {
		case "п.вр.":
			return MorphTensePast
		case "н.вр.":
			return MorphTensePresent
		case "б.вр.":
			return MorphTenseFuture
		}
	}
	return MorphTenseUndefined
}

func (mi *MorphMiscInfo) SetTense(val MorphTense) {
	switch val {
	case MorphTensePast:
		mi.AddAttr("п.вр.")
	case MorphTensePresent:
		mi.AddAttr("н.вр.")
	case MorphTenseFuture:
		mi.AddAttr("б.вр.")
	}
}

func (mi *MorphMiscInfo) Aspect() MorphAspect {
	for _, a := range mi.Attrs {
		switch a {
		case "нес.в.":
			return MorphAspectImperfective
		case "сов.в.":
			return MorphAspectPerfective
		}
	}
	return MorphAspectUndefined
}

func (mi *MorphMiscInfo) SetAspect(val MorphAspect) {
	switch val {
	case MorphAspectImperfective:
		mi.AddAttr("нес.в.")
	case MorphAspectPerfective:
		mi.AddAttr("сов.в.")
	}
}

func (mi *MorphMiscInfo) Mood() MorphMood {
	for _, a := range mi.Attrs {
		if a == "пов.накл." {
			return MorphMoodImperative
		}
	}
	return MorphMoodUndefined
}

func (mi *MorphMiscInfo) SetMood(val MorphMood) {
	if val == MorphMoodImperative {
		mi.AddAttr("пов.накл.")
	}
}

func (mi *MorphMiscInfo) Voice() MorphVoice {
	for _, a := range mi.Attrs {
		switch a {
		case "дейст.з.":
			return MorphVoiceActive
		case "страд.з.":
			return MorphVoicePassive
		}
	}
	return MorphVoiceUndefined
}

func (mi *MorphMiscInfo) SetVoice(val MorphVoice) {
	switch val {
	case MorphVoiceActive:
		mi.AddAttr("дейст.з.")
	case MorphVoicePassive:
		mi.AddAttr("страд.з.")
	}
}

func (mi *MorphMiscInfo) Form() MorphForm {
	for _, a := range mi.Attrs {
		switch a {
		case "к.ф.":
			return MorphFormShort
		case "синоним.форма":
			return MorphFormSynonym
		}
	}
	if mi.IsSynonymForm() {
		return MorphFormSynonym
	}
	return MorphFormUndefined
}

func (mi *MorphMiscInfo) IsSynonymForm() bool {
	return mi.getBoolValue(0)
}

func (mi *MorphMiscInfo) SetIsSynonymForm(val bool) {
	mi.setBoolValue(0, val)
}

func (mi *MorphMiscInfo) String() string {
	if len(mi.Attrs) == 0 && mi.Value == 0 {
		return ""
	}
	res := strings.Builder{}
	if mi.IsSynonymForm() {
		res.WriteString("синоним.форма ")
	}
	for _, a := range mi.Attrs {
		res.WriteString(a + " ")
	}
	return strings.TrimSpace(res.String())
}

func (mi *MorphMiscInfo) Deserialize(str *morphinternal.ByteArrayWrapper, pos *int) {
	sh := str.DeserializeShort(pos)
	mi.Value = int16(sh)
	for {
		s := str.DeserializeString(pos)
		if s == "" {
			break
		}
		mi.AddAttr(s)
	}
}
