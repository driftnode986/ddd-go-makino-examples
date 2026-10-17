package domain

import "time"

// #region behavior
// 推奨パターン: ビヘイビアを表現する
func (b *Booking) Confirm() error { return nil }
func (b *Booking) Cancel(
	reason string,
	clock func() time.Time,
) error {
	return nil
}
func (b *Booking) MarkCompleted() error { return nil }

// #endregion
