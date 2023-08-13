package shandler

import (
	"fmt"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"strings"
)

// Sequence definitions.
const (
	separator    = ";"
	ResetSeq     = "0"
	BoldSeq      = "1"
	FaintSeq     = "2"
	ItalicSeq    = "3"
	UnderlineSeq = "4"
	BlinkSeq     = "5"
	ReverseSeq   = "7"
	CrossOutSeq  = "9"
	OverlineSeq  = "53"
	Foreground   = "38"
	Background   = "48"

	ESC = '\x1b'
	CSI = string(ESC) + "["
)

//go:generate enumer -type ThemeSection -output theme_string.go
type ThemeSection uint

const (
	ThemeTime ThemeSection = iota + 1
	ThemeDebug
	ThemeInfo
	ThemeWarn
	ThemeError
	ThemePrefix
	ThemeCaller
	ThemeKey
	ThemeBracket // only json handler
)

var hasDarkBackground = termenv.HasDarkBackground()

type Theme struct {
	sequences []string
	formatted string
}

func NewTheme() *Theme {
	return &Theme{sequences: make([]string, 0, 10)}
}

// Foreground sets a foreground color.
func (t *Theme) Foreground(light, dark colorful.Color) *Theme {
	if hasDarkBackground {
		t.sequences = append(t.sequences, t.getSequence(dark, false))
	} else {
		t.sequences = append(t.sequences, t.getSequence(light, false))
	}
	return t
}

// Background sets a background color.
func (t *Theme) Background(light, dark colorful.Color) *Theme {
	if hasDarkBackground {
		t.sequences = append(t.sequences, t.getSequence(dark, true))
	} else {
		t.sequences = append(t.sequences, t.getSequence(light, true))
	}
	return t
}

// Bold enables bold rendering.
func (t *Theme) Bold() *Theme {
	t.sequences = append(t.sequences, BoldSeq)
	return t
}

// Faint enables faint rendering.
func (t *Theme) Faint() *Theme {
	t.sequences = append(t.sequences, FaintSeq)
	return t
}

// Italic enables italic rendering.
func (t *Theme) Italic() *Theme {
	t.sequences = append(t.sequences, ItalicSeq)
	return t
}

// Underline enables underline rendering.
func (t *Theme) Underline() *Theme {
	t.sequences = append(t.sequences, UnderlineSeq)
	return t
}

// Overline enables overline rendering.
func (t *Theme) Overline() *Theme {
	t.sequences = append(t.sequences, OverlineSeq)
	return t
}

// Blink enables blink mode.
func (t *Theme) Blink() *Theme {
	t.sequences = append(t.sequences, BlinkSeq)
	return t
}

// Reverse enables reverse color mode.
func (t *Theme) Reverse() *Theme {
	t.sequences = append(t.sequences, ReverseSeq)
	return t
}

// CrossOut enables crossed-out rendering.
func (t *Theme) CrossOut() *Theme {
	t.sequences = append(t.sequences, CrossOutSeq)
	return t
}

func (t *Theme) getSequence(f colorful.Color, bg bool) string {
	prefix := Foreground
	if bg {
		prefix = Background
	}
	r, g, b := f.RGB255()
	return fmt.Sprintf("%s;2;%d;%d;%d", prefix, r, g, b)
}

func (t *Theme) Format() *Theme {
	t.formatted = fmt.Sprintf("%s%sm", CSI, strings.Join(t.sequences, separator))
	return t
}

func (t *Theme) Render(s string) string {
	return fmt.Sprintf("%s%s%sm", t.formatted, s, CSI+ResetSeq)
}

func (h *baseHandler) safeRender(t *Theme, s string) string {
	if !h.tty || t == nil {
		return s
	}
	return t.Render(s)
}
