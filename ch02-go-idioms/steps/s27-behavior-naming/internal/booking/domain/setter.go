package domain

// #region setter
// アンチパターン: Setter は何が起きるかを呼び出し側が知らなければならない
func (b *Booking) SetStatus(s BookingStatus) { b.status = s }
func (b *Booking) SetCancelReason(r string)  { b.cancelReason = r }

// #endregion
