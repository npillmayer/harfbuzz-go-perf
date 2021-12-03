package harfbuzzgoperf

import (
	"bytes"
	"embed"
	"path/filepath"
	"sync"

	tt "github.com/benoitkugler/textlayout/fonts/truetype"
	hb "github.com/benoitkugler/textlayout/harfbuzz"
	"golang.org/x/image/font/gofont/goregular"
)

//go:embed resources/*
var resources embed.FS

type HBFont struct {
	Binary []byte
	GoFont *hb.Font
	CFont  uintptr
}

func LoadEmbeddedFonts() {
	fonts, _ := resources.ReadDir("resources/fonts")
	for _, font := range fonts {
		tracer().Debugf("found embedded font file %s", font.Name())
		f := &HBFont{}
		f.Binary, _ = resources.ReadFile(filepath.Join("resources", "fonts", font.Name()))
		fr := bytes.NewReader(f.Binary)
		hb_face, err := tt.Parse(fr, true)
		if err != nil {
			tracer().Errorf("cannot parse font %s: %s", font.Name(), err)
			continue
		}
		f.GoFont = hb.NewFont(hb_face)
		GlobalFontStore.StoreFont(font.Name(), f)
	}
	f := &HBFont{}
	f.Binary = goregular.TTF
	fr := bytes.NewReader(f.Binary)
	hb_face, _ := tt.Parse(fr, true)
	f.GoFont = hb.NewFont(hb_face)
	GlobalFontStore.StoreFont("Go", f)
}

// --- Font cache ------------------------------------------------------------

// FontCache is a super-simple-minded dictionary for caching fonts. This is just
// for facilitating the tests.
var GlobalFontStore = NewFontCache()
var fontCacheMx = &sync.RWMutex{}

type FontCache map[string]*HBFont

func NewFontCache() FontCache {
	return FontCache(make(map[string]*HBFont))
}

func (cache FontCache) FindFont(name string) *HBFont {
	fontCacheMx.RLock()
	defer fontCacheMx.RUnlock()
	if f, ok := cache[name]; !ok {
		return nil
	} else {
		return f
	}
}

func (cache FontCache) StoreFont(name string, font *HBFont) {
	fontCacheMx.Lock()
	defer fontCacheMx.Unlock()
	cache[name] = font
}
