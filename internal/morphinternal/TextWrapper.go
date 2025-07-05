package morphinternal

import (
	"strings"
	"unicode/utf8"
)

type TextWrapper struct {
	Text   string
	Length int
	Chars  *CharsList
}

func NewTextWrapper(txt string, toUpper bool) *TextWrapper {
	if toUpper {
		txt = strings.ToUpper(txt)
	}
	return &TextWrapper{
		Text:   txt,
		Length: utf8.RuneCountInString(txt),
		Chars:  NewCharsList(txt),
	}
}

func (tw *TextWrapper) String() string {
	return tw.Text
}

// CharsList предоставляет доступ к UnicodeInfo для каждого символа строки
type CharsList struct {
	Text  string
	Count int
}

func NewCharsList(txt string) *CharsList {
	return &CharsList{
		Text:  txt,
		Count: utf8.RuneCountInString(txt),
	}
}

func (cl *CharsList) At(ind int) *UnicodeInfo {
	// получаем i-й rune из строки
	i := 0
	for pos, r := range cl.Text {
		if i == ind {
			return GetChar(r)
		}
		i++
		// при выходе за пределы индекса
		if pos > len(cl.Text) {
			break
		}
	}
	return GetChar('?') // fallback
}
