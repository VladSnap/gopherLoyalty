package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

// Представляет доменную модель транзакции лояльности.
type Transaction struct {
	id               uuid.UUID
	createdAt        time.Time
	transactionType  string
	orderID          uuid.UUID
	amount           CurrencyUnit
}

// Создает новую транзакцию, если данные корректны.
func NewTransaction(
	id uuid.UUID,
	createdAt time.Time,
	transactionType string,
	orderID uuid.UUID,
	amount CurrencyUnit,
) (*Transaction, error) {
	if transactionType != "WITHDRAW" && transactionType != "ACCRUAL" {
		return nil, ErrInvalidTransactionType
	}

	return &Transaction{
		id:               id,
		createdAt:        createdAt,
		transactionType:  transactionType,
		orderID:          orderID,
		amount:           amount,
	}, nil
}

// Getters
func (lat *Transaction) GetID() uuid.UUID {
	return lat.id
}

func (lat *Transaction) GetCreatedAt() time.Time {
	return lat.createdAt
}

func (lat *Transaction) GetTransactionType() string {
	return lat.transactionType
}

func (lat *Transaction) GetOrderID() uuid.UUID {
	return lat.orderID
}

func (lat *Transaction) GetAmount() CurrencyUnit {
	return lat.amount
}

// Setters
func (lat *Transaction) SetTransactionType(transactionType string) error {
	if transactionType != "WITHDRAW" && transactionType != "ACCRUAL" {
		return ErrInvalidTransactionType
	}
	lat.transactionType = transactionType
	return nil
}

func (lat *Transaction) SetOrderID(orderID uuid.UUID) {
	lat.orderID = orderID
}

func (lat *Transaction) SetAmount(amount CurrencyUnit) {
	lat.amount = amount
}
