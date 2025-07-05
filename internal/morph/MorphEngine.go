package morph

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/jhonroun/pullenti/internal/morphinternal"
)

//go:embed ../internal/data/*.dat
var embeddedFS embed.FS

type MorphEngine struct {
	mu           sync.Mutex
	mLazyBuf     *morphinternal.ByteArrayWrapper
	Language     MorphLang
	MRoot        *MorphTreeNode
	MRootReverce *MorphTreeNode
	mRules       []*MorphRule
	mMiscInfos   []*morph.MorphMiscInfo
}

func NewMorphEngine() *MorphEngine {
	return &MorphEngine{
		MRoot:        &MorphTreeNode{},
		MRootReverce: &MorphTreeNode{},
		Language:     morph.MorphLangUndefined,
	}
}

func (me *MorphEngine) GetLazyBuf() *morphinternal.ByteArrayWrapper {
	return me.mLazyBuf
}

func (me *MorphEngine) AddRule(r *MorphRule) {
	me.mRules = append(me.mRules, r)
}

func (me *MorphEngine) GetRule(id int) *MorphRule {
	if id > 0 && id <= len(me.mRules) {
		return me.mRules[id-1]
	}
	return nil
}

func (me *MorphEngine) GetMutRule(id int) *MorphRule {
	return me.GetRule(id)
}

func (me *MorphEngine) GetRuleVar(rid, vid int) *MorphRuleVariant {
	r := me.GetRule(rid)
	if r == nil {
		return nil
	}
	return r.FindVar(vid)
}

func (me *MorphEngine) AddMiscInfo(mi *morph.MorphMiscInfo) {
	if mi.Id == 0 {
		mi.Id = int16(len(me.mMiscInfos) + 1)
	}
	me.mMiscInfos = append(me.mMiscInfos, mi)
}

func (me *MorphEngine) GetMiscInfo(id int) *morph.MorphMiscInfo {
	if id > 0 && id <= len(me.mMiscInfos) {
		return me.mMiscInfos[id-1]
	}
	return nil
}

func (me *MorphEngine) InitializeFromEmbedded(lang morph.MorphLang, lazy bool) bool {
	if !me.Language.IsUndefined() {
		return false
	}

	me.mu.Lock()
	defer me.mu.Unlock()

	if !me.Language.IsUndefined() {
		return false
	}

	me.Language = lang
	filename := fmt.Sprintf("m_%s.dat", lang.String())
	path := filepath.Join("../../internal/data", filename)

	data, err := embeddedFS.ReadFile(path)
	if err != nil {
		me.Language = morph.MorphLangUndefined // сбрасываем при неудаче
		return false
	}

	return me.Deserialize(bytes.NewReader(data), false, lazy)
}

// Deserialize загружает и инициализирует словарь морфологии
func (me *MorphEngine) Deserialize(stream io.Reader, ignoreRevTree, lazy bool) bool {
	if stream == nil {
		return false
	}

	// Декодируем сжатый GZIP-данный поток
	var rawData bytes.Buffer
	err := DeflateGzip(stream, &rawData)
	if err != nil {
		return false
	}

	// Оборачиваем в ByteArrayWrapper
	buf := rawData.Bytes()
	me.mLazyBuf = NewByteArrayWrapper(buf)
	pos := 0

	// Десериализация основного дерева
	if lazy {
		me.MRoot.DeserializeLazy(me.mLazyBuf, me, &pos)
	} else {
		me.MRoot.Deserialize(me.mLazyBuf, &pos)
	}

	if !ignoreRevTree {
		if lazy {
			me.MRootReverce.DeserializeLazy(me.mLazyBuf, me, &pos)
		} else {
			me.MRootReverce.Deserialize(me.mLazyBuf, &pos)
		}
	}

	return true
}
