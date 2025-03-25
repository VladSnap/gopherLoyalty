package withdrawBalance

import (
	"context"
	"errors"
	"net/http"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/google/uuid"
	"github.com/swaggest/rest"
	"github.com/swaggest/usecase/status"
)

type WithdrawBalanceUseCaseImpl struct {
	cmdHandler WithdrawBalanceCmdHandler
}

func NewWithdrawBalanceUseCase(cmdHandler WithdrawBalanceCmdHandler) *WithdrawBalanceUseCaseImpl {
	return &WithdrawBalanceUseCaseImpl{cmdHandler: cmdHandler}
}

type WithdrawBalanceCmdHandler interface {
	Execute(ctx context.Context, orderNumber string,
		withdrawSum float64, currentUser uuid.UUID) error
}

func (uc *WithdrawBalanceUseCaseImpl) Execute(ctx context.Context, input *WithdrawRequest, output *api.EmptyBody) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		return status.Wrap(errors.New("current userID is empty"), status.Unknown)
	}

	err := uc.cmdHandler.Execute(ctx, input.Order, input.Sum, currentUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOrderNumber):
			return rest.HTTPCodeAsError(http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrNotAuthorizeAccessOrder):
			return rest.HTTPCodeAsError(http.StatusUnauthorized)
		case errors.Is(err, domain.ErrInsufficientBalance):
			return rest.HTTPCodeAsError(http.StatusPaymentRequired)
		default:
			return status.Wrap(errors.New("unknown error"), status.Unknown)
		}
	}

	return nil
}
