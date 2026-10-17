// Package test はテストで使う固定の時計などを置く。
package test

import "time"

// FixedClock は常に同じ時刻を返す。
type FixedClock struct{ at time.Time }

func NewFixedClock(at time.Time) *FixedClock { return &FixedClock{at: at} }

func (c *FixedClock) Now() time.Time { return c.at }
