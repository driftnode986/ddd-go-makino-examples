package domain_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/catalog/internal/catalog/domain"
)

func TestNewLesson_Validates(t *testing.T) {
	tests := []struct {
		name     string
		lesson   string
		capacity int
		price    int
		wantErr  bool
	}{
		{"valid", "Yoga", 10, 3000, false},
		{"empty name", "", 10, 3000, true},
		{"zero capacity", "Yoga", 0, 3000, true},
		{"negative price", "Yoga", 10, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewLesson(tt.lesson, tt.capacity, tt.price,
				time.Now(), "trainer-1")
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr = %v, got err = %v", tt.wantErr, err)
			}
		})
	}
}
