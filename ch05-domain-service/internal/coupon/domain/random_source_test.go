package domain_test

import (
	"math/rand/v2"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/coupon/domain"
)

func TestFixedRandom_IntN_ReturnsValueModN(t *testing.T) {
	var r domain.RandomSource = domain.FixedRandom{Value: 34}
	if got := r.IntN(32); got != 2 {
		t.Fatalf("IntN(32): got %d want 2", got)
	}
}

func TestSystemRandom_IntN_StaysInRange(t *testing.T) {
	var r domain.RandomSource = domain.SystemRandom{
		R: rand.New(rand.NewPCG(1, 2)),
	}
	for range 100 {
		if got := r.IntN(32); got < 0 || got >= 32 {
			t.Fatalf("IntN(32) out of range: %d", got)
		}
	}
}
