package MorphForm

import (
	"strings"
)

// MorphCase представляет падеж (именительный, родительный и т.д.).
type MorphCase struct {
	Value int16
}

// Побитовые флаги
const (
	CaseNominative    int = 0
	CaseGenitive      int = 1
	CaseDative        int = 2
	CaseAccusative    int = 3
	CaseInstrumental  int = 4
	CasePrepositional int = 5
	CaseVocative      int = 6
	CasePartial       int = 7
	CaseCommon        int = 8
	CasePossessive    int = 9
)

var (
	Undefined     = MorphCase{Value: 0}
	Nominative    = MorphCase{Value: 1 << CaseNominative}
	Genitive      = MorphCase{Value: 1 << CaseGenitive}
	Dative        = MorphCase{Value: 1 << CaseDative}
	Accusative    = MorphCase{Value: 1 << CaseAccusative}
	Instrumental  = MorphCase{Value: 1 << CaseInstrumental}
	Prepositional = MorphCase{Value: 1 << CasePrepositional}
	Vocative      = MorphCase{Value: 1 << CaseVocative}
	Partial       = MorphCase{Value: 1 << CasePartial}
	Common        = MorphCase{Value: 1 << CaseCommon}
	Possessive    = MorphCase{Value: 1 << CasePossessive}
	AllCases      = MorphCase{Value: 0x3FF}
)

var names = []string{
	"именит.", "родит.", "дател.", "винит.", "творит.", "предлож.",
	"зват.", "частич.", "общ.", "притяж.",
}

// GetValue возвращает true, если флаг с индексом i установлен.
func (mc MorphCase) GetValue(i int) bool {
	return ((mc.Value >> i) & 1) != 0
}

// SetValue устанавливает флаг с индексом i.
func (mc *MorphCase) SetValue(i int, val bool) {
	if val {
		mc.Value |= (1 << i)
	} else {
		mc.Value &= ^(1 << i)
	}
}

// Count возвращает количество установленных падежей.
func (mc MorphCase) Count() int {
	count := 0
	for i := 0; i < 12; i++ {
		if mc.GetValue(i) {
			count++
		}
	}
	return count
}

// Copy возвращает копию.
func (mc MorphCase) Copy() MorphCase {
	return MorphCase{Value: mc.Value}
}

// Equals сравнивает два MorphCase.
func (mc MorphCase) Equals(other MorphCase) bool {
	return mc.Value == other.Value
}

// GetHashCode возвращает хеш.
func (mc MorphCase) GetHashCode() int {
	return int(mc.Value)
}

// IsUndefined возвращает true, если ни один падеж не установлен.
func (mc MorphCase) IsUndefined() bool {
	return mc.Value == 0
}

// String возвращает строковое представление, например "именит.|родит."
func (mc MorphCase) String() string {
	var res []string
	for i := 0; i < len(names); i++ {
		if mc.GetValue(i) {
			res = append(res, names[i])
		}
	}
	return strings.Join(res, "|")
}

// Parse восстанавливает MorphCase из строки вида "родит.|дател."
func Parse(s string) MorphCase {
	var mc MorphCase
	if s == "" {
		return mc
	}
	parts := strings.Split(s, "|")
	for _, p := range parts {
		for i, n := range names {
			if p == n {
				mc.SetValue(i, true)
			}
		}
	}
	return mc
}

// Побитные операции

func (mc MorphCase) And(other MorphCase) MorphCase {
	return MorphCase{Value: mc.Value & other.Value}
}

func (mc MorphCase) Or(other MorphCase) MorphCase {
	return MorphCase{Value: mc.Value | other.Value}
}

func (mc MorphCase) Xor(other MorphCase) MorphCase {
	return MorphCase{Value: mc.Value ^ other.Value}
}

// Методы доступа

func (mc MorphCase) IsNominative() bool    { return mc.GetValue(CaseNominative) }
func (mc *MorphCase) SetNominative(v bool) { mc.SetValue(CaseNominative, v) }

func (mc MorphCase) IsGenitive() bool    { return mc.GetValue(CaseGenitive) }
func (mc *MorphCase) SetGenitive(v bool) { mc.SetValue(CaseGenitive, v) }

func (mc MorphCase) IsDative() bool    { return mc.GetValue(CaseDative) }
func (mc *MorphCase) SetDative(v bool) { mc.SetValue(CaseDative, v) }

func (mc MorphCase) IsAccusative() bool    { return mc.GetValue(CaseAccusative) }
func (mc *MorphCase) SetAccusative(v bool) { mc.SetValue(CaseAccusative, v) }

func (mc MorphCase) IsInstrumental() bool    { return mc.GetValue(CaseInstrumental) }
func (mc *MorphCase) SetInstrumental(v bool) { mc.SetValue(CaseInstrumental, v) }

func (mc MorphCase) IsPrepositional() bool    { return mc.GetValue(CasePrepositional) }
func (mc *MorphCase) SetPrepositional(v bool) { mc.SetValue(CasePrepositional, v) }

func (mc MorphCase) IsVocative() bool    { return mc.GetValue(CaseVocative) }
func (mc *MorphCase) SetVocative(v bool) { mc.SetValue(CaseVocative, v) }

func (mc MorphCase) IsPartial() bool    { return mc.GetValue(CasePartial) }
func (mc *MorphCase) SetPartial(v bool) { mc.SetValue(CasePartial, v) }

func (mc MorphCase) IsCommon() bool    { return mc.GetValue(CaseCommon) }
func (mc *MorphCase) SetCommon(v bool) { mc.SetValue(CaseCommon, v) }

func (mc MorphCase) IsPossessive() bool    { return mc.GetValue(CasePossessive) }
func (mc *MorphCase) SetPossessive(v bool) { mc.SetValue(CasePossessive, v) }
