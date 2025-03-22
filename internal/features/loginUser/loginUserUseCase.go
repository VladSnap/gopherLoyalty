package loginUser

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type LoginUserUseCaseImpl struct {
}

func NewLoginUserUseCase() *LoginUserUseCaseImpl {
	return &LoginUserUseCaseImpl{}
}

func (uc *LoginUserUseCaseImpl) Execute(ctx context.Context, input *LoginUserRequest, output *api.EmptyBody) error {
	output = &api.EmptyBody{}
	return nil
}
