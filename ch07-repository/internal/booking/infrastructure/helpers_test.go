package infrastructure_test

import (
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/internal/booking/domain"
)

// このファイルは in-memory と PostgreSQL の両方のテストが使う補助関数。

func mustBookingID(t *testing.T, s string) domain.BookingID {
	t.Helper()
	id, err := domain.NewBookingID(s)
	if err != nil {
		t.Fatalf("NewBookingID(%q): %v", s, err)
	}
	return id
}

func mustMemberID(t *testing.T, s string) domain.MemberID {
	t.Helper()
	id, err := domain.NewMemberID(s)
	if err != nil {
		t.Fatalf("NewMemberID(%q): %v", s, err)
	}
	return id
}

func mustLessonID(t *testing.T, s string) domain.LessonID {
	t.Helper()
	id, err := domain.NewLessonID(s)
	if err != nil {
		t.Fatalf("NewLessonID(%q): %v", s, err)
	}
	return id
}

// mustNewBooking は ID が bk_001 の予約を作る (version は 1)。
func mustNewBooking(t *testing.T, isoTime string) *domain.Booking {
	t.Helper()
	return mustNewBookingWithID(t, "bk_001", "mem_001", isoTime)
}

func mustNewBookingWithID(
	t *testing.T, id, memberID, isoTime string,
) *domain.Booking {
	t.Helper()
	at, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		t.Fatalf("Parse(%q): %v", isoTime, err)
	}
	return domain.NewBooking(
		mustBookingID(t, id),
		mustMemberID(t, memberID),
		mustLessonID(t, "lsn_001"),
		at,
	)
}

// mustNewBookingWithVersion は指定した version の予約を復元する
// (別のプロセスが持っている古い集約を模す)。
func mustNewBookingWithVersion(
	t *testing.T, id domain.BookingID, version int,
) *domain.Booking {
	t.Helper()
	at := time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)
	return domain.ReconstructBooking(
		id,
		mustMemberID(t, "mem_001"),
		mustLessonID(t, "lsn_001"),
		domain.StatusConfirmed,
		at,
		version,
	)
}
