package getOrders

import (
	"context"
)

type GetOrdersUseCaseImpl struct {
}

func NewGetOrdersUseCase() *GetOrdersUseCaseImpl {
	return &GetOrdersUseCaseImpl{}
}

func (uc *GetOrdersUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *OrderListResponse) error {
	*output = append(*output, OrderResponse{Number: "1"})
	*output = append(*output, OrderResponse{Number: "2"})
	*output = append(*output, OrderResponse{Number: "3"})
	return nil
}
