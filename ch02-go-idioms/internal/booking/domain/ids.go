package domain

import shared "github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/shared/domain"

// BookingID は 2.10 節の型安全な ID (internal/shared/domain) をそのまま使う。
type BookingID = shared.BookingID

// MemberID は会員の ID。本章では文字列で受け取る。
type MemberID string

// String は ID の文字列表現を返す。
func (id MemberID) String() string { return string(id) }

// LessonID はレッスンの ID。本章では文字列で受け取る。
type LessonID string

// String は ID の文字列表現を返す。
func (id LessonID) String() string { return string(id) }

// BookingStatus は予約の状態。
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCompleted BookingStatus = "completed"
	StatusCancelled BookingStatus = "cancelled"
)
