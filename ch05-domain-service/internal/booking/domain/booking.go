package domain

import (
	"errors"
	"time"
)

// Status は予約の状態。
type Status string

const (
	StatusPending   Status = "pending"
	StatusConfirmed Status = "confirmed"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

// Booking は予約エンティティ (本章では衝突検出・有効期限判定の対象として扱う最小実装)。
type Booking struct {
	id          BookingID
	memberID    MemberID
	trainerID   TrainerID
	from        Hour
	to          Hour
	scheduledAt time.Time
	status      Status
	createdAt   time.Time
}

// ErrInvalidPeriod は時間帯が不正なときのドメインエラー。
var ErrInvalidPeriod = errors.New("from must be before to")

// NewBooking はファクトリ関数。
func NewBooking(
	id BookingID,
	memberID MemberID,
	trainerID TrainerID,
	from, to Hour,
	scheduledAt time.Time,
	clock Clock,
) (*Booking, error) {
	if id.IsZero() || memberID.IsZero() || trainerID.IsZero() {
		return nil, errors.New("id, memberID, trainerID required")
	}
	if !from.Before(to) {
		return nil, ErrInvalidPeriod
	}
	return &Booking{
		id:          id,
		memberID:    memberID,
		trainerID:   trainerID,
		from:        from,
		to:          to,
		scheduledAt: scheduledAt,
		status:      StatusPending,
		createdAt:   clock.Now(),
	}, nil
}

// ID は予約 ID。
func (b Booking) ID() BookingID { return b.id }

// MemberID は予約者 ID。
func (b Booking) MemberID() MemberID { return b.memberID }

// TrainerID は担当トレーナー ID。
func (b Booking) TrainerID() TrainerID { return b.trainerID }

// From は開始時刻。
func (b Booking) From() Hour { return b.from }

// To は終了時刻。
func (b Booking) To() Hour { return b.to }

// ScheduledAt は予約日時。
func (b Booking) ScheduledAt() time.Time { return b.scheduledAt }

// Status は状態。
func (b Booking) Status() Status { return b.status }

// CreatedAt は作成日時。
func (b Booking) CreatedAt() time.Time { return b.createdAt }
