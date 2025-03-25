package getBalance

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/helpers"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

// DBTransactionRepository определяет методы для работы с таблицей orders без domain уровня для оптимизации.
type DBTransactionRepository interface {
	CalcBalanceAndWithdraw(ctx context.Context, userID string) (*dbModels.BalanceCalcDTO, error)
}

type GetBalanceUseCaseImpl struct {
	dbLoyaltyRepo DBTransactionRepository
}

func NewGetBalanceUseCase(dbLoyaltyRepo DBTransactionRepository) *GetBalanceUseCaseImpl {
	return &GetBalanceUseCaseImpl{dbLoyaltyRepo: dbLoyaltyRepo}
}

func (uc *GetBalanceUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *BalanceResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		err := errors.New("current userID is empty")
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	}

	calc, err := uc.dbLoyaltyRepo.CalcBalanceAndWithdraw(ctx, currentUserID.String())
	if err != nil {
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	output.Current = float64(domain.CurrencyUnit(helpers.GetOrDefaultInt(calc.Balance, 0)).ToMajorUnit())
	output.Withdrawn = float64(domain.CurrencyUnit(helpers.GetOrDefaultInt(calc.WithdrawTotal, 0)).ToMajorUnit())
	return nil
}
