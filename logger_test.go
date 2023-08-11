package shandler

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	f, err := os.OpenFile("text.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	require.NoError(t, err)
	defer f.Close()
	logger := slog.New(NewTextHandler(
		WithCaller(),
		//WithWriter(f),
		WithPrefix("Application"),
		WithLevel(slog.LevelDebug),
	))
	logger.Debug("debug message", slog.String("str", "string value"))
	logger.Info("info message", slog.Int("int", 888))
	logger.Warn("warn message")
	logger.Error("warn message")
}
