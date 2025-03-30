package getwithdrawals

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbmodels"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

type DBWithdrawRepository interface {
	DBFindByUserID(ctx context.Context, userID string) ([]dbmodels.Withdraw, error)
}

type GetWithdrawalsUseCaseImpl struct {
	withdrawRepo DBWithdrawRepository
}

func NewGetWithdrawalsUseCase(withdrawRepo DBWithdrawRepository) *GetWithdrawalsUseCaseImpl {
	return &GetWithdrawalsUseCaseImpl{withdrawRepo: withdrawRepo}
}

func (uc *GetWithdrawalsUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *WithdrawalListResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		return status.Wrap(errors.New("current userID is empty"), status.Unknown)
	}

	withdrawals, err := uc.withdrawRepo.DBFindByUserID(ctx, currentUserID.String())
	if err != nil {
		log.Zap.Errorf("failed FindByUserID: %v", err)
		return status.Wrap(err, status.Unknown)
	}

	if len(withdrawals) != 0 {
		for _, tr := range withdrawals {
			*output = append(*output, WithdrawalResponse{
				Order:       tr.OrderNumber,
				Sum:         domain.CurrencyUnit(tr.Amount).ToMajorUnit(),
				ProcessedAt: tr.CreatedAt,
			})
		}
	}

	return nil
}
