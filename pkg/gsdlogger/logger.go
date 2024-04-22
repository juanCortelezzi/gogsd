package gsdlogger

import (
	"io"
	"log/slog"
	"net/http"
)

type Logger = *slog.Logger

type LoggerResponseWritter struct {
	http.ResponseWriter

	status      int
	wroteHeader bool
}

func NewLoggerResponseWritter(w http.ResponseWriter) *LoggerResponseWritter {
	return &LoggerResponseWritter{ResponseWriter: w}
}

func (rw *LoggerResponseWritter) Status() int {
	return rw.status
}

func (rw *LoggerResponseWritter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	return
}

func NewLogger(w io.Writer, level slog.Level) Logger {
	options := &slog.HandlerOptions{Level: level}
	handler := slog.NewTextHandler(w, options)

	return slog.New(handler)
}
