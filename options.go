package shandler

import (
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
)

// Replacer is called to rewrite each non-group attribute before it is logged.
// The attribute's value has been resolved (see [Value.Resolve]).
// If Replacer returns an Attr with Key == "", the attribute is discarded.
//
// The built-in attributes with keys "time", "level", "source", and "msg"
// are passed to this function, except that time is omitted
// if zero, and source is omitted if AddSource is false.
//
// The first argument is a list of currently open groups that contain the
// Attr. It must not be retained or modified. Replacer is never called
// for Group attributes, only their contents. For example, the attribute
// list
//
//	Int("a", 1), Group("g", Int("b", 2)), Int("c", 3)
//
// results in consecutive calls to Replacer with the following arguments:
//
//	nil, Int("a", 1)
//	[]string{"g"}, Int("b", 2)
//	nil, Int("c", 3)
//
// Replacer can be used to change the default keys of the built-in
// attributes, convert types (for example, to replace a `time.Time` with the
// integer seconds since the Unix epoch), sanitize personal information, or
// remove attributes from the output.
type Replacer func(groups []string, a slog.Attr) slog.Attr

type Option func(*baseHandler)

func createHandler(opts ...Option) *baseHandler {
	options := &baseHandler{
		timeFormat: time.Kitchen,
		w:          os.Stderr,
		level:      slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithTimeFormat(format string) Option {
	return func(cfg *baseHandler) {
		cfg.timeFormat = format
	}
}

func WithWriter(w io.Writer) Option {
	return func(cfg *baseHandler) {
		cfg.w = w
	}
}

func WithLevel(level slog.Level) Option {
	return func(cfg *baseHandler) {
		cfg.level = level
	}
}

func WithPrefix(prefix string) Option {
	return func(cfg *baseHandler) {
		cfg.prefix = prefix
	}
}

// WithReplacer please refer to Replacer
func WithReplacer(fn Replacer) Option {
	return func(cfg *baseHandler) {
		cfg.replacer = fn
	}
}

func WithCaller() Option {
	return func(cfg *baseHandler) {
		cfg.caller = true
	}
}

func WithFullCaller() Option {
	return func(cfg *baseHandler) {
		cfg.fullCaller = true
	}
}
