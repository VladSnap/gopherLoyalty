package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidAccrual = errors.New("accrual must be non-negative")
)

// Представляет доменную модель расчета бонусов.
type BonusCalculation struct {
	id      uuid.UUID
	orderID uuid.UUID
	accrual CurrencyUnit
}

// Создает новый расчет бонусов, если данные корректны.
func NewBonusCalculation(order *Order) (*BonusCalculation, error) {
	return &BonusCalculation{
		id:      GenerateUniqueID(),
		orderID: order.GetID(),
		accrual: CurrencyUnit(0),
	}, nil
}

// CreateBonusCalculationFromDB создает расчет бонусов из БД, игнорирует валидацию.
func CreateBonusCalculationFromDB(id, orderID uuid.UUID, accrual int) (*BonusCalculation, error) {
	return &BonusCalculation{
		id:      id,
		orderID: orderID,
		accrual: CurrencyUnit(accrual),
	}, nil
}

// Getters
func (bc *BonusCalculation) GetID() uuid.UUID {
	return bc.id
}

func (bc *BonusCalculation) GetOrderID() uuid.UUID {
	return bc.orderID
}

func (bc *BonusCalculation) GetAccrual() CurrencyUnit {
	return bc.accrual
}
