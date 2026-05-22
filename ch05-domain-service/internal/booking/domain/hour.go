package domain

import (
	"errors"
	"fmt"
	"time"
)

// Hour は予約時間帯の値オブジェクト (1 日内の時刻を分単位で保持)。
type Hour struct {
	// hourOfDay は 0-23
	hourOfDay int
	// minute は 0-59
	minute int
}

// ErrInvalidHour は不正な時刻を表すドメインエラー。
var ErrInvalidHour = errors.New("invalid hour")

// NewHour は HH:MM 形式の文字列から Hour を生成する。
func NewHour(s string) (Hour, error) {
	var h, m int
	if _, err := fmt.Sscanf(s, "%d:%d", &h, &m); err != nil {
		return Hour{}, fmt.Errorf("%w: %v", ErrInvalidHour, err)
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return Hour{}, fmt.Errorf("%w: out of range", ErrInvalidHour)
	}
	return Hour{hourOfDay: h, minute: m}, nil
}

// String は HH:MM 形式で時刻を返す。
func (h Hour) String() string { return fmt.Sprintf("%02d:%02d", h.hourOfDay, h.minute) }

// Before は other より前か。
func (h Hour) Before(other Hour) bool {
	return h.totalMinutes() < other.totalMinutes()
}

// After は other より後か。
func (h Hour) After(other Hour) bool {
	return h.totalMinutes() > other.totalMinutes()
}

// Equals は同値比較。
func (h Hour) Equals(other Hour) bool {
	return h.totalMinutes() == other.totalMinutes()
}

// totalMinutes は 1 日の開始からの分数を返す。
func (h Hour) totalMinutes() int { return h.hourOfDay*60 + h.minute }

// HoursOverlap は 2 つの時間帯 [a, b) と [c, d) が重なるかを判定する。
func HoursOverlap(a, b, c, d Hour) bool {
	return a.Before(d) && c.Before(b)
}

// Today は与えられた Clock の「今日」の Hour 時刻を time.Time として返す。
func (h Hour) Today(clock Clock) time.Time {
	now := clock.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), h.hourOfDay, h.minute, 0, 0, now.Location())
}
