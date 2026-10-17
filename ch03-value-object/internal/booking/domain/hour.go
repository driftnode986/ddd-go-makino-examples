// #region new-hour
package domain

import (
	"errors"
	"time"
)

// Hour は予約の開始時刻を表す値オブジェクト。
// 「営業時間内 (9 時から 21 時)」「分単位なし」の
// 業務ルールを内包する。
type Hour struct {
	t time.Time
}

// ErrOutOfBusinessHours は営業時間外の時刻を表すドメインエラー。
var ErrOutOfBusinessHours = errors.New(
	"hour must be between 9 and 21",
)

// ErrHourHasMinutes は分単位を含む時刻を表すドメインエラー。
var ErrHourHasMinutes = errors.New("hour must be on the hour")

// NewHour は業務ルールを満たす時刻を受け取り、Hour を返す。
func NewHour(t time.Time) (Hour, error) {
	if t.Minute() != 0 || t.Second() != 0 || t.Nanosecond() != 0 {
		return Hour{}, ErrHourHasMinutes
	}
	h := t.Hour()
	if h < 9 || h > 21 {
		return Hour{}, ErrOutOfBusinessHours
	}
	return Hour{t: t}, nil
}

// #endregion

// #region hour-methods
// Time は内部の time.Time を返す (Persistence Model 変換用)。
func (h Hour) Time() time.Time { return h.t }

// IsInThePast は現在時刻と比較して過去かどうかを判定する。
// clock は now を返す関数で、テストで差し替えられる。
func (h Hour) IsInThePast(clock func() time.Time) bool {
	return h.t.Before(clock())
}

// #endregion
