package slog_multi

import (
	"io"
	"log/slog"
)

// NewLogger は JSON ハンドラ + テキストハンドラの 2 系統に同時出力するロガーを返します。
func NewLogger(jsonOut, textOut io.Writer) *slog.Logger {
	jsonH := slog.NewJSONHandler(jsonOut, &slog.HandlerOptions{Level: slog.LevelInfo})
	textH := slog.NewTextHandler(textOut, &slog.HandlerOptions{Level: slog.LevelDebug})
	multi := slog.NewMultiHandler(jsonH, textH)
	return slog.New(multi)
}
