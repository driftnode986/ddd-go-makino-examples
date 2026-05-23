// identity/internal/identity/domain/member_test.go
package domain_test

import (
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/domain"
)

func TestNewEmail_Validates(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "alice@example.com", false},
		{"missing @", "alice.example.com", true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr = %v, got err = %v", tt.wantErr, err)
			}
		})
	}
}
