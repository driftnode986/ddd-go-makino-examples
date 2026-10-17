package postgres_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/steps/s22-postgres-model/internal/booking/adapter/postgres"
)

// Persistence Model は公開フィールドなので、パッケージの外から直接組み立てられる。
func TestBookingRow_IsPlainPublicStruct(t *testing.T) {
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	row := postgres.BookingRow{
		ID:          "b-1",
		MemberID:    "m-1",
		LessonID:    "l-1",
		ScheduledAt: at,
		Status:      "pending",
	}
	row.Status = "confirmed"
	if row.Status != "confirmed" || !row.ScheduledAt.Equal(at) {
		t.Fatalf("got %+v", row)
	}
}
