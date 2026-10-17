package application

import (
	"errors"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

func TestMapDomainError(t *testing.T) {
	other := errors.New("other")
	tests := []struct {
		name     string
		in       error
		wantCode string // 空なら ApplicationError にならず、そのまま返る
	}{
		{
			"invalid quantity",
			domain.ErrInvalidQuantity,
			"INVALID_QUANTITY",
		},
		{"cannot modify", domain.ErrCannotModify, "ORDER_LOCKED"},
		{"unmapped", other, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDomainError(tt.in)
			var appErr *ApplicationError
			if tt.wantCode == "" {
				if got != tt.in {
					t.Errorf("got %v, want the same error", got)
				}
				return
			}
			if !errors.As(got, &appErr) ||
				appErr.Code != tt.wantCode {
				t.Errorf("got %v, want code %s", got, tt.wantCode)
			}
			if !errors.Is(got, tt.in) {
				t.Errorf("cause is lost: %v", got)
			}
		})
	}
}
