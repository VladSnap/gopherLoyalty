package getBalance

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/features/services"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

type GetBalanceUseCaseImpl struct {
	bonusAccountServ services.BonusAccountService
}

func NewGetBalanceUseCase(bonusAccountServ services.BonusAccountService) *GetBalanceUseCaseImpl {
	return &GetBalanceUseCaseImpl{bonusAccountServ: bonusAccountServ}
}

func (uc *GetBalanceUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *BalanceResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		err := errors.New("current userID is empty")
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	accountState, err := uc.bonusAccountServ.GetBonusAccountState(ctx, currentUserID)
	if err != nil {
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	output.Current = accountState.GetBalance().ToMajorUnit()
	output.Withdrawn = accountState.GetWithdrawTotal().ToMajorUnit()
	return nil
}
