package dbmodels

import (
	"fmt"
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Order представляет таблицу заказов.
type Order struct {
	ID                 string    `db:"id"`                             // uuid
	Number             string    `db:"number"`                         // string
	UploadedAt         time.Time `db:"uploaded_at"`                    // datetime
	UserID             string    `db:"user_id"`                        // -> user.ID (многие к 1)
	Status             string    `db:"status"`                         // enum (NEW, INVALID, PROCESSING, PROCESSED)
}

// OrderGetDTO модель для маппинга select+join запроса в DTO
type OrderGetDTO struct {
	Number     string    `db:"number"`
	UploadedAt time.Time `db:"uploaded_at"`
	Accrual    *int      `db:"accrual"`
	Status     string    `db:"status"`
}

// Преобразует DBOrder в доменную модель Order.
func (dbo *Order) ToDomain() (*domain.Order, error) {
	id, err := uuid.Parse(dbo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed ParseUniqueID for ID: %w", err)
	}

	userID, err := domain.ParseUniqueID(dbo.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed ParseUniqueID for UserID: %w", err)
	}

	return domain.CreateOrderFromDB(
		id,
		dbo.Number,
		dbo.UploadedAt,
		userID,
		dbo.Status,
	)
}

// Преобразует доменную модель Order в DBOrder.
func DBOrderFromDomain(order *domain.Order) *Order {
	return &Order{
		ID:                 order.GetID().String(),
		Number:             order.GetNumber(),
		UploadedAt:         order.GetUploadedAt(),
		UserID:             order.GetUserID().String(),
		Status:             order.GetStatus().String(),
	}
}

func (dbo *Order) ToGetDTO(accrual *int) *OrderGetDTO {
	return &OrderGetDTO{
		Number:     dbo.Number,
		UploadedAt: dbo.UploadedAt,
		Accrual:    accrual,
		Status:     dbo.Status,
	}
}
