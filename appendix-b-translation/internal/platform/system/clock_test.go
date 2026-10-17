package system_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/system"
)

func TestClock_Now(t *testing.T) {
	before := time.Now()
	got := system.NewClock().Now()
	if got.Before(before) || got.After(time.Now()) {
		t.Errorf("Now() = %v, want between %v and now", got, before)
	}
}
