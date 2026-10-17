package application_test

import (
	"context"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/application"
)

func TestRequestIDContext(t *testing.T) {
	ctx := application.WithRequestID(context.Background(), "req_001")
	id, ok := application.RequestID(ctx)
	if !ok {
		t.Fatal("RequestID returned ok=false")
	}
	if id != "req_001" {
		t.Errorf("RequestID = %q, want %q", id, "req_001")
	}
}

func TestRequestID_Missing(t *testing.T) {
	if id, ok := application.RequestID(context.Background()); ok {
		t.Errorf("RequestID = %q, true; want ok=false", id)
	}
}
