package MorphForm

import (
	"strings"

	"github.com/jhonroun/pullenti/internal/morphinternal"
)

var (
	LatChars      = "ABEKMHOPCTYXIaekmopctyxi"
	CyrChars      = "АВЕКМНОРСТУХІаекморстухі"
	GreekChars    = "ΑΒΓΕΗΙΚΛΜΟПΡТΥΦΧ"
	CyrGreekChars = "АВГЕНІКЛМОПРТУФХ"
	UdarChars     = "ÀÁÈÉËÒÓàáèéëýÝòóЀѐЍѝỲỳ"
	UdarCyrChars  = "ААЕЕЕООааеееуУооЕеИиУу"
	Rus0          = "–ЁѐЀЍѝЎўӢӣ"
	Rus1          = "-ЕЕЕИИУУЙЙ"
	Preps         [][]string
	Cases         []MorphCase
	PrepCases     map[string]MorphCase
	PrepNormsSrc  []string
	PrepNorms     map[string]string
)

func init() {
	Preps = [][]string{
		strings.Split("БЕЗ;ДО;ИЗ;ИЗЗА;ОТ;У;ДЛЯ;РАДИ;ВОЗЛЕ;ПОЗАДИ;ВПЕРЕДИ;БЛИЗ;ВБЛИЗИ;ВГЛУБЬ;ВВИДУ;ВДОЛЬ;ВЗАМЕН;ВКРУГ;ВМЕСТО;ВНЕ;ВНИЗУ;ВНУТРИ;ВНУТРЬ;ВОКРУГ;ВРОДЕ;ВСЛЕД;ВСЛЕДСТВИЕ;ЗАМЕСТО;ИЗНУТРИ;КАСАТЕЛЬНО;КРОМЕ;МИМО;НАВРОДЕ;НАЗАД;НАКАНУНЕ;НАПОДОБИЕ;НАПРОТИВ;НАСЧЕТ;ОКОЛО;ОТНОСИТЕЛЬНО;ПОВЕРХ;ПОДЛЕ;ПОМИМО;ПОПЕРЕК;ПОРЯДКА;ПОСЕРЕДИНЕ;ПОСРЕДИ;ПОСЛЕ;ПРЕВЫШЕ;ПРЕЖДЕ;ПРОТИВ;СВЕРХ;СВЫШЕ;СНАРУЖИ;СРЕДИ;СУПРОТИВ;ПУТЕМ;ПОСРЕДСТВОМ", ";"),
		strings.Split("К;БЛАГОДАРЯ;ВОПРЕКИ;НАВСТРЕЧУ;СОГЛАСНО;СООБРАЗНО;ПАРАЛЛЕЛЬНО;ПОДОБНО;СООТВЕТСТВЕННО;СОРАЗМЕРНО", ";"),
		strings.Split("ПРО;ЧЕРЕЗ;СКВОЗЬ;СПУСТЯ", ";"),
		strings.Split("НАД;ПЕРЕД;ПРЕД", ";"),
		{"ПРИ"},
		strings.Split("В;НА;О;ВКЛЮЧАЯ", ";"),
		{"МЕЖДУ"},
		strings.Split("ЗА;ПОД", ";"),
		{"ПО"},
		{"С"},
	}
	Cases = []MorphCase{
		Genitive,
		Dative,
		Accusative,
		Instrumental,
		Prepositional,
		Accusative.Or(Prepositional),
		Genitive.Or(Instrumental),
		Accusative.Or(Instrumental),
		Dative.Or(Accusative).Or(Prepositional),
		Genitive.Or(Accusative).Or(Instrumental),
	}
	PrepCases = make(map[string]MorphCase)
	for i, preps := range Preps {
		for _, prep := range preps {
			PrepCases[prep] = Cases[i]
		}
	}
	PrepNormsSrc = []string{
		"БЕЗ;БЕЗО", "ВБЛИЗИ;БЛИЗ", "В;ВО", "ВОКРУГ;ВКРУГ", "ВНУТРИ;ВНУТРЬ;ВОВНУТРЬ",
		"ВПЕРЕДИ;ВПЕРЕД", "ВСЛЕД;ВОСЛЕД", "ВМЕСТО;ЗАМЕСТО", "ИЗ;ИЗО", "К;КО",
		"МЕЖДУ;МЕЖ;ПРОМЕЖДУ;ПРОМЕЖ", "НАД;НАДО", "О;ОБ;ОБО", "ОТ;ОТО", "ПЕРЕД;ПРЕД;ПРЕДО;ПЕРЕДО",
		"ПОД;ПОДО", "ПОСЕРЕДИНЕ;ПОСРЕДИ;ПОСЕРЕДЬ", "С;СО", "СРЕДИ;СРЕДЬ;СЕРЕДЬ", "ЧЕРЕЗ;ЧРЕЗ",
	}
	PrepNorms = make(map[string]string)
	for _, grp := range PrepNormsSrc {
		parts := strings.Split(grp, ";")
		if len(parts) == 0 {
			continue
		}
		base := parts[0]
		for _, v := range parts[1:] {
			PrepNorms[v] = base
		}
	}
}

