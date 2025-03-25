package withdrawBalance

import (
	"context"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

type WithdrawBalanceCmdHandlerImpl struct {
	userRepo  domain.UserRepository
	orderRepo domain.OrderRepository
	transRepo domain.WithdrawRepository
}

func NewWithdrawBalanceCmdHandler(userRepo domain.UserRepository,
	orderRepo domain.OrderRepository,
	transRepo domain.WithdrawRepository) *WithdrawBalanceCmdHandlerImpl {
	return &WithdrawBalanceCmdHandlerImpl{userRepo: userRepo, orderRepo: orderRepo, transRepo: transRepo}
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
	// Если заказ есть в системе, то проверим, чтобы он относился к текущему юзеру
	if order != nil {
		if order.GetUserID() != currentUser {
			return domain.ErrNotAuthorizeAccessOrder
		}
	}

	// Проверяем что, на балансе достаточно баллов (эту логику ниже надо убрать в доменный слой)
	balance, err := cmd.transRepo.CalcTotal(ctx, currentUser.String())
	if err != nil {
		return fmt.Errorf("failed calc balance: %w", err)
	}

	withdraw := domain.CurrencyFromMajorUnit(withdrawSum)
	if balance < withdraw {
		return domain.ErrInsufficientBalance
	}

	// Тут возможно надо будет проверить, не существует ли уже какая либо транзакция связанная с заказом
	// Не понятны требования к системе

	newWithdraw, err := domain.NewWithdraw(order.GetNumber(), currentUser, withdraw)
	if err != nil {
		return fmt.Errorf("failed create new withdraw: %w", err)
	}

	err = cmd.transRepo.Create(ctx, newWithdraw)
	if err != nil {
		return fmt.Errorf("failed save new order in DB: %w", err)
	}

	return nil
}
