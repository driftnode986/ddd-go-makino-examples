// #region pricing-service
package pricing

import (
	"fmt"
	"time"
)

type Clock interface {
	Now() time.Time
}

type PricingService struct {
	clock Clock // 時刻に依存するため interface 経由で注入
}

func NewPricingService(clock Clock) *PricingService {
	return &PricingService{clock: clock}
}

// #endregion

// #region calculate-final-price
func (s *PricingService) CalculateFinalPrice(
	base Money,
	coupon *Coupon,
) (Money, error) {
	now := s.clock.Now()
	if !coupon.IsValidAt(now) {
		return Money{}, fmt.Errorf("coupon expired")
	}
	return base.Subtract(coupon.DiscountAmount())
}

// #endregion
