// #region test-imports
package domain_test

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/domain"
)

// #endregion

// #region test-new-email
func TestNewEmail_RejectsInvalidFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "user@example.com", false},
		{"valid with subdomain", "user@mail.example.com", false},
		{"empty", "", true},
		{"missing at", "userexample.com", true},
		{"missing local", "@example.com", true},
		{"spaces", "user @example.com", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewEmail(%q): wantErr=%v, got err=%v",
					tt.input, tt.wantErr, err)
			}
			if err != nil &&
				!errors.Is(err, domain.ErrInvalidEmail) {
				t.Fatalf(
					"err should wrap ErrInvalidEmail, got %v",
					err,
				)
			}
		})
	}
}

// #endregion

func TestEmail_Equals_IsCaseInsensitiveOnDomain(t *testing.T) {
	a, err := domain.NewEmail("Alice@Example.com")
	if err != nil {
		t.Fatal(err)
	}
	b, err := domain.NewEmail("Alice@example.COM")
	if err != nil {
		t.Fatal(err)
	}
	if !a.Equals(b) {
		t.Fatalf(
			"domains differ only in case, want equal: a=%q b=%q",
			a,
			b,
		)
	}
}

func TestEmail_String_RoundTrip(t *testing.T) {
	e, err := domain.NewEmail("user@Example.COM")
	if err != nil {
		t.Fatal(err)
	}
	got := e.String()
	want := "user@example.com"
	if got != want {
		t.Fatalf("String(): got %q want %q", got, want)
	}
}
