package shandler

import (
	"fmt"
	"golang.org/x/exp/slog"
	"runtime"
	"strings"
)

const (
	textComponentSep = '='
	textAttrSep      = ' '
	callerSep        = "/"
)

type TextHandler struct {
	*baseHandler
}

func NewTextHandler(opts ...Option) *TextHandler {
	return &TextHandler{createHandler(false, opts...)}
}

func (t *TextHandler) WithPrefix(prefix string) slog.Handler {
	return &TextHandler{t.withPrefix(prefix)}
}

func (t *TextHandler) WithThemes(themes map[ThemeSection]*Theme) slog.Handler {
	return &TextHandler{t.withThemes(themes)}
}

type textBuilder struct {
	*baseBuilder
}

func (b *textBuilder) start() {}
func (b *textBuilder) close() {}

func (b *textBuilder) appendTime() {
	if b.r.Time.IsZero() {
		return
	}

	b.buf.WriteString(b.h.safeRender(b.h.themes[ThemeTime], b.r.Time.Format(b.h.timeFormat)))
}

func (b *textBuilder) appendLevel() {
	if !b.r.Time.IsZero() {
		b.buf.WriteByte(textAttrSep)
	}
	var level string
	var section ThemeSection
	switch {
	case b.r.Level < slog.LevelInfo:
		section = ThemeDebug
		level = "DBUG"
	case b.r.Level < slog.LevelWarn:
		section = ThemeInfo
		level = "INFO"
	case b.r.Level < slog.LevelError:
		section = ThemeWarn
		level = "WARN"
	default:
		section = ThemeError
		level = "ERRO"
	}
	b.buf.WriteString(b.h.safeRender(b.h.themes[section], level))
}

func (b *textBuilder) appendCaller() {
	if !b.h.caller || b.r.PC <= 0 {
		return
	}

	b.buf.WriteByte(textAttrSep)
	fs := runtime.CallersFrames([]uintptr{b.r.PC})
	f, _ := fs.Next()
	caller := f.Function
	if !b.h.fullCaller {
		paths := strings.Split(caller, callerSep)
		if len(paths) > 1 {
			paths = paths[len(paths)-2:]
			caller = strings.Join(paths, callerSep)
		}
	}
	caller = fmt.Sprintf("<%s:%d>", caller, f.Line)
	b.buf.WriteString(b.h.safeRender(b.h.themes[ThemeCaller], caller))
}

func (b *textBuilder) appendPrefix() {
	if b.h.prefix == "" {
		return
	}

	b.buf.WriteByte(textAttrSep)
	prefix := "[" + b.h.prefix + "]:"
	b.buf.WriteString(b.h.safeRender(b.h.themes[ThemePrefix], prefix))
}

func (b *textBuilder) appendMessage() {
	if b.r.Message == "" {
		return
	}

	b.buf.WriteByte(textAttrSep)
	b.buf.WriteString(b.r.Message)
}

func (b *textBuilder) appendAttrs() {
	if len(b.h.preformatted) > 0 {
		b.buf.WriteByte(textAttrSep)
		_, _ = b.buf.Write(b.h.preformatted)
	}

	b.prefix = NewBuffer()
	defer b.prefix.Free()
	b.prefix.WriteString(b.h.groupPrefix)
	for _, name := range b.h.groups[b.h.nOpenGroups:] {
		b.openGroup(name)
	}
	b.r.Attrs(func(a slog.Attr) bool {
		b.appendAttr(a)
		return true
	})
}

func (b *textBuilder) openGroup(name string) {
	b.prefix.WriteString(name)
	b.prefix.WriteByte(groupKeySep)
	if b.groups != nil {
		*b.groups = append(*b.groups, name)
	}
}

func (b *textBuilder) closeGroup(name string) {
	*b.prefix = (*b.prefix)[:len(*b.prefix)-len(name)-1]
	if b.groups != nil {
		*b.groups = (*b.groups)[:len(*b.groups)-1]
	}
}

func (b *textBuilder) appendAttr(a slog.Attr) {
	a = b.resolve(a)
	if a.Equal(slog.Attr{}) {
		return
	}

	if a.Value.Kind() != slog.KindGroup {
		b.buf.WriteByte(textAttrSep)
		b.buf.WriteString(b.h.safeRender(b.h.themes[ThemeKey], b.quote(string(*b.prefix)+a.Key)))
		b.buf.WriteByte(textComponentSep)
		b.appendValue(a.Value)
		return
	}

	if attrs := a.Value.Group(); len(attrs) > 0 {
		if a.Key != "" {
			b.openGroup(a.Key)
		}
		for _, attr := range attrs {
			b.appendAttr(attr)
		}
		if a.Key != "" {
			b.closeGroup(a.Key)
		}
	}
}

func (b *textBuilder) appendValue(v slog.Value) {
	b.buf.WriteString(b.quote(v.String()))
}

func (b *textBuilder) output() *Buffer {
	b.buf.WriteByte('\n')
	return b.buf
}
