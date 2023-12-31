package shandler

import (
	"runtime"
	"strconv"
	"strings"

	"log/slog"
)

const (
	textComponentSep = '='
	textAttrSep      = ' '
	callerSep        = '/'
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

func (t *TextHandler) WithThemes(themes Themes) slog.Handler {
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

	b.h.WriteColorful(ThemeTime, b.buf, b.r.Time.Format(b.h.timeFormat))
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
	b.h.WriteColorful(section, b.buf, level)
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
		var founded int
		idx := strings.LastIndexFunc(caller, func(r rune) bool {
			if r == callerSep {
				founded++
			}
			if founded == 2 {
				return true
			}
			return false
		})
		caller = caller[idx+1:]
	}
	caller = "<" + caller + ":" + strconv.Itoa(f.Line) + ">"
	b.h.WriteColorful(ThemeCaller, b.buf, caller)
}

func (b *textBuilder) appendPrefix() {
	var prefix string
	if b.h.prefix == "" {
		prefix = ""
	} else {
		prefix = "[" + b.h.prefix + "]:"
	}

	b.buf.WriteByte(textAttrSep)
	b.h.WriteColorful(ThemePrefix, b.buf, prefix)
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
		b.h.WriteColorful(ThemeKey, b.buf, b.quote(string(*b.prefix)+a.Key))
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
		b.baseBuilder.appendTime(v.Time())
	case slog.KindDuration:
		b.buf.WriteString(v.Duration().String())
	default:
		b.buf.WriteString(b.quote(v.String()))
	}
}

func (b *textBuilder) output() *Buffer {
	b.buf.WriteByte('\n')
	return b.buf
}
