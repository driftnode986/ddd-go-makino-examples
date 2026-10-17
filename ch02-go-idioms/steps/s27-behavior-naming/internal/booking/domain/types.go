// Package domain は 2.7 節の「Setter とビヘイビア指向のメソッド」を並べた版。
// ビヘイビア側のメソッドは中身を省略してあり、完成形は ch02-go-idioms/internal/booking/domain にある。
package domain

import canonical "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"

type BookingStatus = canonical.BookingStatus

type Booking struct {
	status       BookingStatus
	cancelReason string
}
