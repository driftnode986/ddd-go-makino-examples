// Package domain は 2.5 節のアンチパターン (ポインタレシーバと値レシーバの混在) を示す版。
package domain

import canonical "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"

type BookingID = canonical.BookingID

type Booking struct {
	id BookingID
}