func CorrectWord(w string) string {
	if w == "" {
		return ""
	}
	res := strings.ToUpper(w)
	changed := false
	for _, ch := range res {
		if strings.ContainsRune(Rus0, ch) {
			changed = true
			break
		}
	}
	if changed {
		sb := strings.Builder{}
		for _, ch := range res {
			if i := strings.IndexRune(Rus0, ch); i >= 0 {
				sb.WriteByte(Rus1[i])
			} else {
				sb.WriteRune(ch)
			}
		}
		res = sb.String()
	}
	if strings.ContainsRune(res, 0x00AD) {
		res = strings.ReplaceAll(res, string(rune(0x00AD)), "-")
	}
	if strings.HasPrefix(res, "АГЕНС") {
		res = "АГЕНТС" + res[5:]
	}
	return res
}

func IsLatinChar(ch rune) bool {
	ui := morphinternal.GetChar(ch)
	return ui.IsLatin()
}

func IsCyrillicChar(ch rune) bool {
	ui := morphinternal.GetChar(ch)
	return ui.IsCyrillic()
}

func IsCyrillic(str string) bool {
	if str == "" {
		return false
	}
	for _, ch := range str {
		if !IsCyrillicChar(ch) {
			if ch != '-' && ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
				return false
			}
		}
	}
	return true
}

