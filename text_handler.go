package shandler

import "time"

const (
	textComponentSep = '='
	textAttrSep      = ' '
)

type textBuilder struct {
	*baseBuilder
}

func (b *textBuilder) start() {}
func (b *textBuilder) end()   {}

func (b *textBuilder) componentSep() byte {
	return textComponentSep
}

func (b *textBuilder) attrSep() byte {
	return textAttrSep
}

func (b *textBuilder) appendTime(t time.Time) {
	if t.IsZero() {
		return
	}

	//b.buf.WriteString(timeStyle.Render(t.Format(b.timeFormat)))
	b.buf.WriteByte(b.componentSep())
}
