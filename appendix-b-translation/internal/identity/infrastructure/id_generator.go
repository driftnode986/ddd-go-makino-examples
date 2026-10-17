package infrastructure

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
)

// RandomIDGenerator は "usr_" に 16 桁の乱数を付けた ID を作る。
type RandomIDGenerator struct{}

func (RandomIDGenerator) New() domain.UserID {
	b := make([]byte, 8)
	_, _ = rand.Read(b) // crypto/rand.Read はエラーを返さない（Go 1.24 以降）
	return domain.UserID("usr_" + hex.EncodeToString(b))
}
