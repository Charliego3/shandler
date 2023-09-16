package shandler

import (
	"fmt"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
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

//go:generate enumer -type ThemeSchema -output theme_string.go
type ThemeSchema uint

type Themes map[ThemeSchema]*Theme

const (
	ThemeTime ThemeSchema = iota + 1
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
func (t *Theme) Bold(bold ...bool) *Theme {
	if len(bold) > 0 && !bold[0] {
		return t
	}
	t.sequences = append(t.sequences, BoldSeq)
	return t
}

// Faint enables faint rendering.
func (t *Theme) Faint(faint ...bool) *Theme {
	if len(faint) > 0 && !faint[0] {
		return t
	}
	t.sequences = append(t.sequences, FaintSeq)
	return t
}

// Italic enables italic rendering.
func (t *Theme) Italic(italic ...bool) *Theme {
	if len(italic) > 0 && !italic[0] {
		return t
	}
	t.sequences = append(t.sequences, ItalicSeq)
	return t
}

// Underline enables underline rendering.
func (t *Theme) Underline(underline ...bool) *Theme {
	if len(underline) > 0 && !underline[0] {
		return t
	}
	t.sequences = append(t.sequences, UnderlineSeq)
	return t
}

// Overline enables overline rendering.
func (t *Theme) Overline(overline ...bool) *Theme {
	if len(overline) > 0 && !overline[0] {
		return t
	}
	t.sequences = append(t.sequences, OverlineSeq)
	return t
}

// Blink enables blink mode.
func (t *Theme) Blink(blink ...bool) *Theme {
	if len(blink) > 0 && !blink[0] {
		return t
	}
	t.sequences = append(t.sequences, BlinkSeq)
	return t
}

// Reverse enables reverse color mode.
func (t *Theme) Reverse(reverse ...bool) *Theme {
	if len(reverse) > 0 && !reverse[0] {
		return t
	}
	t.sequences = append(t.sequences, ReverseSeq)
	return t
}

// CrossOut enables crossed-out rendering.
func (t *Theme) CrossOut(crossOut ...bool) *Theme {
	if len(crossOut) > 0 && !crossOut[0] {
		return t
	}
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
	return t.formatted + s + CSI + ResetSeq + "m"
}

func (h *baseHandler) safeRender(schema ThemeSchema, s string) string {
	if !h.tty {
		return s
	}
	if theme, ok := h.themes[schema]; ok {
		return theme.Render(s)
	}
	return s
}
