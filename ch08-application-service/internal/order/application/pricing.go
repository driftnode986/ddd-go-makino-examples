package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

// PricingPort はアプリケーション層から呼ぶ pricing ドメインサービスの port。
// 実装は infrastructure/stub_pricing.go または別 BC の adapter で行う。
type PricingPort interface {
	UnitPrice(ctx context.Context, productID domain.ProductID) (domain.Money, error)
}
