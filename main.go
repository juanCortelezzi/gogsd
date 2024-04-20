package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type Logger = *slog.Logger

func NewLogger(w io.Writer, level slog.Level) Logger {
	options := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(w, options)

	return slog.New(handler)
}

func main() {
	ctx := context.Background()
	logger := NewLogger(os.Stdout, slog.LevelInfo)

	if err := Run(ctx, logger, os.LookupEnv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		logger.ErrorContext(ctx, "error in top level", "err", err)
		os.Exit(1)
	}
}
