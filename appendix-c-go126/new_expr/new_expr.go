// Package new_expr demonstrates Go 1.26 new(expr) syntax for optional fields.
package new_expr

import (
	"encoding/json"
	"time"
)

type Person struct {
	Name string `json:"name"`
	// Age is the calculated value when born is set, otherwise omitted from JSON.
	Age *int `json:"age,omitempty"`
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (365.25 * 24))
}

// PersonJSON uses Go 1.26 new(expr) to initialize an optional pointer field in one line.
func PersonJSON(name string, born time.Time) ([]byte, error) {
	p := Person{
		Name: name,
		Age:  new(yearsSince(born)),
	}
	return json.Marshal(p)
}
