package MorphForm

import (
	"strings"
)

// MorphBaseInfo — базовая часть морфологической информации.
type MorphBaseInfo struct {
	mCla   MorphClass
	Gender MorphGender
	Number MorphNumber
	mCase  MorphCase
	mLang  MorphLang
}

// GetClass возвращает часть речи.
func (m *MorphBaseInfo) GetClass() MorphClass {
	return m.mCla
}

// SetClass устанавливает часть речи.
func (m *MorphBaseInfo) SetClass(v MorphClass) {
	m.mCla = v
}

// GetCase возвращает падеж.
func (m *MorphBaseInfo) GetCase() MorphCase {
	return m.mCase
}

// SetCase устанавливает падеж.
func (m *MorphBaseInfo) SetCase(v MorphCase) {
	m.mCase = v
}

// GetLanguage возвращает язык.
func (m *MorphBaseInfo) GetLanguage() MorphLang {
	return m.mLang
}

// SetLanguage устанавливает язык.
func (m *MorphBaseInfo) SetLanguage(v MorphLang) {
	m.mLang = v
}

// CopyFrom копирует все поля из другого MorphBaseInfo.
func (m *MorphBaseInfo) CopyFrom(src *MorphBaseInfo) {
	m.mCla = src.mCla.Copy()
	m.Gender = src.Gender
	m.Number = src.Number
	m.mCase = src.mCase.Copy()
	m.mLang = src.mLang.Copy()
}

// ToString возвращает строковое представление морфологической информации.
func (m *MorphBaseInfo) String() string {
	var res strings.Builder
	if !m.mCla.IsUndefined() {
		res.WriteString(m.mCla.String())
		res.WriteRune(' ')
	}
	if m.Number != MorphNumberUndefined {
		switch m.Number {
		case MorphNumberSingular:
			res.WriteString("ед.ч. ")
		case MorphNumberPlural:
			res.WriteString("мн.ч. ")
		default:
			res.WriteString("ед.мн.ч. ")
		}
	}
	if m.Gender != MorphGenderUndefined {
		switch m.Gender {
		case MorphGenderMasculine:
			res.WriteString("муж.р. ")
		case MorphGenderNeuter:
			res.WriteString("ср.р. ")
		case MorphGenderFeminine:
			res.WriteString("жен.р. ")
		case MorphGenderMasculine | MorphGenderNeuter:
			res.WriteString("муж.ср.р. ")
		case MorphGenderFeminine | MorphGenderNeuter:
			res.WriteString("жен.ср.р. ")
		case MorphGenderFeminine | MorphGenderMasculine:
			res.WriteString("муж.жен.р. ")
		case MorphGenderMasculine | MorphGenderFeminine | MorphGenderNeuter:
			res.WriteString("муж.жен.ср.р. ")
		}
	}
	if !m.mCase.IsUndefined() {
		res.WriteString(m.mCase.String())
		res.WriteRune(' ')
	}
	if !m.mLang.IsUndefined() && !m.mLang.Equals(MorphLangRU) {
		res.WriteString(m.mLang.String())
		res.WriteRune(' ')
	}
	return strings.TrimSpace(res.String())
}

// ContainsAttr — виртуальный метод. По умолчанию всегда false.
func (m *MorphBaseInfo) ContainsAttr(attrValue string, cla *MorphClass) bool {
	return false
}

// CheckAccord проверяет согласование с другой морфологической формой.
func (m *MorphBaseInfo) CheckAccord(v *MorphBaseInfo, ignoreGender bool, ignoreNumber bool) bool {
	if !v.mLang.Equals(m.mLang) {
		if v.mLang.IsUndefined() && m.mLang.IsUndefined() {
			return false
		}
	}

	num := v.Number.And(m.Number)
	if num == MorphNumberUndefined && !ignoreNumber {
		if v.Number != MorphNumberUndefined && m.Number != MorphNumberUndefined {
			if v.Number == MorphNumberSingular && v.mCase.IsGenitive() {
				if m.Number == MorphNumberPlural && m.mCase.IsGenitive() {
					if v.Gender.And(MorphGenderMasculine) == MorphGenderMasculine {
						return true
					}
				}
			}
			return false
		}
	}

	if !ignoreGender && num != MorphNumberPlural {
		if v.Gender.And(m.Gender) == MorphGenderUndefined {
			if v.Gender != MorphGenderUndefined && m.Gender != MorphGenderUndefined {
				return false
			}
		}
	}

	if v.mCase.And(m.mCase).IsUndefined() {
		if !v.mCase.IsUndefined() && !m.mCase.IsUndefined() {
			return false
		}
	}

	return true
}
