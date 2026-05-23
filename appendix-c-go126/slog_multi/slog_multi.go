// Package slog_multi demonstrates Go 1.26 slog.NewMultiHandler.
package slog_multi

import (
	"io"
	"log/slog"
)

// NewLogger returns a Logger that dispatches to both JSON and text handlers.
func NewLogger(jsonOut, textOut io.Writer) *slog.Logger {
	jsonH := slog.NewJSONHandler(jsonOut, &slog.HandlerOptions{Level: slog.LevelInfo})
	textH := slog.NewTextHandler(textOut, &slog.HandlerOptions{Level: slog.LevelDebug})
	multi := slog.NewMultiHandler(jsonH, textH)
	return slog.New(multi)
}
