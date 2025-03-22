package getWithdrawals

import (
	"context"
)

type GetWithdrawalsUseCaseImpl struct {
}

func NewGetWithdrawalsUseCase() *GetWithdrawalsUseCaseImpl {
	return &GetWithdrawalsUseCaseImpl{}
}

func (uc *GetWithdrawalsUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *WithdrawalListResponse) error {
	*output = append(*output, WithdrawalResponse{Order: "1"})
	*output = append(*output, WithdrawalResponse{Order: "3"})
	*output = append(*output, WithdrawalResponse{Order: "5"})
	return nil
}
