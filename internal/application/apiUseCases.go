package application

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/features/getBalance"
	"github.com/VladSnap/gopherLoyalty/internal/features/getOrders"
	"github.com/VladSnap/gopherLoyalty/internal/features/getWithdrawals"
	"github.com/VladSnap/gopherLoyalty/internal/features/loginUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/registrationUser"
	"github.com/VladSnap/gopherLoyalty/internal/features/withdrawBalance"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type RegistrationUserUseCase interface {
	Execute(ctx context.Context, input *registrationUser.RegisterUserRequest, output *registrationUser.RegisterUserResponse) error
}

type LoginUserUseCase interface {
	Execute(ctx context.Context, input *loginUser.LoginUserRequest, output *loginUser.LoginUserResponse) error
}

type UploadOrderUseCase interface {
	Execute(ctx context.Context, input string, output *api.EmptyBody) error
}

type GetOrdersUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getOrders.OrderListResponse) error
}

type GetBalanceUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getBalance.BalanceResponse) error
}

type WithdrawBalanceUseCase interface {
	Execute(ctx context.Context, input *withdrawBalance.WithdrawRequest, output *api.EmptyBody) error
}

type GetWithdrawalsUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getWithdrawals.WithdrawalListResponse) error
}
