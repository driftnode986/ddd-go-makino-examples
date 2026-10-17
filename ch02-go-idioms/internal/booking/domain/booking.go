// #region booking-type
package domain

import "time"

type Booking struct {
	id           BookingID
	memberID     MemberID
	lessonID     LessonID
	scheduledAt  time.Time
	status       BookingStatus
	cancelReason string
}

// #endregion

// 状態変更メソッドはポインタレシーバ
func (b *Booking) Confirm() error {
	if b.status != StatusPending {
		return ErrBookingAlreadyConfirmed
	}
	b.status = StatusConfirmed
	return nil
}

// #region cancel
func (b *Booking) Cancel(reason string) error {
	if b.status == StatusCompleted {
		return ErrBookingCannotCancel
	}
	if time.Until(b.scheduledAt) < 24*time.Hour {
		return ErrBookingCannotCancel
	}
	b.status = StatusCancelled
	b.cancelReason = reason
	return nil
}

// #endregion

// Reschedule は予約日時を変更する。キャンセル済み・完了済みの予約は変更できない。
func (b *Booking) Reschedule(newAt time.Time) error {
	if b.status == StatusCancelled || b.status == StatusCompleted {
		return ErrBookingCannotReschedule
	}
	b.scheduledAt = newAt
	return nil
}

// MarkCompleted はレッスンの完了を記録する。確定済みの予約だけが完了できる。
func (b *Booking) MarkCompleted() error {
	if b.status != StatusConfirmed {
		return ErrBookingCannotComplete
	}
	b.status = StatusCompleted
	return nil
}

// 読み取り専用メソッドもポインタレシーバで一貫
func (b *Booking) ID() BookingID         { return b.id }
func (b *Booking) Status() BookingStatus { return b.status }

func (b *Booking) MemberID() MemberID     { return b.memberID }
func (b *Booking) LessonID() LessonID     { return b.lessonID }
func (b *Booking) ScheduledAt() time.Time { return b.scheduledAt }
func (b *Booking) CancelReason() string   { return b.cancelReason }
