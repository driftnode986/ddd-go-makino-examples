package domain

// タグ型 (userTag, bookingTag) は非公開なので、パッケージの外からは
// NewID[T]() の T を書けない。業務用語ごとの生成関数をここに置く。

// NewUserID は新しい UserID を生成する。
func NewUserID() UserID { return NewID[userTag]() }

// NewBookingID は新しい BookingID を生成する。
func NewBookingID() BookingID { return NewID[bookingTag]() }
