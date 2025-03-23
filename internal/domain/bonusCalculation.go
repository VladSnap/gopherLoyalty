package domain

import (
	"errors"

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
	loyaltyStatus string
	accrual       int
}

// Создает новый расчет бонусов, если данные корректны.
func NewBonusCalculation(id, orderID uuid.UUID, loyaltyStatus string, accrual int) (*BonusCalculation, error) {
	if loyaltyStatus == "" {
		return nil, ErrInvalidLoyaltyStatus
	}
	if accrual < 0 {
		return nil, ErrInvalidAccrual
	}

	return &BonusCalculation{
		id:            id,
		orderID:       orderID,
		loyaltyStatus: loyaltyStatus,
		accrual:       accrual,
	}, nil
}

// Getters
func (bc *BonusCalculation) GetID() uuid.UUID {
	return bc.id
}

func (bc *BonusCalculation) GetOrderID() uuid.UUID {
	return bc.orderID
}

func (bc *BonusCalculation) GetLoyaltyStatus() string {
	return bc.loyaltyStatus
}

func (bc *BonusCalculation) GetAccrual() int {
	return bc.accrual
}

// Setters
func (bc *BonusCalculation) SetLoyaltyStatus(loyaltyStatus string) error {
	if loyaltyStatus == "" {
		return ErrInvalidLoyaltyStatus
	}
	bc.loyaltyStatus = loyaltyStatus
	return nil
}

func (bc *BonusCalculation) SetAccrual(accrual int) error {
	if accrual < 0 {
		return ErrInvalidAccrual
	}
	bc.accrual = accrual
	return nil
}
