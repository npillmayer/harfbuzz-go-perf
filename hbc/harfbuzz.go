package hbc

import (
	"encoding/binary"
	"unicode"

	"github.com/npillmayer/harfbuzzgoperf"
	"golang.org/x/text/language"
)

// https://harfbuzz.github.io/shaping-and-shape-plans.html

// Harfbuzz is the de-facto standard for text shaping.
// For further information see
// https://www.freedesktop.org/wiki/Software/HarfBuzz .
//
// A remark to the use of pointers to Harfbuzz-objects: Harfbuzz
// does its own memory management and we must avoid interfering with it.
// The Go garbage collector will therefore be unaware of the memory
// managed by Harfbuzz (in the worst of cases, a fancy Go garbage collector
// may re-locate memory). To hide Harfbuzz-memory from Go, we will use
// 'uintptr' variables instead of 'unsafe.Pointer's.
//
// The downside of this is the need to free() memory whenever we
// hand a Harfbuzz-shaper to GC.
type Harfbuzz struct {
	buffer    uintptr         // central data structure for Harfbuzz
	direction Direction       // L-to-R, R-to-L, T-to-B
	script    language.Script // i.e., Latin, Arabic, Korean, ...
}

// Direction is the direction to typeset text in.
type Direction int

const (
	LeftToRight Direction = iota
	RightToLeft           = 1
	TopToBottom           = 2
	BottomToTop           = 3
)

// Lang4HB returns a script as a HarfBuzz script.
func Script4HB(s language.Script) uint32 {
	b := []byte(s.String())
	b[0] = byte(unicode.ToLower(rune(b[0])))
	h := binary.BigEndian.Uint32(b)
	return h
}

// NewHarfbuzz creates a new Harfbuzz text shaper, fully initialized.
// Defaults are for Latin script, left-to-right.
func NewHarfbuzz() *Harfbuzz {
	hb := &Harfbuzz{}
	hb.buffer = allocHBBuffer()
	hb.direction = LeftToRight
	setHBBufferDirection(hb.buffer, hb.direction)
	hb.script = language.MustParseScript("Latn")
	setHBBufferScript(hb.buffer, Script4HB(hb.script))
	return hb
}

// Cache for font structures prepared for Harfbuzz.
// Harfbuzz uses its own font structure, different from ours.
// Unfortunately this duplicates the binary data of the font.
/*
var harfbuzzFontCache map[*font.TypeCase]uintptr
*/

// TODO: make cache thread-safe
// TODO: return error
/*
func (hb *Harfbuzz) findFont(typecase *font.TypeCase) uintptr {
	var hbfont uintptr
	if harfbuzzFontCache == nil {
		harfbuzzFontCache = make(map[*font.TypeCase]uintptr)
	}
	if hbfont = harfbuzzFontCache[typecase]; hbfont == 0 {
		if hbfont = makeHBFont(typecase); hbfont != 0 {
			harfbuzzFontCache[typecase] = hbfont
		}
	}
	return hbfont
}
*/

// SetScript is part of TextShaper interface.
func (hb *Harfbuzz) SetScript(scr language.Script) {
	setHBBufferScript(hb.buffer, Script4HB(scr))
}

// SetDirection is part of TextShaper interface.
func (hb *Harfbuzz) SetDirection(dir Direction) {
	setHBBufferDirection(hb.buffer, dir)
}

// SetLanguage is part of interface TextShaper.
// Harfbuzz doesn't evaluate a language parameter; method is a NOP.
func (hb *Harfbuzz) SetLanguage(string) {
}

// Shape is part of the  TextShaper interface.
//
// This is where all the heavy lifting is done. We input a font and a
// string of Unicode code-points, and receive a list of glyphs.
func (hb *Harfbuzz) Shape(text string, hbfont uintptr) *HBGlyphSequence {
	if hbfont == 0 {
		panic("cannot find Harfbuzz font")
	}
	harfbuzzShape(hb.buffer, text, hbfont)
	seq := getHBGlyphInfo(hb.buffer)
	return seq
}

func (hb *Harfbuzz) GlyphSequenceString(hbfont *harfbuzzgoperf.HBFont, seq *HBGlyphSequence) string {
	s := hbGlyphString(hbfont.CFont, seq)
	return s
}
