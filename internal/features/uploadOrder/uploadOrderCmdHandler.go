package uploadOrder

import (
	"context"
	"fmt"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

type UploadOrderCmdHandlerImpl struct {
	userRepo  domain.UserRepository
	orderRepo domain.OrderRepository
}

func NewUploadOrderCmdHandler(userRepo domain.UserRepository,
	orderRepo domain.OrderRepository) *UploadOrderCmdHandlerImpl {
	return &UploadOrderCmdHandlerImpl{userRepo: userRepo, orderRepo: orderRepo}
}

func (cmd *UploadOrderCmdHandlerImpl) Execute(ctx context.Context, orderNumber string, currentUser uuid.UUID) (bool, error) {
	if !domain.IsValidLuhn(orderNumber) {
		return false, domain.ErrInvalidOrderNumber
	}

	existOrder, err := cmd.orderRepo.FindByNumber(ctx, orderNumber)
	if err != nil {
		return false, fmt.Errorf("failed FindOrder OrderNumber=%s: %w", orderNumber, err)
	}

	if existOrder != nil {
		if existOrder.GetUserID() == currentUser {
			return false, domain.ErrAlreadyUploadedOrderCurrent
		} else {
			return false, domain.ErrAlreadyUploadedOrderAnother
		}
	}

	newOrder, err := domain.NewOrder(orderNumber, time.Now().UTC(), currentUser)
	if err != nil {
		return false, err
	}

	err = cmd.orderRepo.Create(ctx, newOrder)
	if err != nil {
		return false, fmt.Errorf("failed save new order in DB: %w", err)
	}

	return true, nil // Означает, что новый заказ загружен без ошибок
}
