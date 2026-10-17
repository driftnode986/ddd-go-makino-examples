// #region status
package domain

import (
	"errors"
	"time"
)

// Status は予約の状態を表す。
type Status string

const (
	// StatusPending は予約直後で、確定前の状態。
	StatusPending Status = "pending"
	// StatusConfirmed は予約が確定した状態。
	StatusConfirmed Status = "confirmed"
	// StatusCompleted は予約レッスンを終えた状態。
	StatusCompleted Status = "completed"
	// StatusCancelled はキャンセル済みの状態。
	StatusCancelled Status = "cancelled"
)

// #endregion

// #region booking-type
// Booking は予約エンティティ。ID で識別され、ライフサイクルを持つ。
type Booking struct {
	id          BookingID
	memberID    MemberID
	lessonID    LessonID
	scheduledAt time.Time
	status      Status
	version     int64
	createdAt   time.Time
	updatedAt   time.Time
}

// #endregion

// ドメインエラー (センチネル)。
// #region errors
var (
	// ErrCannotConfirm は確認できない状態遷移を表す。
	ErrCannotConfirm = errors.New(
		"cannot confirm booking from this state",
	)
	// ErrCannotCancel はキャンセル不可な状態遷移を表す。
	ErrCannotCancel = errors.New(
		"cannot cancel booking from this state",
	)
	// ErrCannotComplete は完了できない状態遷移を表す。
	ErrCannotComplete = errors.New(
		"cannot complete booking from this state",
	)
	// ErrScheduledInPast は過去日時の予約を表す。
	ErrScheduledInPast = errors.New(
		"scheduled time must be in the future",
	)
	// ErrMissingID は予約・会員・レッスンの ID の欠けを表す。
	ErrMissingID = errors.New("id, memberID, lessonID required")
)

// #endregion

// #region new-booking
// NewBooking はファクトリ関数。
// 生成時に不変条件 (時刻が未来) を検証する。
// clock は now を返す関数で、テスト時に差し替える。
func NewBooking(
	id BookingID,
	memberID MemberID,
	lessonID LessonID,
	scheduledAt time.Time,
	clock func() time.Time,
) (*Booking, error) {
	if id.IsZero() || memberID.IsZero() || lessonID.IsZero() {
		return nil, ErrMissingID
	}
	now := clock()
	if !scheduledAt.After(now) {
		return nil, ErrScheduledInPast
	}
	return &Booking{
		id:          id,
		memberID:    memberID,
		lessonID:    lessonID,
		scheduledAt: scheduledAt,
		status:      StatusPending,
		version:     1,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// #endregion

// #region confirm-cancel
// Confirm は予約を確定する。pending → confirmed のみ許容。
func (b *Booking) Confirm(clock func() time.Time) error {
	if b.status != StatusPending {
		return ErrCannotConfirm
	}
	b.status = StatusConfirmed
	b.version++
	b.updatedAt = clock()
	return nil
}

// Cancel は予約をキャンセルする。pending または confirmed のみ。
func (b *Booking) Cancel(clock func() time.Time) error {
	if b.status != StatusPending && b.status != StatusConfirmed {
		return ErrCannotCancel
	}
	b.status = StatusCancelled
	b.version++
	b.updatedAt = clock()
	return nil
}

// #endregion

// #region complete
// Complete は予約レッスンの完了を記録する。confirmed のみ。
func (b *Booking) Complete(clock func() time.Time) error {
	if b.status != StatusConfirmed {
		return ErrCannotComplete
	}
	b.status = StatusCompleted
	b.version++
	b.updatedAt = clock()
	return nil
}

// #endregion

// #region accessors
// 状態取得 (値レシーバ、読み取り専用)。

// ID は予約 ID を返す。
func (b Booking) ID() BookingID { return b.id }

// Status は現在の状態を返す。
func (b Booking) Status() Status { return b.status }

// Version は楽観的ロック用の版番号を返す。
func (b Booking) Version() int64 { return b.version }

// #endregion

// MemberID は予約者の Member ID を返す。
func (b Booking) MemberID() MemberID { return b.memberID }

// LessonID は予約レッスンの Lesson ID を返す。
func (b Booking) LessonID() LessonID { return b.lessonID }

// ScheduledAt は予約日時を返す。
func (b Booking) ScheduledAt() time.Time { return b.scheduledAt }

// CreatedAt は作成日時を返す。
func (b Booking) CreatedAt() time.Time { return b.createdAt }

// UpdatedAt は最終更新日時を返す。
func (b Booking) UpdatedAt() time.Time { return b.updatedAt }
