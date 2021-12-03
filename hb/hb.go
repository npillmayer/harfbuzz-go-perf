package hb

import (
	"errors"
	"fmt"

	hb "github.com/benoitkugler/textlayout/harfbuzz"
	"github.com/npillmayer/harfbuzzgoperf"
	"github.com/npillmayer/schuko/tracing"
	"golang.org/x/text/language"
)

// tracer traces to tracing key 'hbperf.hb'.
func tracer() tracing.Trace {
	return tracing.Select("hbperf.hb")
}

type HBParams struct {
	Font      *harfbuzzgoperf.HBFont // use a font at a given point-size
	PtSize    float32                // point size (1/72.27 in) of font to use
	Direction hb.Direction           // writing direction
	Script    language.Script        // 4-letter ISO 15924 script identifier
	Language  language.Tag           // BCP 47 language tag
	Features  []hb.Feature           // OpenType features to apply
}

func GetHBParams(fontname string, ptsize float32) (*HBParams, error) {
	p := &HBParams{}
	p.Font = harfbuzzgoperf.GlobalFontStore.FindFont(fontname)
	if p.Font == nil {
		return nil, fmt.Errorf("cannot find font %q", fontname)
	}
	p.PtSize = ptsize
	tracer().Infof("preparing font %q at %.2fpt", fontname, ptsize)
	return p, nil
}

// Shape calls the HarfBuzz shaper.
//
// Shape shapes a sequence of code-points (runes), turning its Unicode characters to
// positioned glyphs. It will select a shape plan based on params, including the
// selected font, and the properties of the input text.
//
// If `params.Features` is not empty, it will be used to control the
// features applied during shaping. If two features have the same tag but
// overlapping ranges the value of the feature with the higher index takes
// precedence.
//
// params.Font must be set, otherwise no output is created.
//
func Shape(text []rune, buf *hb.Buffer, params *HBParams) (*hb.Buffer, error) {
	if len(text) == 0 || params.Font == nil {
		return buf, errors.New("no input to shape")
	}
	// Prepare HarfBuzz buffer
	if buf == nil {
		buf = hb.NewBuffer()
		//buf.Props = hb_seqProps
	} else {
		buf.Clear()
	}
	buf.AddRunes(text, 0, len(text))
	buf.Shape(params.Font.GoFont, params.Features)
	// Prepare shaped output
	if len(buf.Info) == 0 {
		return buf, fmt.Errorf("nothing got shaped")
	}
	return buf, nil
}
