package application

import (
	"context"
	"fmt"

	"github.com/driftnode986/ddd-go-makino-examples/ch11-practical-project/identity/internal/identity/domain"
)

type RegisterMemberInput struct {
	Email string
	Name  string
}

type RegisterMemberOutput struct {
	MemberID string
}

type RegisterMemberUseCase struct {
	repo domain.MemberRepository
}

func NewRegisterMemberUseCase(repo domain.MemberRepository) *RegisterMemberUseCase {
	return &RegisterMemberUseCase{repo: repo}
}

func (u *RegisterMemberUseCase) Execute(ctx context.Context,
	in RegisterMemberInput) (RegisterMemberOutput, error) {
	email, err := domain.NewEmail(in.Email)
	if err != nil {
		return RegisterMemberOutput{}, fmt.Errorf("invalid input: %w", err)
	}
	m, err := domain.NewMember(email, in.Name)
	if err != nil {
		return RegisterMemberOutput{}, fmt.Errorf("invalid input: %w", err)
	}
	if err := u.repo.Save(ctx, m); err != nil {
		return RegisterMemberOutput{}, fmt.Errorf("save: %w", err)
	}
	return RegisterMemberOutput{MemberID: m.ID().String()}, nil
}
