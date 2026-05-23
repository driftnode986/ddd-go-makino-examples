// Package reflect_iter demonstrates Go 1.26 reflect.Type.Fields() iterator.
package reflect_iter

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

// AssertAllFieldsUnexported reports an error if any struct field is exported.
func AssertAllFieldsUnexported(v any) error {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %s", t.Kind())
	}
	var exported []string
	for f := range t.Fields() {
		if isExported(f.Name) {
			exported = append(exported, f.Name)
		}
	}
	if len(exported) > 0 {
		return fmt.Errorf("exported fields found: %s", strings.Join(exported, ", "))
	}
	return nil
}

func isExported(name string) bool {
	if name == "" {
		return false
	}
	return unicode.IsUpper(rune(name[0]))
}
