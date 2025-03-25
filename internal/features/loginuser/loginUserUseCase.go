package loginuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/services"
	"github.com/swaggest/usecase/status"
)

type LoginUserUseCaseImpl struct {
	cmdHandler LoginUserCmdHandler
	jwtService services.JWTTokenService
}

func NewLoginUserUseCase(cmdHandler LoginUserCmdHandler,
	jwtService services.JWTTokenService) *LoginUserUseCaseImpl {
	return &LoginUserUseCaseImpl{cmdHandler: cmdHandler, jwtService: jwtService}
}

type LoginUserCmdHandler interface {
	Execute(ctx context.Context, login string, password string) (*domain.User, error)
}

func (uc *LoginUserUseCaseImpl) Execute(ctx context.Context, input *LoginUserRequest, output *LoginUserResponse) error {
	user, err := uc.cmdHandler.Execute(ctx, input.Login, input.Password)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrLoginOrPassword):
			return status.Wrap(errors.New("login or password incorrect"), status.Unauthenticated)
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