func GetLanguageForText(text string) string {
	if text == "" {
		return ""
	}
	ruChars := 0
	enChars := 0
	for _, ch := range text {
		if ch >= 0x400 && ch < 0x500 {
			ruChars++
		} else if ch < 0x80 && ((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')) {
			enChars++
		}
	}
	if ruChars > enChars*2 && ruChars > 10 {
		return "ru"
	}
	if ruChars > 0 && enChars == 0 {
		return "ru"
	}
	if enChars > 0 {
		return "en"
	}
	return ""
}

// IsHiphen проверяет, является ли символ дефисом.
func IsHiphen(ch rune) bool {
	ui := morphinternal.GetUnicodeInfo(ch)
	return ui.IsHiphen()
}

// IsCyrillicVowel проверяет, является ли символ кириллической гласной буквой.
func IsCyrillicVowel(ch rune) bool {
	ui := morphinternal.GetUnicodeInfo(ch)
	return ui.IsCyrillic() && ui.IsVowel()
}

// IsLatinVowel проверяет, является ли символ латинской гласной буквой.
func IsLatinVowel(ch rune) bool {
	ui := morphinternal.GetUnicodeInfo(ch)
	return ui.IsLatin() && ui.IsVowel()
}

// GetCyrForLat возвращает соответствующую кириллическую букву для латинской.
func GetCyrForLat(lat rune) rune {
	i := strings.IndexRune(LatChars, lat)
	if i >= 0 && i < len(CyrChars) {
		return rune(CyrChars[i])
	}
	i = strings.IndexRune(GreekChars, lat)
	if i >= 0 && i < len(CyrGreekChars) {
		return rune(CyrGreekChars[i])
	}
	return 0
}

// GetLatForCyr возвращает соответствующую латинскую букву для кириллической.
func GetLatForCyr(cyr rune) rune {
	i := strings.IndexRune(CyrChars, cyr)
	if i >= 0 && i < len(LatChars) {
		return rune(LatChars[i])
	}
	return 0
}

// IsApos проверяет, является ли символ апострофом (одиночная кавычка или знак ударения).
func IsApos(ch rune) bool {
	ui := morphinternal.GetUnicodeInfo(ch)
	return ui.IsApos()
}

// IsQuote проверяет, является ли символ кавычкой (открывающей или закрывающей).
func IsQuote(ch rune) bool {
	ui := morphinternal.GetUnicodeInfo(ch)
	return ui.IsQuot()
}

// TransliteralCorrection пытается скорректировать транслитерацию между латиницей и кириллицей.
// Например, слово может быть написано в латинице с русскими буквами (или наоборот), и функция
// пытается заменить символы, основываясь на вероятности принадлежности к языку.
// Если `always` установлено в true, коррекция будет производиться даже без явных признаков ошибки.
func TransliteralCorrection(value, prevValue string, always bool) string {
	if value == "" {
		return value
	}

	pureCyr, pureLat := 0, 0
	quesCyr, quesLat := 0, 0
	udarCyr := 0
	hasY := false
	hasAccent := false

	for _, ch := range value {
		ui := morphinternal.GetUnicodeInfo(ch)
		if !ui.IsLetter() {
			if ui.IsUdaren() {
				hasAccent = true
				continue
			}
			if ui.IsApos() && len(value) > 2 {
				return TransliteralCorrection(strings.ReplaceAll(value, string(ch), ""), prevValue, false)
			}
			return value
		}
		if ui.IsCyrillic() {
			if strings.ContainsRune(LatChars, ch) {
				quesCyr++
			} else {
				pureCyr++
			}
		} else if ui.IsLatin() {
			if strings.ContainsRune(LatChars, ch) {
				quesLat++
			} else {
				pureLat++
			}
		} else if strings.ContainsRune(UdarChars, ch) {
			udarCyr++
		} else {
			return value
		}
		if ch == 'Ь' && strings.Contains(value, "ЬI") {
			hasY = true
		}
	}

	toRus, toLat := false, false
	if pureLat > 0 && pureCyr > 0 {
		return value
	}
	if (pureLat > 0 || always) && quesCyr > 0 {
		toLat = true
	} else if (pureCyr > 0 || always) && quesLat > 0 {
		toRus = true
	} else if pureCyr == 0 && pureLat == 0 {
		if quesCyr > 0 && quesLat > 0 {
			if len(prevValue) > 0 {
				if IsCyrillicChar(rune(prevValue[0])) {
					toRus = true
				} else if IsLatinChar(rune(prevValue[0])) {
					toLat = true
				}
			}
			if !toLat && !toRus {
				if quesCyr > quesLat {
					toRus = true
				} else if quesLat > quesCyr {
					toLat = true
				}
			}
		}
	}
	if !toRus && !toLat && !hasY && !hasAccent && udarCyr == 0 {
		return value
	}

	tmp := []rune(value)
	res := make([]rune, 0, len(tmp))
	for i := 0; i < len(tmp); i++ {
		ch := tmp[i]
		if ch == 'Ь' && i+1 < len(tmp) && tmp[i+1] == 'I' {
			res = append(res, 'Ы')
			i++
			continue
		}
		cod := int(ch)
		if cod >= 0x300 && cod < 0x370 {
			continue
		}
		if toRus {
			if idx := strings.IndexRune(LatChars, ch); idx >= 0 {
				res = append(res, rune(CyrChars[idx]))
			} else if idx := strings.IndexRune(UdarChars, ch); idx >= 0 {
				res = append(res, rune(UdarCyrChars[idx]))
			} else {
				res = append(res, ch)
			}
		} else if toLat {
			if idx := strings.IndexRune(CyrChars, ch); idx >= 0 {
				res = append(res, rune(LatChars[idx]))
			} else {
				res = append(res, ch)
			}
		} else {
			if idx := strings.IndexRune(UdarChars, ch); idx >= 0 {
				res = append(res, rune(UdarCyrChars[idx]))
			} else {
				res = append(res, ch)
			}
		}
	}
	return string(res)
}

// GetWordLang определяет язык слова по его символам.
// Возвращает Unknown, EN, или набор русскоязычных языков (RU, UA, BY, KZ).
func GetWordLang(word string) *MorphLang {
	cyr, lat, undef := 0, 0, 0

	for _, ch := range word {
		ui := morphinternal.GetUnicodeInfo(ch)
		if !ui.IsLetter() {
			continue
		}
		if ui.IsCyrillic() {
			cyr++
		} else if ui.IsLatin() {
			lat++
		} else {
			undef++
		}
	}

	if undef > 0 || (cyr == 0 && lat == 0) {
		return &MorphLangUnknown
	}
	if cyr == 0 {
		return &MorphLangEN
	}
	if lat > 0 {
		return &MorphLangUnknown
	}

	// Начинаем с полной кириллической группы
	lang := MorphLangRU.Or(MorphLangUA).Or(MorphLangBY).Or(MorphLangKZ)

	for _, ch := range word {
		switch ch {
		case 'Ґ', 'Є', 'Ї':
			lang = lang.Without(MorphLangRU, MorphLangBY)
		case 'І':
			lang = lang.Without(MorphLangRU)
		case 'Ё', 'Э':
			lang = lang.Without(MorphLangUA, MorphLangKZ)
		case 'Ы':
			lang = lang.Without(MorphLangUA)
		case 'Ў':
			lang = lang.Without(MorphLangRU, MorphLangUA)
		case 'Щ':
			lang = lang.Without(MorphLangBY)
		case 'Ъ':
			lang = lang.Without(MorphLangBY, MorphLangUA, MorphLangKZ)
		case 'Ә', 'Ғ', 'Қ', 'Ң', 'Ө', 'Ұ', 'Ү', 'Һ':
			lang = lang.Without(MorphLangBY, MorphLangUA, MorphLangRU)
		case 'В', 'Ф', 'Ц', 'Ч', 'Ь':
			lang = lang.Without(MorphLangKZ)
		}
	}

	return &lang
}

// EndsWith проверяет, заканчивается ли строка подстрокой (без учёта регистра).
func EndsWith(str, substr string) bool {
	if str == "" || substr == "" {
		return false
	}
	if len(substr) > len(str) {
		return false
	}
	return strings.HasSuffix(strings.ToLower(str), strings.ToLower(substr))
}

// EndsWithEx проверяет, заканчивается ли строка на одну из 2–4 подстрок.
func EndsWithEx(str, substr1, substr2, substr3, substr4 string) bool {
	if str == "" {
		return false
	}
	subs := []string{substr1, substr2, substr3, substr4}
	for _, s := range subs {
		if s == "" {
			continue
		}
		if strings.HasSuffix(str, s) {
			return true
		}
	}
	return false
}

// IsLatin возвращает true, если строка состоит из латинских букв, пробелов или дефисов.
func IsLatin(str string) bool {
	if str == "" {
		return false
	}
	for _, ch := range str {
		if !IsLatinChar(ch) && ch != '-' && ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			return false
		}
	}
	return true
}

