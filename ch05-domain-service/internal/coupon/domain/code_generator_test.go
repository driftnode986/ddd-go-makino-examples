package domain_test

import (
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch05-domain-service/internal/coupon/domain"
)

func TestCodeGenerator_Generate_DeterministicWithFixedRandom(t *testing.T) {
	gen := domain.NewCodeGenerator(domain.FixedRandom{Value: 0})
	code := gen.Generate()
	// FixedRandom{Value: 0}.IntN(n) は常に 0 を返す。
	// alphabet[0] = 'A' なので、結果は "AAAAAAAA"。
	if code.String() != "AAAAAAAA" {
		t.Errorf("want=AAAAAAAA got=%s", code)
	}
}

func TestCodeGenerator_Generate_FixedRandomSelectsOtherIndex(t *testing.T) {
	gen := domain.NewCodeGenerator(domain.FixedRandom{Value: 2})
	code := gen.Generate()
	// FixedRandom{Value: 2}.IntN(32) = 2 なので、alphabet[2] = 'C'
	if code.String() != "CCCCCCCC" {
		t.Errorf("want=CCCCCCCC got=%s", code)
	}
}

func TestCodeGenerator_Generate_LengthIsEight(t *testing.T) {
	gen := domain.NewCodeGenerator(domain.FixedRandom{Value: 0})
	code := gen.Generate()
	if len(code.String()) != 8 {
		t.Errorf("length: want=8 got=%d", len(code.String()))
	}
}
