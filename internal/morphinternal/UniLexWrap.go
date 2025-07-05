package morphinternal

import "main/internal/morph"

// UniLexWrap соответствует обёртке словарных форм для конкретного языка
type UniLexWrap struct {
	Lang      morph.MorphLang        // Язык
	WordForms []*morph.MorphWordForm // Список словоформ
}

// NewUniLexWrap конструктор для UniLexWrap
func NewUniLexWrap(lng morph.MorphLang) *UniLexWrap {
	return &UniLexWrap{
		Lang:      lng,
		WordForms: make([]*morph.MorphWordForm, 0),
	}
}
