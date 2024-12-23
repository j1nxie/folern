package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

type CustomHandler struct {
	out io.Writer
}

func (h *CustomHandler) Enabled(_ context.Context, l slog.Level) bool {
	return true
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006/01/02 15:04:05")
	level := r.Level.String()

	switch r.Level {
	case slog.LevelDebug:
		level = blue + level + reset
	case slog.LevelInfo:
		level = green + level + reset
	case slog.LevelWarn:
		level = yellow + level + reset
	case slog.LevelError:
		level = red + level + reset
	}

	fmt.Fprintf(h.out, "%s %s %s\n", timeStr, level, r.Message)

	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			fmt.Fprintf(h.out, "\t%s=%v\n", a.Key, a.Value)
			return true
		})
	}

	return nil
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h
}

func InitLogger() {
	handler := &CustomHandler{out: os.Stdout}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func Operation(op string, details ...any) {
	op = yellow + "[" + op + "]" + reset
	msg := fmt.Sprintf("%s %s", op, fmt.Sprint(details...))
	slog.Info(msg)
}

func Error(op string, err error, details ...any) {
	op = yellow + "[" + op + "]" + reset
	msg := fmt.Sprintf("%s %s", op, fmt.Sprint(details...))
	slog.Error(msg, "error", err)
}
