package shandler

import (
	"golang.org/x/exp/slog"
	"testing"
	"time"
)

func BenchmarkLogger(b *testing.B) {
	slog.SetDefault(slog.New(NewTextHandler(
		WithCaller(),
		WithPrefix("Application"),
		//WithFullCaller(),
		WithTimeFormat(time.DateTime),
		WithLevel(slog.LevelDebug),
	)))

	b.ResetTimer()
	slog.Warn("warn-message")
}

func TestOutput(t *testing.T) {
	slog.SetDefault(slog.New(NewTextHandler(
		WithCaller(),
		WithPrefix("Application"),
		//WithFullCaller(),
		WithTimeFormat(time.DateTime),
		WithLevel(slog.LevelDebug),
	)))

	slog.Debug("debug message", slog.String("str", "string value"))
	slog.Info("info message", slog.Int("int", 888))
	slog.Warn("warn-message",
		slog.Group("group",
			slog.String("one", "value1"),
			slog.Int("two", 2),
			slog.Group("inner",
				slog.String("inner key", "inner value"))),
		slog.Bool("b", true))

	slog.Error("error message")
	logger := slog.New(slog.Default().Handler().(Handler).WithPrefix("another"))
	logger.Info("with another prefix logged")

	logger = slog.New(slog.Default().Handler().(Handler).WithThemes(map[ThemeSection]*Theme{
		ThemeCaller: NewTheme().Bold().Underline(),
	}))
	logger.Info("with another prefix logged")
}
