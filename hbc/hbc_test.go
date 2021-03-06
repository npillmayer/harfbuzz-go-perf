package hbc_test

import (
	"testing"

	"github.com/npillmayer/harfbuzzgoperf"
	"github.com/npillmayer/harfbuzzgoperf/hbc"
)

func TestHarfbuzzShape(t *testing.T) {
	var hb *hbc.Harfbuzz
	if hb = hbc.NewHarfbuzz(nil); hb == nil {
		t.Error("cannot create Harfbuzz")
	}
	harfbuzzgoperf.LoadEmbeddedFonts()
	font := harfbuzzgoperf.GlobalFontStore.FindFont("Go")
	if font == nil {
		t.Fatal("expected to find font Go Sans")
	}
	font.CFont = hbc.MakeHBFont(font.Binary, 12.0)
	seq := hb.Shape("Wäffle", font.CFont)
	if seq == nil {
		t.Fail()
	}
	t.Logf("Input is \"Wäffle\",\nHarfbuzz returns %d glyphs: %s\n",
		seq.GlyphCount(), hb.GlyphSequenceString(font, seq))
	if seq.GlyphCount() != 6 {
		t.Errorf("expected to have 6 shaped glyphs, have %d", seq.GlyphCount())
	}
}

func TestHarfbuzzShapeResult(t *testing.T) {
	var seq *hbc.HBGlyphSequence
	harfbuzzgoperf.LoadEmbeddedFonts()
	font := harfbuzzgoperf.GlobalFontStore.FindFont("Calibri.ttf")
	if font == nil {
		t.Fatal("expected to find font Calibri")
	}
	font.CFont = hbc.MakeHBFont(font.Binary, 12.0)
	var hb *hbc.Harfbuzz
	if hb = hbc.NewHarfbuzz(nil); hb == nil {
		t.Fatal("failed to create Harfbuzz instance")
	}
	if seq = hb.Shape("Fifig", font.CFont); seq != nil {
		t.Logf("Input is \"Fifig\",\nHarfbuzz returns %d glyphs: %s\n",
			seq.GlyphCount(), hb.GlyphSequenceString(font, seq))
		cnt := seq.GlyphCount()
		for i := 0; i < cnt; i++ {
			gi := seq.GetGlyphInfoAt(i)
			t.Logf("glyph info #%d/%d: x-advance %.2f\n", i, gi.Cluster(), gi.XAdvance())
		}
	}
	if seq == nil {
		t.Error("expected shaping output to be non-nil")
	}
}

// --- Benchmarking ----------------------------------------------------------

var Cnt int

func BenchmarkHBShape(b *testing.B) {
	var seq *hbc.HBGlyphSequence
	harfbuzzgoperf.LoadEmbeddedFonts()
	fontname := "Calibri.ttf"
	//fontname := "Go"
	font := harfbuzzgoperf.GlobalFontStore.FindFont(fontname)
	if font == nil {
		b.Fatalf("cannot prepare font %q", fontname)
	}
	font.CFont = hbc.MakeHBFont(font.Binary, 12.0)
	buf := hbc.AllocHBBuffer()
	var hb *hbc.Harfbuzz
	for i := 0; i < b.N; i++ {
		for _, line := range harfbuzzgoperf.Corpus {
			if hb = hbc.NewHarfbuzz(buf); hb == nil {
				b.Fatal("failed to create Harfbuzz instance")
			}
			seq = hb.Shape(line, font.CFont)
			Cnt = seq.GlyphCount()
		}
	}
}
