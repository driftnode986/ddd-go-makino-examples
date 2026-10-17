// #region email-test-imports
package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/identity/domain"
)

// #endregion

// #region email-test
func TestNewEmail_RejectsInvalidFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid email", "user@example.com", nil},
		{"empty string", "", domain.ErrInvalidEmail},
		{"no at sign", "userexample.com", domain.ErrInvalidEmail},
		{"only at sign", "@", domain.ErrInvalidEmail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewEmail(tt.input)
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("want error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// #endregion
