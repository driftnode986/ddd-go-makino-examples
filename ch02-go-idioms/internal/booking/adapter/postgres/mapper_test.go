package postgres

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch02-go-idioms/internal/booking/domain"
)

func TestToRow_CopiesDomainFields(t *testing.T) {
	at := time.Now().Add(48 * time.Hour)
	b, err := domain.NewBooking("member-1", "lesson-1", at)
	if err != nil {
		t.Fatal(err)
	}
	if err := b.Cancel("schedule changed", time.Now); err != nil {
		t.Fatal(err)
	}
	row := toRow(b)
	want := BookingRow{
		ID:           b.ID().String(),
		MemberID:     "member-1",
		LessonID:     "lesson-1",
		ScheduledAt:  at,
		Status:       "cancelled",
		CancelReason: "schedule changed",
	}
	if row != want {
		t.Fatalf("toRow:\n got %+v\nwant %+v", row, want)
	}
}
