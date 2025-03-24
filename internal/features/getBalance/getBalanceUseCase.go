package getBalance

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

// DBLoyaltyAccountTransactionRepository определяет методы для работы с таблицей orders без domain уровня для оптимизации.
type DBLoyaltyAccountTransactionRepository interface {
	CalcBalanceAndWithdraw(ctx context.Context, userID string) (*dbModels.LoyaltyAccountCalcDTO, error)
}

type GetBalanceUseCaseImpl struct {
	dbLoyaltyRepo DBLoyaltyAccountTransactionRepository
}

func NewGetBalanceUseCase(dbLoyaltyRepo DBLoyaltyAccountTransactionRepository) *GetBalanceUseCaseImpl {
	return &GetBalanceUseCaseImpl{dbLoyaltyRepo: dbLoyaltyRepo}
}

func (uc *GetBalanceUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *BalanceResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		return status.Wrap(errors.New("current userID is empty"), status.Unknown)
	}

	calc, err := uc.dbLoyaltyRepo.CalcBalanceAndWithdraw(ctx, currentUserID.String())
	if err != nil {
		return status.Wrap(err, status.Unknown)
	}

	output.Current = float64(domain.CurrencyUnit(calc.Balance).ToMajorUnit())
	output.Withdrawn = float64(domain.CurrencyUnit(calc.WithdrawTotal).ToMajorUnit())
	return nil
}
