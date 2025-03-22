package getBalance

import (
	"context"
)

type GetBalanceUseCaseImpl struct {
}

func NewGetBalanceUseCase() *GetBalanceUseCaseImpl {
	return &GetBalanceUseCaseImpl{}
}

func (uc *GetBalanceUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *BalanceResponse) error {
	output = &BalanceResponse{}
	return nil
}