// NormalizePreposition нормализует вариант предлога по словарю.
// Например, "ВО" → "В"
func NormalizePreposition(prep string) string {
	if norm, ok := PrepNorms[strings.ToUpper(prep)]; ok {
		return norm
	}
	return prep
}

// GetCaseAfterPreposition возвращает падеж, который требует данный предлог.
// Если предлог не найден — возвращается MorphCaseUndefined.
func GetCaseAfterPreposition(prep string) MorphCase {
	if mc, ok := PrepCases[strings.ToUpper(prep)]; ok {
		return mc
	}
	return Undefined
}

// ToStringMorphTense возвращает строковое представление времени (прошедшее, настоящее, будущее).
// Пример: MorphTensePast | MorphTenseFuture => "прошедшее|будущее"
func ToStringMorphTense(tense MorphTense) string {
	var res []string
	if tense.And(MorphTensePast).IsDefined() {
		res = append(res, "прошедшее")
	}
	if tense.And(MorphTensePresent).IsDefined() {
		res = append(res, "настоящее")
	}
	if tense.And(MorphTenseFuture).IsDefined() {
		res = append(res, "будущее")
	}
	return strings.Join(res, "|")
}

// ToStringMorphPerson возвращает строковое представление лица глагола (1, 2, 3).
// Пример: MorphPersonFirst | MorphPersonThird => "1лицо|3лицо"
func ToStringMorphPerson(person MorphPerson) string {
	var res []string
	if person.And(MorphPersonFirst).IsDefined() {
		res = append(res, "1лицо")
	}
	if person.And(MorphPersonSecond).IsDefined() {
		res = append(res, "2лицо")
	}
	if person.And(MorphPersonThird).IsDefined() {
		res = append(res, "3лицо")
	}
	return strings.Join(res, "|")
}

