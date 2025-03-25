package withdrawBalance

import (
	"context"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/features/services"
	"github.com/google/uuid"
)

type WithdrawBalanceCmdHandlerImpl struct {
	userRepo         domain.UserRepository
	orderRepo        domain.OrderRepository
	withdrawRepo     domain.WithdrawRepository
	bonusAccountServ services.BonusAccountService
}

func NewWithdrawBalanceCmdHandler(userRepo domain.UserRepository,
	orderRepo domain.OrderRepository,
	withdrawRepo domain.WithdrawRepository,
	bonusAccountServ services.BonusAccountService,
) *WithdrawBalanceCmdHandlerImpl {
	return &WithdrawBalanceCmdHandlerImpl{userRepo: userRepo, orderRepo: orderRepo,
		withdrawRepo: withdrawRepo, bonusAccountServ: bonusAccountServ}
}

func (cmd *WithdrawBalanceCmdHandlerImpl) Execute(ctx context.Context, orderNumber string,
	withdrawSum float64, currentUser uuid.UUID) error {
	if !domain.IsValidLuhn(orderNumber) {
		return domain.ErrInvalidOrderNumber
	}

	order, err := cmd.orderRepo.FindByNumber(ctx, orderNumber)
	if err != nil {
		return fmt.Errorf("failed FindOrder OrderNumber=%s: %w", orderNumber, err)
	}
	// Если заказ есть в системе, то проверим, чтобы он относился к текущему юзеру (не было в требовании)
	if order != nil {
		if order.GetUserID() != currentUser {
			return domain.ErrNotAuthorizeAccessOrder
		}
	}

	bonusAccount, err := cmd.bonusAccountServ.GetBonusAccount(ctx, currentUser)
	if err != nil {
		return fmt.Errorf("failed GetBonusAccount: %w", err)
	}

	newWithdraw, err := bonusAccount.AddWithdraw(orderNumber, domain.CurrencyFromMajorUnit(withdrawSum))
	if err != nil {
		return fmt.Errorf("failed AddWithdraw: %w", err)
	}

	err = cmd.withdrawRepo.Create(ctx, newWithdraw)
	if err != nil {
		return fmt.Errorf("failed save new order in DB: %w", err)
	}

	return nil
}
