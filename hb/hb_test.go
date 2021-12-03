package hb

import (
	"testing"

	hb "github.com/benoitkugler/textlayout/harfbuzz"
	"github.com/npillmayer/harfbuzzgoperf"
	"github.com/npillmayer/schuko/tracing/gotestingadapter"
)

func TestShapeSimple(t *testing.T) {
	teardown := gotestingadapter.QuickConfig(t, "hbperf.base")
	defer teardown()
	//
	harfbuzzgoperf.LoadEmbeddedFonts()
	fontname := "Go"
	text := []rune("The quick brown fox jumps over the lazy dog")
	params, err := GetHBParams(fontname, 12.0)
	if err != nil {
		t.Fatalf("cannot prepare font %q", fontname)
	}
	buf, err := Shape(text, nil, params)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("# of glyphs = %d", len(buf.Info))
}

var buf *hb.Buffer

func BenchmarkHBShape(b *testing.B) {
	harfbuzzgoperf.LoadEmbeddedFonts()
	fontname := "Calibri.ttf"
	params, err := GetHBParams(fontname, 12.0)
	if err != nil {
		b.Fatalf("cannot prepare font %q", fontname)
	}
	for i := 0; i < b.N; i++ {
		for _, line := range harfbuzzgoperf.Corpus {
			buf, err = Shape(line, nil, params)
			if err != nil {
				b.Fatal("expected shaping output to be non-nil")
			}
		}
	}
}
