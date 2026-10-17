package test_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

func TestFixedClock_Now(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := test.NewFixedClock(at)
	if !c.Now().Equal(at) {
		t.Errorf("Now() = %v, want %v", c.Now(), at)
	}
}
