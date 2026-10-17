// Package after は 2.15 節の演習の「リファクタ後」のコード。
// フィールドを非公開にし、コンストラクタで検証し、Setter をビヘイビア指向のメソッドに変え、
// エラーをセンチネルエラーで返し、レシーバをポインタレシーバに揃えた。
package after

import "errors"

// Status は予約の状態。
type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

var (
	ErrBookingIDRequired = errors.New("booking id is required")
	ErrCannotConfirm     = errors.New("booking cannot be confirmed in current state")
	ErrCannotCancel      = errors.New("booking cannot be cancelled in current state")
	ErrCannotComplete    = errors.New("booking cannot be completed in current state")
)

// Booking は予約エンティティ。
type Booking struct {
	id           string
	memberID     string
	status       Status
	cancelReason string
}

// NewBooking は予約を作る。ID が空ならエラーを返す (panic しない)。
func NewBooking(id, memberID string) (*Booking, error) {
	if id == "" {
		return nil, ErrBookingIDRequired
	}
	return &Booking{id: id, memberID: memberID, status: StatusPending}, nil
}

// Confirm は予約を確定する。確定前の予約だけが確定できる。
func (b *Booking) Confirm() error {
	if b.status != StatusPending {
		return ErrCannotConfirm
	}
	b.status = StatusConfirmed
	return nil
}

// Cancel は予約をキャンセルし、理由を記録する。完了・キャンセル済みの予約はキャンセルできない。
func (b *Booking) Cancel(reason string) error {
	if b.status == StatusCompleted || b.status == StatusCancelled {
		return ErrCannotCancel
	}
	b.status = StatusCancelled
	b.cancelReason = reason
	return nil
}

// MarkCompleted はレッスンの完了を記録する。確定済みの予約だけが完了できる。
func (b *Booking) MarkCompleted() error {
	if b.status != StatusConfirmed {
		return ErrCannotComplete
	}
	b.status = StatusCompleted
	return nil
}

func (b *Booking) ID() string           { return b.id }
func (b *Booking) MemberID() string     { return b.memberID }
func (b *Booking) Status() Status       { return b.status }
func (b *Booking) CancelReason() string { return b.cancelReason }
