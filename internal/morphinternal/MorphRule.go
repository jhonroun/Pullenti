package morphinternal

import (
	"strings"
)

// MorphRule соответствует MorphRule из Pullenti
type MorphRule struct {
	Id        int
	Tails     []string
	MorphVars [][]*MorphRuleVariant
	LazyPos   int
}

// Add добавляет хвост и связанные варианты
func (m *MorphRule) Add(tail string, vars []*MorphRuleVariant) {
	m.Tails = append(m.Tails, tail)
	m.MorphVars = append(m.MorphVars, vars)
}

// ContainsVar проверяет наличие хвоста
func (m *MorphRule) ContainsVar(tail string) bool {
	for _, t := range m.Tails {
		if t == tail {
			return true
		}
	}
	return false
}

// GetVars возвращает список вариантов по хвосту
func (m *MorphRule) GetVars(key string) []*MorphRuleVariant {
	for i, t := range m.Tails {
		if t == key {
			return m.MorphVars[i]
		}
	}
	return nil
}

// FindVar ищет вариант по Id
func (m *MorphRule) FindVar(id int16) *MorphRuleVariant {
	for _, group := range m.MorphVars {
		for _, v := range group {
			if v.Id == id {
				return v
			}
		}
	}
	return nil
}

// String возвращает строковое представление правил
func (m *MorphRule) String() string {
	var sb strings.Builder
	for i, tail := range m.Tails {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("-" + tail)
	}
	return sb.String()
}

// Deserialize десериализует MorphRule из ByteArrayWrapper
func (m *MorphRule) Deserialize(str *ByteArrayWrapper, pos *int) {
	m.Id = int(str.DeserializeShort(pos))
	var idCounter int16 = 1
	for !str.IsEOF(*pos) {
		b := str.DeserializeByte(pos)
		if b == 0xFF {
			break
		}
		*pos-- // вернём байт обратно
		key := str.DeserializeString(pos)
		if key == "" {
			key = ""
		}
		var li []*MorphRuleVariant
		for !str.IsEOF(*pos) {
			mrv := &MorphRuleVariant{}
			if !mrv.Deserialize(str, pos) {
				break
			}
			mrv.Tail = key
			mrv.RuleId = int16(m.Id)
			mrv.Id = idCounter
			idCounter++
			li = append(li, mrv)
		}
		m.Add(key, li)
	}
}
