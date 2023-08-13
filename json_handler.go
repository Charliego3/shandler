package shandler

import (
	"golang.org/x/exp/slog"
	"time"
)

const (
	jsonComponentSep = ':'
	jsonAttrSep      = ','
)

type JsonHandler struct {
	*baseHandler
}

func NewJsonHandler(opts ...Option) *JsonHandler {
	return &JsonHandler{createHandler(true, opts...)}
}

func (j *JsonHandler) WithPrefix(prefix string) slog.Handler {
	h := j.clone()
	h.initThemes()
	h.prefix = prefix
	return &JsonHandler{h}
}

type jsonBuilder struct {
	*baseBuilder
}

func (b *jsonBuilder) start() {
	b.buf.WriteString(b.h.themes[ThemeBracket].Render("{"))
}

func (b *jsonBuilder) close() {
	b.buf.WriteString(b.h.themes[ThemeBracket].Render("}"))
}

func (b *jsonBuilder) componentSep() byte {
	return jsonComponentSep
}

func (b *jsonBuilder) attrSep() byte {
	return jsonAttrSep
}

func (b *jsonBuilder) appendTime(t time.Time) {
	if t.IsZero() {
		return
	}

	//b.buf.WriteString(timeTheme.Render(t.Format(b.opts.timeFormat)))
	b.buf.WriteByte(b.componentSep())
}

func (b *jsonBuilder) appendLevel(slog.Level) {

}
