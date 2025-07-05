package morph

import (
	"errors"
	"strings"

	"github.com/jhonroun/pullenti/internal/MorphForm"
	"github.com/jhonroun/pullenti/internal/morphinternal"
)

var lazyLoad = true
var initialized = false
var morphs = InnerMorphology{}
var emptyWordForms = []*MorphForm.MorphWordForm{}
var emptyMisc = &MorphForm.MorphMiscInfo{}

// Initialize загружает словари для указанных языков.
// Если не вызвать, будет вызвано автоматически при первом использовании.
func Initialize(langs *MorphForm.MorphLang) {
	if initialized {
		return
	}
	morphinternal.Initialize()
	if langs == nil || langs.IsUndefined() {
		l := MorphForm.NewMorphLang()
		l = l.Or(MorphForm.MorphLangRU)
		l = l.Or(MorphForm.MorphLangEN)
		langs = &l
	}
	morphs.LoadLanguages(*langs, lazyLoad)
	initialized = true
}

// LoadedLanguages возвращает языки, для которых загружены словари.
func LoadedLanguages() MorphForm.MorphLang {
	return morphs.LoadedLanguages()
}

// LoadLanguages загружает указанные языки.
func LoadLanguages(langs MorphForm.MorphLang) {
	morphs.LoadLanguages(langs, lazyLoad)
}

// UnloadLanguages выгружает указанные языки.
func UnloadLanguages(langs MorphForm.MorphLang) {
	morphs.UnloadLanguages(langs)
}

// Tokenize выполняет токенизацию текста без морфологического анализа.
func Tokenize(text string) []*MorphForm.MorphToken {
	if text == "" {
		return nil
	}
	lang := MorphForm.NewMorphLang()
	res := morphs.Run(text, true, &lang, false, nil)
	for _, r := range res {
		if r.WordForms == nil {
			r.WordForms = emptyWordForms
		}
		for _, wf := range r.WordForms {
			if wf.Misc == nil {
				wf.Misc = emptyMisc
			}
		}
	}
	return res
}

// Process выполняет морфологический анализ текста.
func Process(text string, lang *MorphForm.MorphLang) ([]*MorphForm.MorphToken, error) {
	if text == "" {
		return nil, nil
	}
	if !initialized {
		return nil, errors.New("pullenti Morphology Service not initialized")
	}
	res := morphs.Run(text, false, lang, false, nil)
	for _, r := range res {
		if r.WordForms == nil {
			r.WordForms = emptyWordForms
		}
		for _, wf := range r.WordForms {
			if wf.Misc == nil {
				wf.Misc = emptyMisc
			}
		}
	}
	return res, nil
}

// GetAllWordforms возвращает все словоформы для слова.
func GetAllWordforms(word string, lang *MorphForm.MorphLang) []*MorphForm.MorphWordForm {
	if word == "" {
		return nil
	}
	if !initialized {
		panic("Pullenti Morphology Service not initialized")
	}
	if strings.ToLower(word) != word {
		word = strings.ToUpper(word)
	}
	res := morphs.GetAllWordforms(word, *lang)
	for _, wf := range res {
		if wf.Misc == nil {
			wf.Misc = emptyMisc
		}
	}
	return res
}

// GetAllWordsByClass возвращает все словоформы указанного класса.
func GetAllWordsByClass(cla *MorphForm.MorphClass, lang *MorphForm.MorphLang) []*MorphForm.MorphWordForm {
	if !initialized {
		panic("Pullenti Morphology Service not initialized")
	}
	return morphs.GetAllWordsByClass(*cla, *lang)
}
