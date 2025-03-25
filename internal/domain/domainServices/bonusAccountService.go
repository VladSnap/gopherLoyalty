package domainServices

import (
	"context"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

type BonusAccountServiceImpl struct {
	withdrawRepo  domain.WithdrawRepository
	bonusCalcRepo domain.BonusCalculationRepository
}

func NewBonusAccountServiceImpl(withdrawRepo domain.WithdrawRepository,
	bonusCalcRepo domain.BonusCalculationRepository) *BonusAccountServiceImpl {
	return &BonusAccountServiceImpl{withdrawRepo: withdrawRepo, bonusCalcRepo: bonusCalcRepo}
}

func (bs *BonusAccountServiceImpl) GetBonusAccount(userID uuid.UUID) {

}

func (bs *BonusAccountServiceImpl) GetBonusAccountState(ctx context.Context, userID uuid.UUID) (*domain.BonusAccountState, error) {
	withdrawTotal, err := bs.withdrawRepo.CalcTotal(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed withdrawRepo.CalcTotal: %w", err)
	}

	bonusTotal, err := bs.bonusCalcRepo.CalcTotal(ctx, userID.String())
	if err != nil {
		return nil, fmt.Errorf("failed bonusCalcRepo.CalcTotal: %w", err)
	}

	accountState, err := domain.NewBonusAccountState(bonusTotal-withdrawTotal, bonusTotal, withdrawTotal)
	if err != nil {
		return nil, fmt.Errorf("failed NewBonusAccountState: %w", err)
	}

	return accountState, nil
}
