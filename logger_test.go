package shandler

import (
	"golang.org/x/exp/slog"
	"math"
	"os"
	"testing"
	"time"
)

func BenchmarkLogger(b *testing.B) {
	f, _ := os.OpenFile("shandler.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer f.Close()
	slog.SetDefault(slog.New(NewTextHandler(
		//WithCaller(),
		//WithPrefix("Application"),
		//WithFullCaller(),
		WithWriter(f),
		//WithTimeFormat(time.DateTime),
		//WithLevel(slog.LevelDebug),
	)))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		slog.Warn("warn-message")
	}
}

func BenchmarkSlog(b *testing.B) {
	f, _ := os.OpenFile("slog.log", os.O_CREATE|os.O_RDWR, os.ModePerm)
	defer f.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(f, &slog.HandlerOptions{
		//AddSource: true,
		//Level:     slog.LevelDebug,
	})))

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		slog.Warn("warn-message")
	}
}

func TestOutput(t *testing.T) {
	slog.SetDefault(slog.New(NewTextHandler(
		WithCaller(),
		//WithPrefix("Application"),
		//WithFullCaller(),
		WithTimeFormat(time.DateTime),
		WithLevel(slog.LevelDebug),
	)))

	slog.Debug("debug message", slog.String("str", "string value"))
	slog.Info("info message", slog.Int("int", 888), slog.Time("ts", time.Now()))
	slog.Warn("warn-message",
		slog.Group("group",
			slog.String("one", "value1"),
			slog.Int("two", 2),
			slog.Group("inner",
				slog.String("key", "inner value"),
				slog.Float64("f", math.MaxFloat64),
			)),
		slog.Bool("b", true))

	slog.Error("error message")
	logger := slog.New(slog.Default().Handler().(Handler).WithPrefix("another"))
	logger.Info("with another prefix logged")

	logger = slog.New(slog.Default().Handler().(Handler).WithThemes(map[ThemeSchema]*Theme{
		ThemeCaller: NewTheme().Bold().Underline(),
	}))
	logger.Info("with another prefix logged")
}
