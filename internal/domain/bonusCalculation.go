package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrInvalidLoyaltyStatus = errors.New("loyalty status is required")
	ErrInvalidAccrual       = errors.New("accrual must be non-negative")
)

// Представляет доменную модель расчета бонусов.
type BonusCalculation struct {
	id            uuid.UUID
	orderID       uuid.UUID
	loyaltyStatus LoyaltyStatus
	accrual       CurrencyUnit
}

// Создает новый расчет бонусов, если данные корректны.
func NewBonusCalculation(order Order) (*BonusCalculation, error) {
	return &BonusCalculation{
		id:            GenerateUniqueID(),
		orderID:       order.GetID(),
		loyaltyStatus: LoyaltyStatusRegistered,
		accrual:       CurrencyUnit(0),
	}, nil
}

// CreateBonusCalculationFromDB создает расчет бонусов из БД, игнорирует валидацию.
func CreateBonusCalculationFromDB(id, orderID uuid.UUID, loyaltyStatus string, accrual int) (*BonusCalculation, error) {
	status, err := ParseLoyaltyStatus(loyaltyStatus)
	if err != nil {
		return nil, fmt.Errorf("failed ParseOrderStatus: %w", err)
	}

	return &BonusCalculation{
		id:            id,
		orderID:       orderID,
		loyaltyStatus: status,
		accrual:       CurrencyUnit(accrual),
	}, nil
}

// Getters
func (bc *BonusCalculation) GetID() uuid.UUID {
	return bc.id
}

func (bc *BonusCalculation) GetOrderID() uuid.UUID {
	return bc.orderID
}

func (bc *BonusCalculation) GetLoyaltyStatus() LoyaltyStatus {
	return bc.loyaltyStatus
}

func (bc *BonusCalculation) GetAccrual() CurrencyUnit {
	return bc.accrual
}
