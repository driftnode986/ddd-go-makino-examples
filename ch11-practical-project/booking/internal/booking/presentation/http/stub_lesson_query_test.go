package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
	"github.com/google/uuid"
)

// stubLessonQuery は Catalog BC の代わりに決まったレッスンを返すテスト用の偽物。
// booking_handler_test.go の setupTestServer が使う。
type stubLessonQuery struct{ view dto.LessonView }

func (s *stubLessonQuery) FindLessonView(_ context.Context,
	_ string) (dto.LessonView, error) {
	return s.view, nil
}

func TestListBookingsByMember_HTTP(t *testing.T) {
	srv := setupTestServer(t)
	memberID := uuid.NewString()
	for range 2 {
		body := `{"member_id":"` + memberID + `","lesson_id":"` + uuid.NewString() + `"}`
		res, err := http.Post(srv.URL+"/booking/bookings", "application/json",
			strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		_ = res.Body.Close()
	}

	res, err := http.Get(srv.URL + "/booking/members/" + memberID + "/bookings")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var views []struct{ ID, Status string }
	if err := json.NewDecoder(res.Body).Decode(&views); err != nil {
		t.Fatal(err)
	}
	if len(views) != 2 {
		t.Errorf("bookings = %d, want 2", len(views))
	}
}
