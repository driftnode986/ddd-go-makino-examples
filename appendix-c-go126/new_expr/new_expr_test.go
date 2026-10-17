package new_expr_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-c-go126/new_expr"
)

func TestPersonJSON(t *testing.T) {
	born := time.Now().AddDate(-30, 0, 0)
	data, err := new_expr.PersonJSON("Alice", born)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", got["name"])
	}
	// age is calculated; 30 +/- a few based on leap years.
	age, ok := got["age"].(float64)
	if !ok {
		t.Fatalf("age missing or wrong type: %T", got["age"])
	}
	if age < 28 || age > 32 {
		t.Errorf("age = %v, want around 30", age)
	}
	if !strings.Contains(string(data), `"age"`) {
		t.Errorf("expected age field in JSON: %s", data)
	}
}

func TestPersonJSON_NoBirthday(t *testing.T) {
	data, err := new_expr.PersonJSON("Bob", time.Time{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(string(data), `"age"`) {
		t.Errorf("age must be omitted without a birthday: %s", data)
	}
}
