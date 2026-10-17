package domain

// UserID は会員の ID。
type UserID string

func (id UserID) String() string { return string(id) }
