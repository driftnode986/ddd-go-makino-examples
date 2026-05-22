package application

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

// FixedClock はテストで使う固定時刻 Clock。
type FixedClock struct {
	At time.Time
}

func (c FixedClock) Now() time.Time { return c.At }