// ToStringMorphGender возвращает строку для грамматического рода (мужской, женский, средний).
func ToStringMorphGender(g MorphGender) string {
	var res []string
	if g.And(MorphGenderMasculine).IsDefined() {
		res = append(res, "муж.")
	}
	if g.And(MorphGenderFeminine).IsDefined() {
		res = append(res, "жен.")
	}
	if g.And(MorphGenderNeuter).IsDefined() {
		res = append(res, "средн.")
	}
	return strings.Join(res, "|")
}

// ToStringMorphNumber возвращает число: единственное, множественное или оба.
func ToStringMorphNumber(n MorphNumber) string {
	var res []string
	if n.And(MorphNumberSingular).IsDefined() {
		res = append(res, "единств.")
	}
	if n.And(MorphNumberPlural).IsDefined() {
		res = append(res, "множеств.")
	}
	return strings.Join(res, "|")
}

// ToStringMorphVoice возвращает вид глагола: действительный, страдательный, средний (рефлексив).
func ToStringMorphVoice(v MorphVoice) string {
	var res []string
	if v.And(MorphVoiceActive).IsDefined() {
		res = append(res, "действит.")
	}
	if v.And(MorphVoicePassive).IsDefined() {
		res = append(res, "страдат.")
	}
	if v.And(MorphVoiceMiddle).IsDefined() {
		res = append(res, "средн.")
	}
	return strings.Join(res, "|")
}

// ToStringMorphMood возвращает наклонение глагола: изъявительное, повелительное, условное.
func ToStringMorphMood(m MorphMood) string {
	var res []string
	if m.And(MorphMoodIndicative).IsDefined() {
		res = append(res, "изъявит.")
	}
	if m.And(MorphMoodImperative).IsDefined() {
		res = append(res, "повелит.")
	}
	if m.And(MorphMoodSubjunctive).IsDefined() {
		res = append(res, "условн.")
	}
	return strings.Join(res, "|")
}

// ToStringMorphAspect возвращает аспект глагола: совершенный или несовершенный.
func ToStringMorphAspect(a MorphAspect) string {
	var res []string
	if a.And(MorphAspectImperfective).IsDefined() {
		res = append(res, "несоверш.")
	}
	if a.And(MorphAspectPerfective).IsDefined() {
		res = append(res, "соверш.")
	}
	return strings.Join(res, "|")
}

// ToStringMorphFinite возвращает тип глагольной формы: личная, деепричастие, инфинитив, причастие.
func ToStringMorphFinite(f MorphFinite) string {
	var res []string
	if f.And(MorphFiniteFinite).IsDefined() {
		res = append(res, "finite")
	}
	if f.And(MorphFiniteGerund).IsDefined() {
		res = append(res, "gerund")
	}
	if f.And(MorphFiniteInfinitive).IsDefined() {
		res = append(res, "инфинитив")
	}
	if f.And(MorphFiniteParticiple).IsDefined() {
		res = append(res, "participle")
	}
	return strings.Join(res, "|")
}

// ToStringMorphForm возвращает форму слова: краткая, синонимичная и т.п.
func ToStringMorphForm(f MorphForm) string {
	var res []string
	if f.And(MorphFormShort).IsDefined() {
		res = append(res, "кратк.")
	}
	if f.And(MorphFormSynonym).IsDefined() {
		res = append(res, "синонимич.")
	}
	return strings.Join(res, "|")
}
