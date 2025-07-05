package MorphForm

import (
	"fmt"
	"strings"
)

// MorphToken представляет токен морфологического анализа.
type MorphToken struct {
	BeginChar int              // Начальная позиция
	EndChar   int              // Конечная позиция
	Term      string           // Нормализованное представление (верхний регистр)
	WordForms []*MorphWordForm // Варианты словоформ
	CharInfo  CharsInfo        // Информация о символах
	Tag       interface{}      // Произвольный тег
	Language  MorphLang        // Внутренний язык (устанавливается вручную)
	Lemma     string           // Кэшированная лемма
}

// Length возвращает длину нормализованного слова.
func (m *MorphToken) Length() int {
	if m.Term == "" {
		return 0
	}
	return len(m.Term)
}

// GetSourceText возвращает фрагмент текста из исходной строки.
func (m *MorphToken) GetSourceText(text string) string {
	if m.BeginChar >= 0 && m.EndChar < len(text) && m.EndChar >= m.BeginChar {
		return text[m.BeginChar : m.EndChar+1]
	}
	return ""
}

// GetLemma возвращает лемму (нормальную форму) термина.
func (m *MorphToken) GetLemma() string {
	if m.Lemma != "" {
		return m.Lemma
	}

	var res string

	if len(m.WordForms) > 0 {
		if len(m.WordForms) == 1 {
			wf := m.WordForms[0]
			if wf.NormalFull != "" {
				res = wf.NormalFull
			} else {
				res = wf.NormalCase
			}
		}

		if res == "" && !m.CharInfo.IsAllLower() {
			for _, wf := range m.WordForms {
				if wf.Class.IsProperSurname() {
					s := wf.NormalFull
					if s == "" {
						s = wf.NormalCase
					}
					if EndsWithEx(s, "ОВ", "ЕВ", "", "") {
						res = s
						break
					}
				} else if wf.Class.IsProperName() && wf.IsInDictionary() {
					return wf.NormalCase
				}
			}
		}

		if res == "" {
			var best *MorphWordForm
			for _, wf := range m.WordForms {
				if best == nil || m.CompareForms(best, wf) > 0 {
					best = wf
				}
			}
			if best != nil {
				if best.NormalFull != "" {
					res = best.NormalFull
				} else {
					res = best.NormalCase
				}
			}
		}
	}

	if res != "" {
		if EndsWithEx(res, "АНЫЙ", "ЕНЫЙ", "", "") {
			res = res[:len(res)-3] + "ННЫЙ"
		} else if EndsWith(res, "ЙСЯ") {
			res = res[:len(res)-2]
		} else if EndsWith(res, "АНИЙ") && res == m.Term {
			found := false
			for _, wf := range m.WordForms {
				if wf.IsInDictionary() {
					found = true
					break
				}
			}
			if !found {
				return res[:len(res)-1] + "Е"
			}
		}
		return res
	}
	if m.Term != "" {
		return m.Term
	}
	return "?"
}

// CompareForms сравнивает две словоформы и возвращает приоритет между ними.
func (m *MorphToken) CompareForms(x, y *MorphWordForm) int {
	vx := x.NormalFull
	if vx == "" {
		vx = x.NormalCase
	}
	vy := y.NormalFull
	if vy == "" {
		vy = y.NormalCase
	}
	if vx == vy {
		return 0
	}
	if vx == "" {
		return 1
	}
	if vy == "" {
		return -1
	}
	lastx := rune(vx[len(vx)-1])
	lasty := rune(vy[len(vy)-1])

	if x.Class.IsProperSurname() && !m.CharInfo.IsAllLower() {
		if EndsWithEx(vx, "ОВ", "ЕВ", "ИН", "") {
			if !y.Class.IsProperSurname() {
				return -1
			}
		}
	}
	if y.Class.IsProperSurname() && !m.CharInfo.IsAllLower() {
		if EndsWithEx(vy, "ОВ", "ЕВ", "ИН", "") {
			if !x.Class.IsProperSurname() {
				return 1
			}
			if len(vx) > len(vy) {
				return -1
			}
			if len(vx) < len(vy) {
				return 1
			}
			return 0
		}
	}
	if x.Class.Equals(y.Class) {
		if x.Class.IsAdjective() {
			if lastx == 'Й' && lasty != 'Й' {
				return -1
			}
			if lastx != 'Й' && lasty == 'Й' {
				return 1
			}
			if !EndsWith(vx, "ОЙ") && EndsWith(vy, "ОЙ") {
				return -1
			}
			if EndsWith(vx, "ОЙ") && !EndsWith(vy, "ОЙ") {
				return 1
			}
		}
		if x.Class.IsNoun() {
			if x.Number == MorphNumberSingular && y.Number == MorphNumberPlural && len(vx) <= len(vy)+1 {
				return -1
			}
			if x.Number == MorphNumberPlural && y.Number == MorphNumberSingular && len(vx) >= len(vy)-1 {
				return 1
			}
		}
		if len(vx) < len(vy) {
			return -1
		}
		if len(vx) > len(vy) {
			return 1
		}
		return 0
	}
	if x.Class.IsAdverb() {
		return 1
	}
	if x.Class.IsNoun() && x.IsInDictionary() {
		if y.Class.IsAdjective() && y.IsInDictionary() {
			if !y.Misc.ContainsAttr("к.ф.") {
				return 1
			}
		}
		return -1
	}
	if x.Class.IsAdjective() {
		if !x.IsInDictionary() && y.Class.IsNoun() && y.IsInDictionary() {
			return 1
		}
		return -1
	}
	if x.Class.IsVerb() {
		if y.Class.IsNoun() || y.Class.IsAdjective() || y.Class.IsPreposition() {
			return 1
		}
		return -1
	}
	if y.Class.IsAdverb() {
		return -1
	}
	if y.Class.IsNoun() && y.IsInDictionary() {
		return 1
	}
	if y.Class.IsAdjective() {
		if (x.Class.IsNoun() || x.Class.IsProperSecname()) && x.IsInDictionary() {
			return -1
		}
		if x.Class.IsNoun() && !y.IsInDictionary() {
			if len(vx) < len(vy) {
				return -1
			}
		}
		return 1
	}
	if y.Class.IsVerb() {
		if x.Class.IsNoun() || x.Class.IsAdjective() || x.Class.IsPreposition() {
			return -1
		}
		if x.Class.IsProper() {
			return -1
		}
		return 1
	}
	if len(vx) < len(vy) {
		return -1
	}
	if len(vx) > len(vy) {
		return 1
	}
	return 0
}

