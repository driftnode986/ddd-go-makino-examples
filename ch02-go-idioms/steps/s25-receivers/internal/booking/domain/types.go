// Package domain は 2.5 節の「レシーバをポインタで揃える」を示すため、Booking を ID と状態だけに絞った版。
// 完成形は ch02-go-idioms/internal/booking/domain にある。
package domain

import canonical "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"

// ID・状態・エラーは完成形と同じものを使う。
type (
	BookingID     = canonical.BookingID
	BookingStatus = canonical.BookingStatus
)

const (
	StatusPending   = canonical.StatusPending
	StatusConfirmed = canonical.StatusConfirmed
)

var ErrBookingAlreadyConfirmed = canonical.ErrBookingAlreadyConfirmed
