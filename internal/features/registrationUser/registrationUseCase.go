package registrationUser

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type RegistrationUserUseCaseImpl struct {
}

func NewRegistrationUserUseCase() *RegistrationUserUseCaseImpl {
	return &RegistrationUserUseCaseImpl{}
}

func (uc *RegistrationUserUseCaseImpl) Execute(ctx context.Context, input *RegisterRequest, output *api.EmptyBody) error {
	output = &api.EmptyBody{}
	return nil
}
