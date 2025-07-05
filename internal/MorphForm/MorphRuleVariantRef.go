package MorphForm

import (
	"fmt"
)

// MorphRuleVariantRef ссылается на вариант морфологического правила с коэффициентом.
type MorphRuleVariantRef struct {
	RuleId    int   // ID правила
	VariantId int16 // ID варианта
	Coef      int16 // Коэффициент (чем выше, тем приоритетнее)
}

// NewMorphRuleVariantRef — конструктор.
func NewMorphRuleVariantRef(rid int, vid int16, co int16) *MorphRuleVariantRef {
	return &MorphRuleVariantRef{
		RuleId:    rid,
		VariantId: vid,
		Coef:      co,
	}
}

// String возвращает строковое представление: "RuleId VariantId".
func (r *MorphRuleVariantRef) String() string {
	return fmt.Sprintf("%d %d", r.RuleId, r.VariantId)
}

// CompareTo сравнивает по коэффициенту (по убыванию).
func (r *MorphRuleVariantRef) CompareTo(other *MorphRuleVariantRef) int {
	if r.Coef > other.Coef {
		return -1
	}
	if r.Coef < other.Coef {
		return 1
	}
	return 0
}
