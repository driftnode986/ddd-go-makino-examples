package pricing

import "time"

// Coupon は割引クーポン (本章では有効期限と割引額だけを持つ)。
type Coupon struct {
	discount   Money
	validUntil time.Time
}

// NewCoupon は validUntil まで使える、discount だけ割り引くクーポンを返す。
func NewCoupon(discount Money, validUntil time.Time) *Coupon {
	return &Coupon{discount: discount, validUntil: validUntil}
}

// IsValidAt は now の時点でクーポンが使えるかを返す (validUntil ちょうどまで有効)。
func (c *Coupon) IsValidAt(
	now time.Time,
) bool {
	return !now.After(c.validUntil)
}

// DiscountAmount は割引額を返す。
func (c *Coupon) DiscountAmount() Money { return c.discount }
