package harfbuzzgoperf

import (
	"testing"

	"github.com/npillmayer/schuko/tracing/gotestingadapter"
)

func TestEmbeddedFonts(t *testing.T) {
	teardown := gotestingadapter.QuickConfig(t, "hbperf.base")
	defer teardown()
	//
	LoadEmbeddedFonts()
	if len(GlobalFontStore) != 3 {
		t.Fatalf("expected 3 fonts to be pre-loaded, have %d", len(GlobalFontStore))
	}
	f := GlobalFontStore.FindFont("Go")
	if f == nil {
		t.Fatal("expected to find font Go, could not")
	}
}
