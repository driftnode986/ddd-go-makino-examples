// #region register-user
package application

import (
	"context"

	"github.com/driftnode986/ddd-go-makino-examples/appendix-b-translation/internal/identity/domain"
)

type RegisterUserUseCase struct {
	txMgr TxManager
	users domain.UserRepository
	clock Clock
	idGen IDGenerator
}

func (uc *RegisterUserUseCase) Execute(
	ctx context.Context,
	cmd RegisterUserCommand,
) (domain.UserID, error) {
	var resultID domain.UserID
	err := uc.txMgr.RunInTx(ctx, func(ctx context.Context) error {
		user, err := domain.NewUser(
			uc.idGen.New(),
			cmd.Email,
			cmd.Name,
			uc.clock.Now(),
		)
		if err != nil {
			return err
		}
		if err := uc.users.Add(ctx, user); err != nil {
			return err
		}
		resultID = user.ID()
		return nil
	})
	return resultID, err
}

// #endregion

func NewRegisterUserUseCase(
	txMgr TxManager,
	users domain.UserRepository,
	clock Clock,
	idGen IDGenerator,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		txMgr: txMgr,
		users: users,
		clock: clock,
		idGen: idGen,
	}
}
