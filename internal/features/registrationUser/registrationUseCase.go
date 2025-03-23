package registrationUser

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/swaggest/usecase/status"
)

type RegistrationUserUseCaseImpl struct {
	cmdHandler RegistrationUserCmdHandler
}

func NewRegistrationUserUseCase(cmdHandler RegistrationUserCmdHandler) *RegistrationUserUseCaseImpl {
	return &RegistrationUserUseCaseImpl{cmdHandler: cmdHandler}
}

type RegistrationUserCmdHandler interface {
	Execute(ctx context.Context, login string, password string) error
}

func (uc *RegistrationUserUseCaseImpl) Execute(ctx context.Context, input *RegisterRequest, output *api.EmptyBody) error {
	// Вызов юзкейса регистрации пользователя
	err := uc.cmdHandler.Execute(ctx, input.Login, input.Password)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLoginAlreadyUsed):
			return status.Wrap(errors.New("login is already used"), status.AlreadyExists)
		case errors.Is(err, domain.ErrInvalidLogin):
			return status.Wrap(errors.New("login is required"), status.InvalidArgument)
		case errors.Is(err, domain.ErrInvalidPassword):
			return status.Wrap(errors.New("password is required"), status.InvalidArgument)
		default:
			return status.Wrap(errors.New("unknown error"), status.Unknown)
		}
	}

	output = &api.EmptyBody{}
	return nil
}