// CompareForms сравнивает две словоформы и возвращает:
// -1 если x лучше, 1 если y лучше, 0 если равны.
func CompareForms(x, y *MorphWordForm, info CharsInfo) int {
	if x == nil {
		if y == nil {
			return 0
		}
		return 1
	}
	if y == nil {
		return -1
	}

	// 1. Нормальная форма без полной — предпочтительнее
	if x.NormalCase != x.NormalFull && y.NormalCase == y.NormalFull {
		return 1
	}
	if y.NormalCase != y.NormalFull && x.NormalCase == x.NormalFull {
		return -1
	}

	// 2. Оба — собственные имена
	if x.Class.IsProper() && y.Class.IsProper() {
		if x.Misc != nil && y.Misc != nil {
			if x.Misc.IsProperSurname && !y.Misc.IsProperSurname {
				return -1
			}
			if y.Misc.IsProperSurname && !x.Misc.IsProperSurname {
				return 1
			}
			if x.Misc.IsProperName && !y.Misc.IsProperName {
				return -1
			}
			if y.Misc.IsProperName && !x.Misc.IsProperName {
				return 1
			}
			if x.Misc.IsProperSecname && !y.Misc.IsProperSecname {
				return -1
			}
			if y.Misc.IsProperSecname && !x.Misc.IsProperSecname {
				return 1
			}
		}
	}

	// 3. Предпочтение существительных > прилагательных > глаголов
	if x.Class.IsNoun() && !y.Class.IsNoun() {
		return -1
	}
	if y.Class.IsNoun() && !x.Class.IsNoun() {
		return 1
	}
	if x.Class.IsAdjective() && !y.Class.IsAdjective() {
		return -1
	}
	if y.Class.IsAdjective() && !x.Class.IsAdjective() {
		return 1
	}
	if x.Class.IsVerb() && !y.Class.IsVerb() {
		return -1
	}
	if y.Class.IsVerb() && !x.Class.IsVerb() {
		return 1
	}

	// 4. Формы из словаря предпочтительнее
	if x.IsInDictionary() && !y.IsInDictionary() {
		return -1
	}
	if y.IsInDictionary() && !x.IsInDictionary() {
		return 1
	}

	// 5. Краткие формы прилагательных — менее предпочтительны
	if x.Misc != nil && x.Misc.IsShortForm && (y.Misc == nil || !y.Misc.IsShortForm) {
		return 1
	}
	if y.Misc != nil && y.Misc.IsShortForm && (x.Misc == nil || !x.Misc.IsShortForm) {
		return -1
	}

	// 6. Строчное написание предпочтительнее
	if info.IsAllLower() {
		if isCapitalForm(x) && !isCapitalForm(y) {
			return 1
		}
		if isCapitalForm(y) && !isCapitalForm(x) {
			return -1
		}
	}

	// 7. Лексикографическое сравнение (на всякий случай)
	if x.NormalCase < y.NormalCase {
		return -1
	}
	if x.NormalCase > y.NormalCase {
		return 1
	}

	return 0
}

// isCapitalForm проверяет, начинается ли форма с заглавной (по сути — Proper с признаком имя/фамилия).
func isCapitalForm(wf *MorphWordForm) bool {
	if wf == nil || !wf.Class.IsProper() {
		return false
	}
	if wf.Misc == nil {
		return false
	}
	return wf.Misc.IsProperSurname || wf.Misc.IsProperName || wf.Misc.IsProperSecname
}

// GetLanguage возвращает язык токена. Если явно не задан, берётся из первой словоформы.
func (m *MorphToken) GetLanguage() MorphLang {
	if !m.Language.IsUndefined() {
		return m.Language
	}
	if len(m.WordForms) > 0 {
		return m.WordForms[0].mLang
	}
	return MorphLangUnknown
}

// String возвращает отладочное представление токена.
func (m *MorphToken) String() string {
	var forms []string
	for _, wf := range m.WordForms {
		forms = append(forms, wf.String())
	}
	return fmt.Sprintf("%s [%s]", m.Term, strings.Join(forms, "; "))
}
