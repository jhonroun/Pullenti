package morphinternal

import (
	"fmt"
	"strings"
	"unicode"
)

type UnicodeInfo struct {
	UniChar rune
	Code    int
	value   uint16
}

var allChars []*UnicodeInfo
var inited bool

func GetChar(ch rune) *UnicodeInfo {
	if !inited {
		Initialize()
	}
	if int(ch) >= 0x10000 {
		ch = '?'
	}
	return allChars[int(ch)]
}

func Initialize() {
	if inited {
		return
	}
	inited = true
	allChars = make([]*UnicodeInfo, 0x10000)
	cyrVowel := "АЕЁИОУЮЯЫЭЄІЇЎӘӨҰҮІ" + strings.ToLower("АЕЁИОУЮЯЫЭЄІЇЎӘӨҰҮІ")

	for i := 0; i < 0x10000; i++ {
		ch := rune(i)
		ui := &UnicodeInfo{UniChar: ch, Code: i}

		if unicode.IsSpace(ch) {
			ui.SetWhitespace(true)
		} else if unicode.IsDigit(ch) {
			ui.SetDigit(true)
		} else if unicode.IsLetter(ch) {
			ui.SetLetter(true)
			if i >= 0x400 && i < 0x500 {
				ui.SetCyrillic(true)
				if strings.ContainsRune(cyrVowel, ch) {
					ui.SetVowel(true)
				}
			} else if i < 0x200 {
				ui.SetLatin(true)
				if strings.ContainsRune("AEIOUYaeiouy", ch) {
					ui.SetVowel(true)
				}
			}
			if unicode.IsUpper(ch) {
				ui.SetUpper(true)
			}
			if unicode.IsLower(ch) {
				ui.SetLower(true)
			}
		} else {
			if strings.ContainsRune("-–—¬\u00AD\u2011−", ch) {
				ui.SetHiphen(true)
			}
			if strings.ContainsRune(`"'`+"`“”’", ch) {
				ui.SetQuot(true)
			}
			if strings.ContainsRune(`'`+"`’", ch) {
				ui.SetApos(true)
				ui.SetQuot(true)
			}
		}
		if i >= 0x300 && i < 0x370 {
			ui.SetUdaren(true)
		}
		allChars[i] = ui
	}
}

// Флаги: каждый флаг — бит от 0 до 11
const (
	bWhitespace = 1 << iota
	bDigit
	bLetter
	bUpper
	bLower
	bLatin
	bCyrillic
	bHiphen
	bVowel
	bQuot
	bApos
	bUdaren
)

// Геттеры
func (u *UnicodeInfo) IsWhitespace() bool { return u.value&bWhitespace != 0 }
func (u *UnicodeInfo) IsDigit() bool      { return u.value&bDigit != 0 }
func (u *UnicodeInfo) IsLetter() bool     { return u.value&bLetter != 0 }
func (u *UnicodeInfo) IsUpper() bool      { return u.value&bUpper != 0 }
func (u *UnicodeInfo) IsLower() bool      { return u.value&bLower != 0 }
func (u *UnicodeInfo) IsLatin() bool      { return u.value&bLatin != 0 }
func (u *UnicodeInfo) IsCyrillic() bool   { return u.value&bCyrillic != 0 }
func (u *UnicodeInfo) IsHiphen() bool     { return u.value&bHiphen != 0 }
func (u *UnicodeInfo) IsVowel() bool      { return u.value&bVowel != 0 }
func (u *UnicodeInfo) IsQuot() bool       { return u.value&bQuot != 0 }
func (u *UnicodeInfo) IsApos() bool       { return u.value&bApos != 0 }
func (u *UnicodeInfo) IsUdaren() bool     { return u.value&bUdaren != 0 }

// Сеттеры
func (u *UnicodeInfo) SetWhitespace(v bool) { u.setFlag(bWhitespace, v) }
func (u *UnicodeInfo) SetDigit(v bool)      { u.setFlag(bDigit, v) }
func (u *UnicodeInfo) SetLetter(v bool)     { u.setFlag(bLetter, v) }
func (u *UnicodeInfo) SetUpper(v bool)      { u.setFlag(bUpper, v) }
func (u *UnicodeInfo) SetLower(v bool)      { u.setFlag(bLower, v) }
func (u *UnicodeInfo) SetLatin(v bool)      { u.setFlag(bLatin, v) }
func (u *UnicodeInfo) SetCyrillic(v bool)   { u.setFlag(bCyrillic, v) }
func (u *UnicodeInfo) SetHiphen(v bool)     { u.setFlag(bHiphen, v) }
func (u *UnicodeInfo) SetVowel(v bool)      { u.setFlag(bVowel, v) }
func (u *UnicodeInfo) SetQuot(v bool)       { u.setFlag(bQuot, v) }
func (u *UnicodeInfo) SetApos(v bool)       { u.setFlag(bApos, v) }
func (u *UnicodeInfo) SetUdaren(v bool)     { u.setFlag(bUdaren, v) }

func (u *UnicodeInfo) setFlag(bit uint16, v bool) {
	if v {
		u.value |= bit
	} else {
		u.value &= ^bit
	}
}

// Отладочное представление
func (u *UnicodeInfo) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("'%c'(%d)", u.UniChar, u.Code))
	if u.IsWhitespace() {
		sb.WriteString(", whitespace")
	}
	if u.IsDigit() {
		sb.WriteString(", digit")
	}
	if u.IsLetter() {
		sb.WriteString(", letter")
	}
	if u.IsLatin() {
		sb.WriteString(", latin")
	}
	if u.IsCyrillic() {
		sb.WriteString(", cyrillic")
	}
	if u.IsUpper() {
		sb.WriteString(", upper")
	}
	if u.IsLower() {
		sb.WriteString(", lower")
	}
	if u.IsHiphen() {
		sb.WriteString(", hiphen")
	}
	if u.IsQuot() {
		sb.WriteString(", quot")
	}
	if u.IsApos() {
		sb.WriteString(", apos")
	}
	if u.IsVowel() {
		sb.WriteString(", vowel")
	}
	if u.IsUdaren() {
		sb.WriteString(", udaren")
	}
	return sb.String()
}

func GetUnicodeInfo(ch rune) *UnicodeInfo {
	return GetChar(ch)
}

func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func InStrArr(s string, arr ...string) bool {
	for _, item := range arr {
		if s == item {
			return true
		}
	}
	return false
}
