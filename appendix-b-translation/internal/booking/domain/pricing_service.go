package domain

// PricingService は注文の税込み合計を計算するドメインサービス。
type PricingService struct {
	taxRatePercent int64
}

// NewPricingService は税率（%）を受け取る。
func NewPricingService(taxRatePercent int64) PricingService {
	return PricingService{taxRatePercent: taxRatePercent}
}

// Total は税込みの合計金額（円）。1 円未満は切り捨てる。
func (s PricingService) Total(o *Order) int64 {
	sub := o.Subtotal()
	return sub + sub*s.taxRatePercent/100
}
