package shandler

import (
	"io"
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/exp/slog"
)

// Replacer is called to rewrite each non-group attribute before it is logged.
// The attribute's value has been resolved (see [slog.Value.Resolve]).
// If Replacer returns an Attr with Key == "", the attribute is discarded.
//
// The built-in attributes with keys "time", "level", "source", and "msg"
// are passed to this function, except that time is omitted
// if zero, and source is omitted if AddSource is false.
//
// The first argument is a list of currently open groups that contain the
// Attr. It must not be retained or modified. Replacer is never called
// for Group attributes, only their contents. For example, the attribute
// list
//
//	Int("a", 1), Group("g", Int("b", 2)), Int("c", 3)
//
// results in consecutive calls to Replacer with the following arguments:
//
//	nil, Int("a", 1)
//	[]string{"g"}, Int("b", 2)
//	nil, Int("c", 3)
//
// Replacer can be used to change the default keys of the built-in
// attributes, convert types (for example, to replace a `time.Time` with the
// integer seconds since the Unix epoch), sanitize personal information, or
// remove attributes from the output.
type Replacer func(groups []string, a slog.Attr) slog.Attr

type Option func(*baseHandler)

func createHandler(json bool, opts ...Option) *baseHandler {
	h := &baseHandler{
		timeFormat: time.Kitchen,
		w:          os.Stderr,
		level:      slog.LevelInfo,
		json:       json,
		themes:     make(map[ThemeSchema]*Theme, 9),
	}
	for _, opt := range opts {
		opt(h)
	}
	h.initThemes()
	return h
}

func (h *baseHandler) initThemes() {
	if h.tty = h.isTTY(); !h.tty {
		return
	}

	h.themes[ThemeTime] = fillTheme(h.themes[ThemeTime], "#6085b9", "#7d467c", false, true, false)
	h.themes[ThemeDebug] = fillTheme(h.themes[ThemeDebug], "#4746ff", "#2f81ff", true, false, false)
	h.themes[ThemeInfo] = fillTheme(h.themes[ThemeInfo], "#009adc", "#00FFD5", true, false, false)
	h.themes[ThemeWarn] = fillTheme(h.themes[ThemeWarn], "#e16c00", "#ff9c01", true, false, false)
	h.themes[ThemeError] = fillTheme(h.themes[ThemeError], "#ff000a", "#FF4F86", true, false, false)
	h.themes[ThemePrefix] = fillTheme(h.themes[ThemePrefix], "#579159", "#008708", true, false, false)
	h.themes[ThemeCaller] = fillTheme(h.themes[ThemeCaller], "#765ea5", "#2f6e87", false, false, false)
	h.themes[ThemeKey] = fillTheme(h.themes[ThemeKey], "#7F7F7F", "#7F7F7F", true, false, false)
	if h.json {
		h.themes[ThemeBracket] = fillTheme(h.themes[ThemeBracket], "#000000", "#ffffff", true, false, false)
	}
}

func fillTheme(source *Theme, l, d string, bold, underline, overline bool) *Theme {
	if source == nil {
		light, _ := colorful.Hex(l)
		dark, _ := colorful.Hex(d)
		source = NewTheme().Foreground(light, dark)
		source.Bold(bold).Underline(underline).Overline(overline)
	}
	return source.Format()
}

func WithTimeFormat(format string) Option {
	return func(cfg *baseHandler) {
		cfg.timeFormat = format
	}
}

func WithWriter(w io.Writer) Option {
	return func(cfg *baseHandler) {
		cfg.w = w
	}
}

func WithLevel(level slog.Level) Option {
	return func(cfg *baseHandler) {
		cfg.level = level
	}
}

func WithPrefix(prefix string) Option {
	return func(cfg *baseHandler) {
		cfg.prefix = prefix
	}
}

// WithReplacer please refer to Replacer
func WithReplacer(fn Replacer) Option {
	return func(cfg *baseHandler) {
		cfg.replacer = fn
	}
}

func WithCaller() Option {
	return func(cfg *baseHandler) {
		cfg.caller = true
	}
}

func WithFullCaller() Option {
	return func(cfg *baseHandler) {
		cfg.fullCaller = true
	}
}

func WithTheme(section ThemeSchema, theme *Theme) Option {
	return func(cfg *baseHandler) {
		if theme == nil {
			return
		}
		cfg.themes[section] = theme.Format()
	}
}
