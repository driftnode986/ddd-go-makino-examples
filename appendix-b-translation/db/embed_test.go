package db_test

import (
	"strings"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/db"
)

func TestSchemasAreEmbedded(t *testing.T) {
	if !strings.Contains(db.IdentitySchema, "CREATE TABLE users") {
		t.Errorf("IdentitySchema = %q", db.IdentitySchema)
	}
	if !strings.Contains(db.BookingSchema, "CREATE TABLE orders") {
		t.Errorf("BookingSchema = %q", db.BookingSchema)
	}
}
