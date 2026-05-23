// identity/internal/identity/domain/repository.go
package domain

import (
	"context"
	"errors"
)

var ErrMemberNotFound = errors.New("member not found")

type MemberRepository interface {
	Save(ctx context.Context, m *Member) error
	FindByID(ctx context.Context, id MemberID) (*Member, error)
}
