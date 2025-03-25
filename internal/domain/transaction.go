package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTransactionType   = errors.New("invalid transaction type")
	ErrInvalidTransactionAmount = errors.New("invalid transaction amount")
	ErrInsufficientBalance      = errors.New("insufficient balance") // Это по идее надо в доменную модель счёта положить
)

// Представляет доменную модель транзакции лояльности.
type Transaction struct {
	id              uuid.UUID
	createdAt       time.Time
	transactionType TransactionType
	orderID         uuid.UUID
	amount          CurrencyUnit
}

// Создает новую транзакцию, если данные корректны.
func NewTransaction(
	transactionType TransactionType,
	orderID uuid.UUID,
	amount CurrencyUnit,
) (*Transaction, error) {
	if transactionType != TransactionTypeWithdraw &&
		transactionType != TransactionTypeAccrual {
		return nil, ErrInvalidTransactionType
	}

	if amount < 0 {
		return nil, ErrInvalidTransactionAmount
	}

	return &Transaction{
		id:              GenerateUniqueID(),
		createdAt:       time.Now().UTC(),
		transactionType: transactionType,
		orderID:         orderID,
		amount:          amount,
	}, nil
}

// CreateTransactionFromDB создает транзацию из БД, игнорирует валидацию.
func CreateTransactionFromDB(
	id uuid.UUID,
	createdAt time.Time,
	transactionType string,
	orderID uuid.UUID,
	amount int,
) (*Transaction, error) {
	trType, err := ParseTransactionType(transactionType)
	if err != nil {
		return nil, ErrInvalidTransactionType
	}

	return &Transaction{
		id:              id,
		createdAt:       createdAt,
		transactionType: trType,
		orderID:         orderID,
		amount:          CurrencyUnit(amount),
	}, nil
}

// Getters
func (lat *Transaction) GetID() uuid.UUID {
	return lat.id
}

func (lat *Transaction) GetCreatedAt() time.Time {
	return lat.createdAt
}

func (lat *Transaction) GetTransactionType() TransactionType {
	return lat.transactionType
}

func (lat *Transaction) GetOrderID() uuid.UUID {
	return lat.orderID
}

func (lat *Transaction) GetAmount() CurrencyUnit {
	return lat.amount
}
