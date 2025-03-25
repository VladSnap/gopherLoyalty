package withdrawBalance

import (
	"context"
	"errors"
	"net/http"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
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
		err := errors.New("current userID is empty")
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	err := uc.cmdHandler.Execute(ctx, input.Order, input.Sum, currentUserID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidOrderNumber):
			log.Zap.Error(err)
			return rest.HTTPCodeAsError(http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrNotAuthorizeAccessOrder):
			log.Zap.Error(err)
			return rest.HTTPCodeAsError(http.StatusUnauthorized)
		case errors.Is(err, domain.ErrInsufficientBalance):
			log.Zap.Error(err)
			return rest.HTTPCodeAsError(http.StatusPaymentRequired)
		default:
			log.Zap.Error(err)
			return status.Wrap(errors.New("unknown error"), status.Unknown)
		}
	}

	return nil
}
