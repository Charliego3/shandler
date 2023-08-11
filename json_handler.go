package shandler

import "time"

const (
	jsonComponentSep = ':'
	jsonAttrSep      = ','
)

type jsonBuilder struct {
	*baseBuilder
}

func (b *jsonBuilder) appendTime(t time.Time) {
	if t.IsZero() {
		return
	}

	//b.buf.WriteString(timeStyle.Render(t.Format(b.opts.timeFormat)))
	b.buf.WriteByte(b.componentSep())
}

func (b *jsonBuilder) componentSep() byte {
	return jsonComponentSep
}

func (b *jsonBuilder) attrSep() byte {
	return jsonAttrSep
}

func (b *jsonBuilder) start() {
	b.buf.WriteByte('{')
}

func (b *jsonBuilder) end() {
	b.buf.WriteByte('}')
}
