package application

import (
	"errors"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch08-application-service/internal/order/domain"
)

type ApplicationError struct {
	Code    string
	Message string
	Cause   error
}

func (e *ApplicationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *ApplicationError) Unwrap() error {
	return e.Cause
}

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
