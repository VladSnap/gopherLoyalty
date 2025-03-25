package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTransactionAmount = errors.New("invalid transaction amount")
	ErrInsufficientBalance      = errors.New("insufficient balance") // Это по идее надо в доменную модель счёта положить
)

// Представляет доменную модель транзакции лояльности.
type Withdraw struct {
	id          uuid.UUID
	createdAt   time.Time
	orderNumber string
	userID      uuid.UUID
	amount      CurrencyUnit
}

// Создает новую транзакцию, если данные корректны.
func NewWithdraw(
	orderNumber string,
	userID uuid.UUID,
	amount CurrencyUnit,
) (*Withdraw, error) {
	if orderNumber == "" || !IsValidLuhn(orderNumber) {
		return nil, ErrInvalidOrderNumber
	}
	if amount < 0 {
		return nil, ErrInvalidTransactionAmount
	}

	return &Withdraw{
		id:          GenerateUniqueID(),
		createdAt:   time.Now().UTC(),
		orderNumber: orderNumber,
		userID:      userID,
		amount:      amount,
	}, nil
}

// CreateTransactionFromDB создает транзацию из БД, игнорирует валидацию.
func CreateTransactionFromDB(
	id uuid.UUID,
	createdAt time.Time,
	orderNumber string,
	amount int,
) (*Withdraw, error) {
	return &Withdraw{
		id:          id,
		createdAt:   createdAt,
		orderNumber: orderNumber,
		amount:      CurrencyUnit(amount),
	}, nil
}

// Getters
func (w *Withdraw) GetID() uuid.UUID {
	return w.id
}

func (w *Withdraw) GetCreatedAt() time.Time {
	return w.createdAt
}

func (w *Withdraw) GetOrderNumber() string {
	return w.orderNumber
}

func (w *Withdraw) GetUserID() uuid.UUID {
	return w.userID
}

func (w *Withdraw) GetAmount() CurrencyUnit {
	return w.amount
}
