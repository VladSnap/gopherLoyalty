package dbModels

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
	BonusCalculationID *string   `db:"bonus_calculation_id,omitempty"` // -> bonusCalculation.ID (1 к 1)
	Status             string    `db:"status"`                         // enum (NEW, INVALID, PROCESSING, PROCESSED)
}

// Преобразует DBOrder в доменную модель Order.
func (dbo *Order) ToDomain() (*domain.Order, error) {
	id, err := uuid.Parse(dbo.ID)
	if err != nil {
		return nil, fmt.Errorf("failed ParseUniqueID for ID: %w", err)
	}

	userId, err := domain.ParseUniqueID(dbo.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed ParseUniqueID for UserID: %w", err)
	}

	var bonusCalculationID *uuid.UUID
	if dbo.BonusCalculationID != nil {
		parsedUUID, err := uuid.Parse(*dbo.BonusCalculationID)
		if err != nil {
			return nil, err
		}
		bonusCalculationID = &parsedUUID
	}

	return domain.CreateOrderFromDB(
		id,
		dbo.Number,
		dbo.UploadedAt,
		userId,
		dbo.Status,
		bonusCalculationID,
	)
}

// Преобразует доменную модель Order в DBOrder.
func DBOrderFromDomain(order *domain.Order) *Order {
	var bonusCalculationID *string
	if order.GetBonusCalculationID() != nil {
		strUUID := order.GetBonusCalculationID().String()
		bonusCalculationID = &strUUID
	}

	return &Order{
		ID:                 order.GetID().String(),
		Number:             order.GetNumber(),
		UploadedAt:         order.GetUploadedAt(),
		UserID:             order.GetUserID().String(),
		BonusCalculationID: bonusCalculationID,
		Status:             order.GetStatus().String(),
	}
}
