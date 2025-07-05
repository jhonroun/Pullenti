package morph

// UniLexWrap соответствует обёртке словарных форм для конкретного языка
type UniLexWrap struct {
	Lang      MorphLang        // Язык
	WordForms []*MorphWordForm // Список словоформ
}

// NewUniLexWrap конструктор для UniLexWrap
func NewUniLexWrap(lng MorphLang) *UniLexWrap {
	return &UniLexWrap{
		Lang:      lng,
		WordForms: make([]*MorphWordForm, 0),
	}
}
