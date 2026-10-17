package application

import (
	"errors"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidQuantity):
		return &ApplicationError{
			Code:    "INVALID_QUANTITY",
			Message: "数量が正しくありません",
			Cause:   err,
		}
	case errors.Is(err, domain.ErrCannotModify):
		return &ApplicationError{
			Code:    "ORDER_LOCKED",
			Message: "注文は確定済みのため変更できません",
			Cause:   err,
		}
	default:
		return err
	}
}
