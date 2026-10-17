package domain

import "time"

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

// AreSlotsConflicting は 2 つの時間帯の重なりを判定する
// 状態を持たないので関数で OK
func AreSlotsConflicting(a, b TimeSlot) bool {
	return a.Start.Before(b.End) && b.Start.Before(a.End)
}
