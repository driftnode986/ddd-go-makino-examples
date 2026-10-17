package publishing_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/application/publishing"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/domain"
	infra "github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/booking/infrastructure"
	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/platform/test"
)

type fakeTx struct{ calls int }

func (f *fakeTx) RunInTx(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	f.calls++
	return fn(ctx)
}

type recordingPublisher struct {
	events []domain.DomainEvent
	err    error
}

func (p *recordingPublisher) Publish(
	_ context.Context,
	e domain.DomainEvent,
) error {
	if p.err != nil {
		return p.err
	}
	p.events = append(p.events, e)
	return nil
}

func TestPlaceOrderUseCase_PublishesOrderPlaced(t *testing.T) {
	ctx := context.Background()
	tx := &fakeTx{}
	orders := infra.NewInMemoryOrderRepository()
	pub := &recordingPublisher{}
	uc := publishing.NewPlaceOrderUseCase(
		tx,
		orders,
		pub,
		test.NewFixedClock(time.Now()),
	)

	err := uc.Execute(ctx, application.PlaceOrderCommand{
		CustomerID: "cus_001",
		Items: []domain.OrderItem{
			{ProductID: "p1", Quantity: 1, UnitPrice: 500},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if tx.calls != 1 || len(pub.events) != 1 {
		t.Fatalf(
			"tx calls = %d, events = %d",
			tx.calls,
			len(pub.events),
		)
	}
	placed, ok := pub.events[0].(domain.OrderPlaced)
	if !ok {
		t.Fatalf("event = %T, want OrderPlaced", pub.events[0])
	}
	if _, err := orders.FindByID(ctx, placed.OrderID); err != nil {
		t.Errorf("published order is not saved: %v", err)
	}
}

func TestPlaceOrderUseCase_ReturnsPublishError(t *testing.T) {
	errBroker := errors.New("broker down")
	uc := publishing.NewPlaceOrderUseCase(
		&fakeTx{},
		infra.NewInMemoryOrderRepository(),
		&recordingPublisher{
			err: errBroker,
		},
		test.NewFixedClock(time.Now()),
	)
	err := uc.Execute(
		context.Background(),
		application.PlaceOrderCommand{
			CustomerID: "cus_001",
			Items: []domain.OrderItem{
				{ProductID: "p1", Quantity: 1, UnitPrice: 500},
			},
		},
	)
	if !errors.Is(err, errBroker) {
		t.Errorf(
			"err = %v, want broker down"+
				"（トランザクションをロールバックさせるためにエラーを返す）",
			err,
		)
	}
}
