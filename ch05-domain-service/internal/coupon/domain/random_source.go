package domain

import "math/rand/v2"

// RandomSource は乱数取得の依存。
type RandomSource interface {
	IntN(n int) int
}

// SystemRandom は本番用の乱数源。
type SystemRandom struct {
	R *rand.Rand
}

// IntN は [0, n) の整数を返す。
func (s SystemRandom) IntN(n int) int { return s.R.IntN(n) }

// FixedRandom はテスト用の乱数源。固定値を返す。
type FixedRandom struct {
	Value int
}

// IntN は固定値を n で割った余りを返す (テスト時の境界調整用)。
func (f FixedRandom) IntN(n int) int { return f.Value % n }
