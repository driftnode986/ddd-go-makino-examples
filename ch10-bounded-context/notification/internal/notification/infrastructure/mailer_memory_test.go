package infrastructure_test

import (
	"context"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch10-bounded-context/notification/internal/notification/infrastructure"
)

func TestMemoryMailer_RecordsMail(t *testing.T) {
	m := infrastructure.NewMemoryMailer()
	if err := m.Send(context.Background(), "cus_001", "hello"); err != nil {
		t.Fatal(err)
	}
	if m.SentCount() != 1 {
		t.Fatalf("SentCount = %d, want 1", m.SentCount())
	}
	got := m.Records()[0]
	if got.CustomerID != "cus_001" || got.Body != "hello" {
		t.Errorf("record = %+v", got)
	}
}
