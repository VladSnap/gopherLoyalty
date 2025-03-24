package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTransactionType = errors.New("invalid transaction type")
)

// Представляет доменную модель транзакции счета лояльности.
type LoyaltyAccountTransaction struct {
	id               uuid.UUID
	createdAt        time.Time
	loyaltyAccountID uuid.UUID
	transactionType  string
	orderID          uuid.UUID
	amount           CurrencyUnit
}

// Создает новую транзакцию, если данные корректны.
func NewLoyaltyAccountTransaction(
	id uuid.UUID,
	createdAt time.Time,
	loyaltyAccountID uuid.UUID,
	transactionType string,
	orderID uuid.UUID,
	amount CurrencyUnit,
) (*LoyaltyAccountTransaction, error) {
	if transactionType != "withdraw" && transactionType != "accrual" {
		return nil, ErrInvalidTransactionType
	}

	return &LoyaltyAccountTransaction{
		id:               id,
		createdAt:        createdAt,
		loyaltyAccountID: loyaltyAccountID,
		transactionType:  transactionType,
		orderID:          orderID,
		amount:           amount,
	}, nil
}

// Getters
func (lat *LoyaltyAccountTransaction) GetID() uuid.UUID {
	return lat.id
}

func (lat *LoyaltyAccountTransaction) GetCreatedAt() time.Time {
	return lat.createdAt
}

func (lat *LoyaltyAccountTransaction) GetLoyaltyAccountID() uuid.UUID {
	return lat.loyaltyAccountID
}

func (lat *LoyaltyAccountTransaction) GetTransactionType() string {
	return lat.transactionType
}

func (lat *LoyaltyAccountTransaction) GetOrderID() uuid.UUID {
	return lat.orderID
}

func (lat *LoyaltyAccountTransaction) GetAmount() CurrencyUnit {
	return lat.amount
}

// Setters
func (lat *LoyaltyAccountTransaction) SetTransactionType(transactionType string) error {
	if transactionType != "withdraw" && transactionType != "accrual" {
		return ErrInvalidTransactionType
	}
	lat.transactionType = transactionType
	return nil
}

func (lat *LoyaltyAccountTransaction) SetOrderID(orderID uuid.UUID) {
	lat.orderID = orderID
}

func (lat *LoyaltyAccountTransaction) SetAmount(amount CurrencyUnit) {
	lat.amount = amount
}
