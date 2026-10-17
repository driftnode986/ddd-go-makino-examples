// Package before は 2.15 節の演習の「リファクタ前」のコード。
// 公開フィールドと Setter で状態を書き換え、不正な入力は panic で止める。
package before

// Booking は予約。フィールドが公開されているので、誰でも直接書き換えられる。
type Booking struct {
	ID           string
	MemberID     string
	Status       string
	CancelReason string
}

// NewBooking は予約を作る。ID が空なら panic する。
func NewBooking(id, memberID string) *Booking {
	if id == "" {
		panic("booking id is required")
	}
	return &Booking{ID: id, MemberID: memberID, Status: "pending"}
}

// SetStatus は状態を書き換える。どの遷移が許されるかは呼び出し側が知っている必要がある。
func (b *Booking) SetStatus(s string) { b.Status = s }

// GetCancelReason はキャンセル理由を返す。値レシーバとポインタレシーバが混在している。
func (b Booking) GetCancelReason() string { return b.CancelReason }

// SetCancelReason はキャンセル理由を書き換える。
func (b *Booking) SetCancelReason(r string) { b.CancelReason = r }
