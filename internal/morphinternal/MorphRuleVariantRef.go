package morphinternal

import (
	"fmt"
)

type MorphRuleVariantRef struct {
	RuleId    int
	VariantId int16
	Coef      int16
}

func NewMorphRuleVariantRef(rid int, vid int16, coef int16) *MorphRuleVariantRef {
	return &MorphRuleVariantRef{
		RuleId:    rid,
		VariantId: vid,
		Coef:      coef,
	}
}

func (m *MorphRuleVariantRef) String() string {
	return fmt.Sprintf("%d %d", m.RuleId, m.VariantId)
}

// Для сортировки по Coef (по убыванию)
type MorphRuleVariantRefList []*MorphRuleVariantRef

func (list MorphRuleVariantRefList) Len() int {
	return len(list)
}

func (list MorphRuleVariantRefList) Less(i, j int) bool {
	return list[i].Coef > list[j].Coef // сортировка по убыванию
}

func (list MorphRuleVariantRefList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

// Использование:
// sort.Sort(MorphRuleVariantRefList(yourSlice))
