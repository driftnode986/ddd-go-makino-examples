package infrastructure

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/application"
	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

// StubPricing は本書では商品マスタを持たない簡易実装。
// 本番ではドメインサービス (Ch5) のインターフェース実装に差し替える。
type StubPricing struct {
	defaultUnit int64
}

func NewStubPricing() *StubPricing {
	return &StubPricing{defaultUnit: 1000}
}

func (p *StubPricing) UnitPrice(
	ctx context.Context, productID domain.ProductID,
) (domain.Money, error) {
	return domain.MustMoney(p.defaultUnit, "JPY"), nil
}

var _ application.PricingPort = (*StubPricing)(nil)
