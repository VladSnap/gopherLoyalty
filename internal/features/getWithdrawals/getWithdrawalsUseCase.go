package getWithdrawals

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

type DBTransactionRepository interface {
	FindWithdrawalByUserID(ctx context.Context, userID string) ([]dbModels.TransactionDTO, error)
}

type GetWithdrawalsUseCaseImpl struct {
	transactRepo DBTransactionRepository
}

func NewGetWithdrawalsUseCase(transactRepo DBTransactionRepository) *GetWithdrawalsUseCaseImpl {
	return &GetWithdrawalsUseCaseImpl{transactRepo: transactRepo}
}

func (uc *GetWithdrawalsUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *WithdrawalListResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		return status.Wrap(errors.New("current userID is empty"), status.Unknown)
	}

	trans, err := uc.transactRepo.FindWithdrawalByUserID(ctx, currentUserID.String())
	if err != nil {
		return status.Wrap(err, status.Unknown)
	}

	if len(trans) != 0 {
		for _, tr := range trans {
			*output = append(*output, WithdrawalResponse{
				Order:       tr.OrderNumber,
				Sum:         domain.CurrencyUnit(tr.Amount).ToMajorUnit(),
				ProcessedAt: tr.CreatedAt,
			})
		}
	}

	return nil
}
