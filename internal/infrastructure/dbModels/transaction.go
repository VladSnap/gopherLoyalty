package dbModels

import (
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Представляет таблицу транзакций по счетам лояльности.
type Transaction struct {
	ID              string    `db:"id"`         // uuid
	CreatedAt       time.Time `db:"created_at"` // datetime
	TransactionType string    `db:"type"`       // enum (withdraw, accrual)
	OrderID         string    `db:"order_id"`   // -> order.ID (многие к 1)
	Amount          int       `db:"amount"`     // int
}

type TransactionDTO struct {
	OrderNumber string    `db:"order_number"`
	Amount      int       `db:"amount"`
	CreatedAt   time.Time `db:"created_at"`
}

// Представляет DTO для расчета баланса и списаний.
type BalanceCalcDTO struct {
	Balance       *int `db:"balance"`
	WithdrawTotal *int `db:"withdrawtotal"`
}

// Преобразует DBTransaction в доменную модель Transaction.
func (dlat *Transaction) ToDomain() (*domain.Transaction, error) {
	id, err := uuid.Parse(dlat.ID)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.Parse(dlat.OrderID)
	if err != nil {
		return nil, err
	}

	return domain.NewTransaction(
		id,
		dlat.CreatedAt,
		dlat.TransactionType,
		orderID,
		domain.CurrencyUnit(dlat.Amount),
	)
}

// Преобразует доменную модель Transaction в DBTransaction.
func DBTransactionFromDomain(lat *domain.Transaction) *Transaction {
	return &Transaction{
		ID:              lat.GetID().String(),
		CreatedAt:       lat.GetCreatedAt(),
		TransactionType: lat.GetTransactionType(),
		OrderID:         lat.GetOrderID().String(),
	}
}
