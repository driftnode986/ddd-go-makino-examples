package reflect_iter_test

import (
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/reflect_iter"
)

type goodValueObject struct {
	value string
	count int
}

type badValueObject struct {
	Value string
	count int
}

func TestAssertAllFieldsUnexported(t *testing.T) {
	// Reference fields to satisfy staticcheck
	// (the reflect helper reads them dynamically).
	_ = goodValueObject{value: "v", count: 1}
	_ = badValueObject{Value: "v", count: 1}

	if err := reflect_iter.AssertAllFieldsUnexported(
		goodValueObject{},
	); err != nil {
		t.Errorf("unexpected error on good struct: %v", err)
	}
	err := reflect_iter.AssertAllFieldsUnexported(badValueObject{})
	if err == nil {
		t.Fatalf("expected error on bad struct, got nil")
	}
	if !strings.Contains(err.Error(), "Value") {
		t.Errorf("expected error to mention 'Value', got: %v", err)
	}
}
