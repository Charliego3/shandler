package shandler

import (
	"context"
	"io"
	"sync"

	"github.com/mattn/go-isatty"
	"log/slog"
	"slices"
)

type Handler interface {
	WithPrefix(prefix string) slog.Handler
	WithThemes(themes Themes) slog.Handler
}

func getHandler() Handler {
	h := slog.Default().Handler()
	if sh, ok := h.(Handler); ok {
		return sh
	}
	return nil
}

func CopyWithPrefix(prefix string) *slog.Logger {
	h := getHandler()
	if h == nil {
		return nil
	}
	return slog.New(h.WithPrefix(prefix))
}

func CopyWithThemes(themes Themes) *slog.Logger {
	h := getHandler()
	if h == nil {
		return nil
	}
	return slog.New(h.WithThemes(themes))
}

// File represents a file descriptor.
type File interface {
	io.ReadWriter
	Fd() uintptr
}

type baseHandler struct {
	preformatted []byte
	groupPrefix  string   // for text: prefix of groups opened in preformatting
	groups       []string // all groups started from WithGroup
	nOpenGroups  int      // the number of groups opened in preformattedAttrs
	json         bool
	mux          sync.Mutex

	// timeFormat specify what's pattern to be formatted
	// default using time.Kitchen
	//
	// eg: time.DateTime
	timeFormat string

	// w is output writer, default using os.Stderr
	w io.Writer

	// tty only tty can be colored output
	tty bool

	// level is logger min Level, default is slog.LevelInfo
	level slog.Level

	// prefix output prefix in every record
	prefix string

	// replacer refer to Replacer
	replacer Replacer

	// caller if true caller will be logged.
	caller bool

	// fullCaller: <mod/package.FunctionName:Line>
	fullCaller bool

	themes Themes
}

func (h *baseHandler) isTTY() bool {
	if f, ok := h.w.(File); ok {
		return isatty.IsTerminal(f.Fd())
	}
	return false
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
// It is called early, before any arguments are processed,
// to save effort if the log event should be discarded.
// If called from a Logger method, the first argument is the context
// passed to that method, or context.Background() if nil was passed
// or the method does not take a context.
// The context is passed so Enabled can use its values
// to make a decision.
func (h *baseHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle handles the Record.
// It will only be called when Enabled returns true.
// The Context argument is as for Enabled.
// It is present solely to provide Handlers access to the context's values.
// Canceling the context should not affect record processing.
// (Among other things, log messages may be necessary to debug a
// cancellation-related problem.)
//
// Handle methods that produce output should observe the following rules:
//   - If r.Time is the zero time, ignore the time.
//   - If r.PC is zero, ignore it.
//   - Attr's values should be resolved.
//   - If an Attr's key and value are both the zero value, ignore the Attr.
//     This can be tested with attr.Equal(Attr{}).
//   - If a group's key is empty, inline the group's Attrs.
//   - If a group has no Attrs (even if it has a non-empty key),
//     ignore it.
func (h *baseHandler) Handle(_ context.Context, r slog.Record) error {
	b := h.createBuilder(NewBuffer(), r)
	defer b.free()
	b.start()
	b.appendTime()
	b.appendLevel()
	b.appendCaller()
	b.appendPrefix()
	b.appendMessage()
	b.appendAttrs()
	b.close()
	buf := b.output()
	h.mux.Lock()
	defer h.mux.Unlock()
	_, err := h.w.Write(*buf)
	return err
}

// WithAttrs returns a new Handler whose attributes consist of
// both the receiver's attributes and the arguments.
// The Handler owns the slice: it may retain, modify or discard it.
func (h *baseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return nil
}

// WithGroup returns a new Handler with the given group appended to
// the receiver's existing groups.
// The keys of all subsequent attributes, whether added by With or in a
// Record, should be qualified by the sequence of group names.
//
// How this qualification happens is up to the Handler, so long as
// this Handler's attribute keys differ from those of another Handler
// with a different sequence of group names.
//
// A Handler should treat WithGroup as starting a Group of Attrs that ends
// at the end of the log event. That is,
//
//	logger.WithGroup("s").LogAttrs(level, msg, slog.Int("a", 1), slog.Int("b", 2))
//
// should behave like
//
//	logger.LogAttrs(level, msg, slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
//
// If the name is empty, WithGroup returns the receiver.
func (h *baseHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := h.clone()
	h2.initThemes()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *baseHandler) withPrefix(prefix string) *baseHandler {
	h2 := h.clone()
	h2.initThemes()
	h2.prefix = prefix
	return h2
}

func (h *baseHandler) withThemes(themes Themes) *baseHandler {
	h2 := h.clone()
	for k, v := range themes {
		h2.themes[k] = v
	}
	h2.initThemes()
	return h2
}

func (h *baseHandler) createBuilder(buf *Buffer, r slog.Record) Builder {
	if h.json {
		return nil
	}
	return &textBuilder{h.createBaseBuilder(buf, r)}
}

func (h *baseHandler) clone() *baseHandler {
	return &baseHandler{
		preformatted: slices.Clip(h.preformatted),
		groupPrefix:  h.groupPrefix,
		groups:       slices.Clip(h.groups),
		nOpenGroups:  h.nOpenGroups,
		json:         h.json,
		timeFormat:   h.timeFormat,
		w:            h.w,
		level:        h.level,
		prefix:       h.prefix,
		replacer:     h.replacer,
		caller:       h.caller,
		fullCaller:   h.fullCaller,
		themes:       h.themes,
	}
}
