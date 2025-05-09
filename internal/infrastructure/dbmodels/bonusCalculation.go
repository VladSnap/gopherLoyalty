package dbmodels

import (
	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Представляет таблицу расчета бонусов.
type BonusCalculation struct {
	ID      string `db:"id"`       // uuid
	OrderID string `db:"order_id"` // -> order.ID (1 к 1)
	Accrual int    `db:"accrual"`  // int
}

// Преобразует DBBonusCalculation в доменную модель BonusCalculation.
func (dbc *BonusCalculation) ToDomain() (*domain.BonusCalculation, error) {
	id, err := uuid.Parse(dbc.ID)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.Parse(dbc.OrderID)
	if err != nil {
		return nil, err
	}

	return domain.CreateBonusCalculationFromDB(
		id,
		orderID,
		dbc.Accrual,
	)
}

// Преобразует доменную модель BonusCalculation в DBBonusCalculation.
func DBBonusCalculationFromDomain(bc *domain.BonusCalculation) *BonusCalculation {
	return &BonusCalculation{
		ID:      bc.GetID().String(),
		OrderID: bc.GetOrderID().String(),
		Accrual: int(bc.GetAccrual()),
	}
}
