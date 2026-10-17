package domain

type Booking struct {
	id     BookingID
	status BookingStatus
}

// 状態変更メソッドはポインタレシーバ
func (b *Booking) Confirm() error {
	if b.status != StatusPending {
		return ErrBookingAlreadyConfirmed
	}
	b.status = StatusConfirmed
	return nil
}

// 読み取り専用メソッドもポインタレシーバで一貫
func (b *Booking) ID() BookingID         { return b.id }
func (b *Booking) Status() BookingStatus { return b.status }
