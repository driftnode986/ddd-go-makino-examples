package slog_multi_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/slog_multi"
)

func TestNewLogger_DispatchesToBoth(t *testing.T) {
	var jsonBuf, textBuf bytes.Buffer
	logger := slog_multi.NewLogger(&jsonBuf, &textBuf)
	logger.Info("hello", "key", "value")

	jsonOut := jsonBuf.String()
	textOut := textBuf.String()

	if !strings.Contains(jsonOut, `"msg":"hello"`) {
		t.Errorf("JSON handler did not receive record: %q", jsonOut)
	}
	if !strings.Contains(textOut, "msg=hello") {
		t.Errorf("Text handler did not receive record: %q", textOut)
	}
}
