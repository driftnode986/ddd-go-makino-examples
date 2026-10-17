package persistence_test

import (
	"reflect"
	"testing"

	"github.com/driftnode986/ddd-go-makino-examples/ch03-value-object/internal/booking/adapter/persistence"
)

// Persistence Model の db タグがカラム名に対応している
func TestUserRow_DBTags(t *testing.T) {
	typ := reflect.TypeFor[persistence.UserRow]()
	for field, want := range map[string]string{"ID": "id", "Email": "email"} {
		f, ok := typ.FieldByName(field)
		if !ok {
			t.Fatalf("field %s not found", field)
		}
		if got := f.Tag.Get("db"); got != want {
			t.Errorf("%s: db tag got %q want %q", field, got, want)
		}
	}
}
