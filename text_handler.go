package shandler

import (
	"golang.org/x/exp/slog"
	"runtime"
	"strconv"
	"time"
	"unicode/utf8"
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

func (t *TextHandler) WithThemes(themes map[ThemeSchema]*Theme) slog.Handler {
	return &TextHandler{t.withThemes(themes)}
}

type textBuilder struct {
	*baseBuilder
}

func (b *textBuilder) start() {}
func (b *textBuilder) close() {}

// appendTime If r.Time is the zero time, ignore the time.
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
	var section ThemeSchema
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

// appendCaller If r.PC is zero or disabled caller, ignore it.
func (b *textBuilder) appendCaller() {
	if !b.h.caller || b.r.PC <= 0 {
		return
	}

	b.buf.WriteByte(textAttrSep)
	fs := runtime.CallersFrames([]uintptr{b.r.PC})
	f, _ := fs.Next()
	caller := f.Function
	if !b.h.fullCaller {
		var idx, founded int
		for i := utf8.RuneCountInString(caller) - 1; i >= 0; i-- {
			if caller[i] == callerSep[0] {
				founded++
			}
			idx = i
			if founded == 2 {
				break
			}
		}
		if idx == 0 {
			idx -= 1
		}
		caller = caller[idx+1:]
	}
	caller = "<" + caller + ":" + strconv.Itoa(f.Line) + ">"
	b.buf.WriteString(b.h.safeRender(b.h.themes[ThemeCaller], caller))
}

func (b *textBuilder) appendPrefix() {
	var prefix string
	if b.h.prefix == "" {
		prefix = "\uF444"
	} else {
		prefix = "[" + b.h.prefix + "]:"
	}

	b.buf.WriteByte(textAttrSep)
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

// appendAttr If an Attr's key and value are both the zero value, ignore the Attr.
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
	switch v.Kind() {
	case slog.KindString:
		b.buf.WriteString(b.quote(v.String()))
	case slog.KindInt64:
		b.buf.WriteString(strconv.FormatInt(v.Int64(), 10))
	case slog.KindUint64:
		b.buf.WriteString(strconv.FormatUint(v.Uint64(), 10))
	case slog.KindFloat64:
		b.buf.WriteString(strconv.FormatFloat(v.Float64(), 'f', 10, 64))
	case slog.KindBool:
		b.buf.WriteString(strconv.FormatBool(v.Bool()))
	case slog.KindTime:
		b.buf.WriteString(v.Time().Format(time.RFC3339Nano))
	case slog.KindDuration:
		b.buf.WriteString(v.Duration().String())
	case slog.KindAny:
	default:
		b.buf.WriteString(b.quote(v.String()))
	}
}

func (b *textBuilder) output() *Buffer {
	b.buf.WriteByte('\n')
	return b.buf
}
