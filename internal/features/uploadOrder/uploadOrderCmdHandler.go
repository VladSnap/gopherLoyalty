package uploadOrder

import (
	"context"
	"fmt"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

type UploadOrderCmdHandlerImpl struct {
	userRepo      domain.UserRepository
	orderRepo     domain.OrderRepository
	bonusCalcRepo domain.BonusCalculationRepository
}

func NewUploadOrderCmdHandler(userRepo domain.UserRepository,
	orderRepo domain.OrderRepository,
	bonusCalcRepo domain.BonusCalculationRepository,
) *UploadOrderCmdHandlerImpl {
	return &UploadOrderCmdHandlerImpl{userRepo: userRepo, orderRepo: orderRepo, bonusCalcRepo: bonusCalcRepo}
}

func (cmd *UploadOrderCmdHandlerImpl) Execute(ctx context.Context, orderNumber string, currentUser uuid.UUID) error {
	if !domain.IsValidLuhn(orderNumber) {
		return domain.ErrInvalidOrderNumber
	}

	existOrder, err := cmd.orderRepo.FindByNumber(ctx, orderNumber)
	if err != nil {
		return fmt.Errorf("failed FindOrder OrderNumber=%s: %w", orderNumber, err)
	}

	if existOrder != nil {
		if existOrder.GetUserID() == currentUser {
			return domain.ErrAlreadyUploadedOrderCurrent
		} else {
			return domain.ErrAlreadyUploadedOrderAnother
		}
	}

	newOrder, err := domain.NewOrder(orderNumber, time.Now().UTC(), currentUser)
	if err != nil {
		return err
	}

	err = cmd.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return fmt.Errorf("failed save new order in DB: %w", err)
	}

	return nil // Означает, что новый заказ загружен без ошибок
}
