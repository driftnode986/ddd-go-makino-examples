package errors_astype_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/errors_astype"
)

func TestMapStatus(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"validation", &errors_astype.ValidationError{Field: "email", Message: "invalid"}, http.StatusBadRequest},
		{"not found", &errors_astype.NotFoundError{Resource: "user", ID: "u-1"}, http.StatusNotFound},
		{"other", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors_astype.MapStatus(tt.err)
			if got != tt.want {
				t.Errorf("MapStatus(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}
