package morphinternal

import (
	"fmt"
)

type MorphTreeNode struct {
	Nodes           map[int16]*MorphTreeNode
	RuleIds         []int
	ReverceVariants []*MorphRuleVariantRef
	LazyPos         int
}

func (m *MorphTreeNode) CalcTotalNodes() int {
	res := 0
	if m.Nodes != nil {
		for _, v := range m.Nodes {
			res += v.CalcTotalNodes() + 1
		}
	}
	return res
}

func (m *MorphTreeNode) String() string {
	count := 0
	if m.RuleIds != nil {
		count = len(m.RuleIds)
	}
	return fmt.Sprintf("? (%d, %d)", m.CalcTotalNodes(), count)
}

func (m *MorphTreeNode) deserializeBase(str *ByteArrayWrapper, pos *int) {
	cou := str.DeserializeShort(pos)
	if cou > 0 {
		m.RuleIds = make([]int, 0, cou)
		for i := 0; i < int(cou); i++ {
			id := str.DeserializeShort(pos)
			if id != 0 {
				m.RuleIds = append(m.RuleIds, int(id))
			}
		}
	}

	cou = str.DeserializeShort(pos)
	if cou > 0 {
		m.ReverceVariants = make([]*MorphRuleVariantRef, 0, cou)
		for i := 0; i < int(cou); i++ {
			rid := str.DeserializeShort(pos)
			id := str.DeserializeShort(pos)
			co := str.DeserializeShort(pos)
			m.ReverceVariants = append(m.ReverceVariants, &MorphRuleVariantRef{
				RuleId:    int(rid),
				VariantId: int16(id),
				Coef:      int16(co),
			})
		}
	}
}

func (m *MorphTreeNode) Deserialize(str *ByteArrayWrapper, pos *int) int {
	res := 0
	m.deserializeBase(str, pos)
	cou := str.DeserializeShort(pos)
	if cou > 0 {
		m.Nodes = make(map[int16]*MorphTreeNode)
		for i := 0; i < int(cou); i++ {
			key := str.DeserializeShort(pos)
			_ = str.DeserializeInt(pos) // skip 'pp' field, not used in full load
			child := &MorphTreeNode{}
			res1 := child.Deserialize(str, pos)
			res += res1 + 1
			m.Nodes[int16(key)] = child
		}
	}
	return res
}

func (m *MorphTreeNode) DeserializeLazy(str *ByteArrayWrapper, me *MorphEngine, pos *int) {
	m.deserializeBase(str, pos)
	cou := str.DeserializeShort(pos)
	if cou > 0 {
		m.Nodes = make(map[int16]*MorphTreeNode)
		for i := 0; i < int(cou); i++ {
			key := str.DeserializeShort(pos)
			pp := str.DeserializeInt(pos)
			child := &MorphTreeNode{LazyPos: *pos}
			m.Nodes[int16(key)] = child
			*pos = pp
		}
	}

	savedPos := *pos
	if m.RuleIds != nil {
		for _, rid := range m.RuleIds {
			r := me.GetMutRule(rid)
			if r != nil && r.LazyPos > 0 {
				*pos = r.LazyPos
				r.Deserialize(str, pos)
				r.LazyPos = 0
			}
		}
		*pos = savedPos
	}

	if m.ReverceVariants != nil {
		for _, rv := range m.ReverceVariants {
			r := me.GetMutRule(rv.RuleId)
			if r != nil && r.LazyPos > 0 {
				*pos = r.LazyPos
				r.Deserialize(str, pos)
				r.LazyPos = 0
			}
		}
		*pos = savedPos
	}
}
