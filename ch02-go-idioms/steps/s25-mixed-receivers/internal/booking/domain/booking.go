package domain

// #region mixed-receivers
// アンチパターン: ポインタと値の混在
func (b *Booking) Confirm() error { return nil }
func (b Booking) ID() BookingID   { return b.id } // 値レシーバ
// #endregion
