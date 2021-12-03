package harfbuzzgoperf

import (
	"encoding/binary"
	"unicode"

	hblang "github.com/benoitkugler/textlayout/language"
	"golang.org/x/text/language"
)

// --- Type conversion -------------------------------------------------------

// Lang4HB returns a language tag as a HarfBuzz language.
func Lang4HB(l language.Tag) hblang.Language {
	return hblang.NewLanguage(l.String())
}

// Lang4HB returns a script as a HarfBuzz script.
func Script4HB(s language.Script) hblang.Script {
	b := []byte(s.String())
	b[0] = byte(unicode.ToLower(rune(b[0])))
	h := binary.BigEndian.Uint32(b)
	return hblang.Script(h)
}
