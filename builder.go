package shandler

import (
	"sync"
	"time"
)

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

type Builder interface {
	start()
	end()
	appendTime(time.Time)
	componentSep() byte
	attrSep() byte
	free()
}

type baseBuilder struct {
	buf     *Buffer
	freeBuf bool      // should buf be freed?
	sep     string    // separator to write before next key
	prefix  *Buffer   // for text: key prefix
	groups  *[]string // pool-allocated slice of active groups, for ReplaceAttr
	json    bool
}

func (h *baseHandler) newBuilder(buf *Buffer, freeBuf bool, sep string, prefix *Buffer) Builder {
	builder := &baseBuilder{buf: buf, freeBuf: freeBuf, sep: sep, prefix: prefix}
	if h.replacer != nil {
		builder.groups = groupPool.Get().(*[]string)
		*builder.groups = append(*builder.groups, h.groups[:h.nOpenGroups]...)
	}
	if h.json {
		return &jsonBuilder{baseBuilder: builder}
	}
	return &textBuilder{baseBuilder: builder}
}

func (b *baseBuilder) free() {
	if b.freeBuf {
		b.buf.Free()
	}
	if gs := b.groups; gs != nil {
		*gs = (*gs)[:0]
		groupPool.Put(gs)
	}
}
