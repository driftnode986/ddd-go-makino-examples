package new_expr

import (
	"encoding/json"
	"time"
)

type Person struct {
	Name string `json:"name"`
	// Age は生年月日があれば値、なければ nil で JSON に出さない。
	Age *int `json:"age,omitempty"`
}

func yearsSince(t time.Time) int {
	return int(time.Since(t).Hours() / (365.25 * 24))
}

// PersonJSON は Go 1.26 の new(expr) で
// Optional フィールドを 1 行で初期化します。
func PersonJSON(name string, born time.Time) ([]byte, error) {
	p := Person{Name: name}
	if !born.IsZero() { // 生年月日が未登録なら Age は nil のまま
		p.Age = new(yearsSince(born))
	}
	return json.Marshal(p)
}
