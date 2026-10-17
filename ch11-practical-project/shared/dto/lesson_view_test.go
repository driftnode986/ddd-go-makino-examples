package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/shared/dto"
)

// LessonView は BC 間の「公表された言語」なので、JSON の名前が変わらないことを確かめる。
func TestLessonView_JSONNames(t *testing.T) {
	v := dto.LessonView{ID: "l1", Name: "Yoga", Capacity: 10, PriceJPY: 2000,
		StartsAt: time.Date(2026, 6, 1, 9, 0, 0, 0, time.UTC), TrainerID: "t1"}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	want := `{"id":"l1","name":"Yoga","capacity":10,"price_jpy":2000,` +
		`"starts_at":"2026-06-01T09:00:00Z","trainer_id":"t1"}`
	if string(b) != want {
		t.Errorf("json = %s\nwant   %s", b, want)
	}
}
