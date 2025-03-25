package getOrders

import (
	"context"
	"errors"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/api"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/log"
	"github.com/google/uuid"
	"github.com/swaggest/usecase/status"
)

// DBOrderRepository определяет методы для работы с таблицей orders без domain уровня для оптимизации.
type DBOrderRepository interface {
	FindByUserID(ctx context.Context, userID string) ([]dbModels.OrderGetDTO, error)
}

type GetOrdersUseCaseImpl struct {
	dbOrderRepo DBOrderRepository
}

func NewGetOrdersUseCase(dbOrderRepo DBOrderRepository) *GetOrdersUseCaseImpl {
	return &GetOrdersUseCaseImpl{dbOrderRepo: dbOrderRepo}
}

func (uc *GetOrdersUseCaseImpl) Execute(ctx context.Context, input *interface{}, output *OrderListResponse) error {
	currentUserID, ok := ctx.Value(api.KeyContext("UserID")).(uuid.UUID)
	if !ok {
		err := errors.New("current userID is empty")
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	orders, err := uc.dbOrderRepo.FindByUserID(ctx, currentUserID.String())
	if err != nil {
		log.Zap.Error(err)
		return status.Wrap(err, status.Unknown)
	}

	if len(orders) != 0 {
		for _, ord := range orders {
			*output = append(*output, OrderResponse{
				Number:     ord.Number,
				Status:     ord.Status,
				UploadedAt: ord.UploadedAt,
				Accrual: func() *float64 {
					if ord.Accrual == nil {
						return nil
					}
					acc := domain.CurrencyUnit(*ord.Accrual).ToMajorUnit()
					return &acc
				}(),
			})
		}
	}

	return nil
}
