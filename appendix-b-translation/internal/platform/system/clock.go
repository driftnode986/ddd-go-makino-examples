// Package system は本番で使う時計などの実装を置く。
package system

import "time"

// Clock は実際の現在時刻を返す。
type Clock struct{}

func NewClock() Clock { return Clock{} }

func (Clock) Now() time.Time { return time.Now() }
