package logger

import (
	"fmt"
	"log/slog"
)

func Operation(op string, details ...any) {
	msg := fmt.Sprintf("[%s] %s", op, fmt.Sprint(details...))
	slog.Info(msg)
}

func Error(op string, err error, details ...any) {
	msg := fmt.Sprintf("[%s] %s", op, fmt.Sprint(details...))
	slog.Error(msg, "error", err)
}
