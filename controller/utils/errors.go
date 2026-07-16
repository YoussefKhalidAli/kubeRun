package utils

import (
	"context"
	"log/slog"
)

func HandelError(err error, t string, details string) {
	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	var severity string
	if len(t) > 0 {
		severity = string(t[len(t)-1])
	}

	var level slog.Level
	switch severity {
	case "H":
		level = slog.LevelError
	case "M":
		level = slog.LevelWarn
	case "L":
		level = slog.LevelDebug
	default:
		level = slog.LevelError
	}

	Logger.Log(context.Background(), level, "handled error",
		slog.String("error_code", t),
		slog.String("severity", severity),
		slog.String("details", details),
		slog.String("error", errStr),
	)

	if severity == "H" {
		panic(err)
	}
}
