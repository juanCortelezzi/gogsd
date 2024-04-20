package gsdlogger

import (
	"io"
	"log/slog"
)

type Logger = *slog.Logger

func NewLogger(w io.Writer, level slog.Level) Logger {
	options := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(w, options)

	return slog.New(handler)
}
