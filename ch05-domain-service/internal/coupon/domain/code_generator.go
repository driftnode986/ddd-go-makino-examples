package domain

// CodeGenerator はクーポンコード生成のドメインサービス。
type CodeGenerator struct {
	rand RandomSource
}

// NewCodeGenerator はファクトリ関数。
func NewCodeGenerator(rand RandomSource) *CodeGenerator {
	return &CodeGenerator{rand: rand}
}

// Generate は 8 文字のクーポンコードを生成する。
func (g *CodeGenerator) Generate() CouponCode {
	// O/0/I/1 を除外
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = alphabet[g.rand.IntN(len(alphabet))]
	}
	return CouponCode{value: string(buf)}
}
