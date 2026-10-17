package http_test

import (
	"encoding/json"
	"testing"

	bookinghttp "github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/adapter/http"
)

// DTO は公開フィールドなので、JSON から検証なしでそのまま値が入る
func TestRegisterUserRequest_DecodesJSON(t *testing.T) {
	var req bookinghttp.RegisterUserRequest
	if err := json.Unmarshal([]byte(`{"email":"not-an-email"}`), &req); err != nil {
		t.Fatal(err)
	}
	if req.Email != "not-an-email" {
		t.Fatalf("Email: got %q", req.Email)
	}
}
