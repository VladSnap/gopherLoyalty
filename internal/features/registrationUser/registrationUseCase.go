package registrationUser

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
	"github.com/swaggest/usecase/status"
)

type RegistrationUserUseCaseImpl struct {
	cmdHandler RegistrationUserCmdHandler
	jwtService services.JWTTokenService
}

func NewRegistrationUserUseCase(cmdHandler RegistrationUserCmdHandler,
	jwtService services.JWTTokenService) *RegistrationUserUseCaseImpl {
	return &RegistrationUserUseCaseImpl{cmdHandler: cmdHandler, jwtService: jwtService}
}

type RegistrationUserCmdHandler interface {
	Execute(ctx context.Context, login string, password string) (*domain.User, error)
}

func (uc *RegistrationUserUseCaseImpl) Execute(ctx context.Context, input *RegisterUserRequest, output *RegisterUserResponse) error {
	user, err := uc.cmdHandler.Execute(ctx, input.Login, input.Password)

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
	// Аутентифицируем пользователя присвоив токен авторизации заголовку Authorization в ответе.
	authToken, err := uc.jwtService.CreateToken(user.GetID())
	if err != nil {
		return fmt.Errorf("failed create auth token: %w", err)
	}
	output.Authorization = authToken
	return nil
}
