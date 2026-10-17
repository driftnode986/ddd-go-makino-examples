package db_test

import (
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch07-repository/db"
)

func TestSchema_EmbedsBookingsTable(t *testing.T) {
	if !strings.Contains(db.Schema, "CREATE TABLE bookings") {
		t.Errorf(
			"Schema does not contain the bookings table:\n%s",
			db.Schema,
		)
	}
}
