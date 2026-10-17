package domain

import "time"

// ExpirationChecker は予約の有効期限切れを判定するドメインサービス。
type ExpirationChecker struct {
	clock Clock
}

// NewExpirationChecker はファクトリ関数。
func NewExpirationChecker(clock Clock) *ExpirationChecker {
	return &ExpirationChecker{clock: clock}
}

// IsExpired は予約が有効期限切れかを判定する。
// 「予約日時の 30 分後を過ぎても確定されなかった予約は失効」というドメインルール。
func (c *ExpirationChecker) IsExpired(b *Booking) bool {
	if b.Status() != StatusPending {
		return false
	}
	deadline := b.ScheduledAt().Add(30 * time.Minute)
	return c.clock.Now().After(deadline)
}
