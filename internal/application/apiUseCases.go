package application

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/features/getbalance"
	"github.com/VladSnap/gopherLoyalty/internal/features/getorders"
	"github.com/VladSnap/gopherLoyalty/internal/features/getwithdrawals"
	"github.com/VladSnap/gopherLoyalty/internal/features/loginuser"
	"github.com/VladSnap/gopherLoyalty/internal/features/registrationuser"
	"github.com/VladSnap/gopherLoyalty/internal/features/uploadorder"
	"github.com/VladSnap/gopherLoyalty/internal/features/withdrawbalance"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
)

type RegistrationUserUseCase interface {
	Execute(ctx context.Context, input *registrationuser.RegisterUserRequest, output *registrationuser.RegisterUserResponse) error
}

type LoginUserUseCase interface {
	Execute(ctx context.Context, input *loginuser.LoginUserRequest, output *loginuser.LoginUserResponse) error
}

type UploadOrderUseCase interface {
	Execute(ctx context.Context, input uploadorder.UploadOrderRequest, output *api.HTTPStatusResponse) error
}

type GetOrdersUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getorders.OrderListResponse) error
}

type GetBalanceUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getbalance.BalanceResponse) error
}

type WithdrawBalanceUseCase interface {
	Execute(ctx context.Context, input *withdrawbalance.WithdrawRequest, output *api.EmptyBody) error
}

type GetWithdrawalsUseCase interface {
	Execute(ctx context.Context, input *interface{}, output *getwithdrawals.WithdrawalListResponse) error
}
