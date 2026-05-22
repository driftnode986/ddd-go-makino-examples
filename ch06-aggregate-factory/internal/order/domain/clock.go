package domain

import "time"

// Clock は時刻取得の依存。
type Clock interface {
	Now() time.Time
}

// SystemClock は本番用。
type SystemClock struct{}

// Now は現在時刻を返す。
func (SystemClock) Now() time.Time { return time.Now() }

// FixedClock はテスト用。
type FixedClock struct {
	T time.Time
}

// Now は固定時刻を返す。
func (f FixedClock) Now() time.Time { return f.T }
