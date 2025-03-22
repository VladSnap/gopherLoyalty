package withdrawBalance

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type WithdrawBalanceUseCaseImpl struct {
}

func NewWithdrawBalanceUseCase() *WithdrawBalanceUseCaseImpl {
	return &WithdrawBalanceUseCaseImpl{}
}

func (uc *WithdrawBalanceUseCaseImpl) Execute(ctx context.Context, input *WithdrawRequest, output *api.EmptyBody) error {
	output = &api.EmptyBody{}
	return nil
}
