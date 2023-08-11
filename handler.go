package shandler

import (
	"context"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
	"io"
	"sync"
)

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

	// level is logger min Level, default is slog.LevelInfo
	level slog.Level

	// prefix output prefix in every record
	prefix string

	// replacer report to Replacer
	replacer Replacer

	// caller if true caller will be logged.
	caller bool

	// fullCaller: <mod/package.FunctionName:Line>
	fullCaller bool
}

func NewTextHandler(opts ...Option) slog.Handler {
	return createHandler(opts...)
}

func NewJsonHandler(opts ...Option) slog.Handler {
	return createHandler(opts...)
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
	return nil
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
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *baseHandler) getLevelString(l slog.Level) string {
	switch l {
	case slog.LevelInfo:
		return infoStyle.Render("INFO")
	case slog.LevelDebug:
		return debugStyle.Render("DBUG")
	case slog.LevelWarn:
		return warnStyle.Render("WARN")
	case slog.LevelError:
		return errorStyle.Render("ERRO")
	default:
		return boldStyle.Render(l.String()[:4])
	}
}

func (h *baseHandler) clone() *baseHandler {
	return &baseHandler{
		preformatted: slices.Clip(h.preformatted),
		groupPrefix:  h.groupPrefix,
		groups:       slices.Clip(h.groups),
		nOpenGroups:  h.nOpenGroups,
		timeFormat:   h.timeFormat,
		w:            h.w,
		level:        h.level,
		prefix:       h.prefix,
		replacer:     h.replacer,
		caller:       h.caller,
		fullCaller:   h.fullCaller,
	}
}
