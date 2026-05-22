package application

import (
	"errors"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch09-layered-architecture/internal/order/domain"
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
			Code:    "INVALID_INPUT",
			Message: "数量が正しくありません",
			Cause:   err,
		}
	case errors.Is(err, domain.ErrCannotModify):
		return &ApplicationError{
			Code:    "CONFLICT",
			Message: "注文は確定済みのため変更できません",
			Cause:   err,
		}
	case errors.Is(err, domain.ErrNotFound):
		return &ApplicationError{
			Code:    "NOT_FOUND",
			Message: "注文が見つかりません",
			Cause:   err,
		}
	case errors.Is(err, domain.ErrInvalidID):
		return &ApplicationError{
			Code:    "INVALID_INPUT",
			Message: "ID が正しくありません",
			Cause:   err,
		}
	default:
		return err
	}
}
