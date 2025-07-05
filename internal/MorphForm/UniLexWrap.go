package MorphForm

// UniLexWrap — вспомогательная структура, содержащая список словоформ и язык.
type UniLexWrap struct {
	Lang      MorphLang        // Язык словоформ
	WordForms []*MorphWordForm // Список словоформ
}

// NewUniLexWrap — конструктор.
func NewUniLexWrap(lang MorphLang) *UniLexWrap {
	return &UniLexWrap{
		Lang:      lang,
		WordForms: []*MorphWordForm{},
	}
}
