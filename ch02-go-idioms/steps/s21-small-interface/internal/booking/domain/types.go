// Package domain は 2.1 節の「小さい interface とメソッドの多い struct」を示す版。
// メソッドの中身は省略 (/* ... */) してあり、
// 完成形は ch02-go-idioms/internal/booking/domain にある。
package domain

import canonical "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"

// ID と状態の型は完成形と同じものを使う。
type (
	BookingID     = canonical.BookingID
	MemberID      = canonical.MemberID
	LessonID      = canonical.LessonID
	BookingStatus = canonical.BookingStatus
)
