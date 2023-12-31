package shandler

import (
	"strconv"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"log/slog"
)

const groupKeySep = '.'

var groupPool = sync.Pool{New: func() any {
	s := make([]string, 0, 10)
	return &s
}}

type Builder interface {
	start()
	close()
	free()
	appendTime()
	appendLevel()
	appendCaller()
	appendPrefix()
	appendMessage()
	appendAttrs()
	output() *Buffer
}

type baseBuilder struct {
	h      *baseHandler
	r      slog.Record
	buf    *Buffer
	prefix *Buffer   // for text: key prefix
	groups *[]string // pool-allocated slice of active groups, for ReplaceAttr
}

func (h *baseHandler) createBaseBuilder(buf *Buffer, r slog.Record) *baseBuilder {
	b := &baseBuilder{h: h, r: r, buf: buf}
	if h.replacer != nil {
		b.groups = groupPool.Get().(*[]string)
		*b.groups = append(*b.groups, h.groups[:h.nOpenGroups]...)
	}
	return b
}

func (b *baseBuilder) resolve(a slog.Attr) slog.Attr {
	if rep := b.h.replacer; rep != nil && a.Value.Kind() != slog.KindGroup {
		var gs []string
		if b.groups != nil {
			gs = *b.groups
		}
		a.Value = a.Value.Resolve()
		a = rep(gs, a)
	}
	a.Value = a.Value.Resolve()
	return a
}

func (b *baseBuilder) quote(str string) string {
	if !needsQuoting(str) {
		return str
	}

	return strconv.Quote(str)
}

func needsQuoting(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); {
		b := s[i]
		if b < utf8.RuneSelf {
			if b == '\t' || b == '\n' {
				return false
			}
			if b != '\\' && (b == '=' || !safeSet[b]) {
				return true
			}
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError || unicode.IsSpace(r) || !unicode.IsPrint(r) {
			return true
		}
		i += size
	}
	return false
}

func (b *baseBuilder) free() {
	b.buf.Free()
	if gs := b.groups; gs != nil {
		*gs = (*gs)[:0]
		groupPool.Put(gs)
	}
}

func (b *baseBuilder) appendTime(t time.Time) {
	year, month, day := t.Date()
	b.buf.WritePosIntWidth(year, 4)
	b.buf.WriteByte('-')
	b.buf.WritePosIntWidth(int(month), 2)
	b.buf.WriteByte('-')
	b.buf.WritePosIntWidth(day, 2)
	b.buf.WriteByte('T')
	hour, minute, sec := t.Clock()
	b.buf.WritePosIntWidth(hour, 2)
	b.buf.WriteByte(':')
	b.buf.WritePosIntWidth(minute, 2)
	b.buf.WriteByte(':')
	b.buf.WritePosIntWidth(sec, 2)
	ns := t.Nanosecond()
	b.buf.WriteByte('.')
	b.buf.WritePosIntWidth(ns/1e6, 3)
	_, offsetSeconds := t.Zone()
	if offsetSeconds == 0 {
		b.buf.WriteByte('Z')
	} else {
		offsetMinutes := offsetSeconds / 60
		if offsetMinutes < 0 {
			b.buf.WriteByte('-')
			offsetMinutes = -offsetMinutes
		} else {
			b.buf.WriteByte('+')
		}
		b.buf.WritePosIntWidth(offsetMinutes/60, 2)
		b.buf.WriteByte(':')
		b.buf.WritePosIntWidth(offsetMinutes%60, 2)
	}
}

var safeSet = [utf8.RuneSelf]bool{
	' ':  true,
	'!':  true,
	'"':  false,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'(':  true,
	')':  true,
	'*':  true,
	'+':  true,
	',':  true,
	'-':  true,
	'.':  true,
	'/':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	':':  true,
	';':  true,
	'<':  true,
	'=':  true,
	'>':  true,
	'?':  true,
	'@':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'V':  true,
	'W':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'[':  true,
	'\\': false,
	']':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'{':  true,
	'|':  true,
	'}':  true,
	'~':  true,

	'\u007f': true,
}
