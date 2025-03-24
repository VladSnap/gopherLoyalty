package dbModels

import (
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Представляет таблицу транзакций по счетам лояльности.
type LoyaltyAccountTransaction struct {
	ID               string    `db:"id"`                 // uuid
	CreatedAt        time.Time `db:"created_at"`         // datetime
	LoyaltyAccountID string    `db:"loyalty_account_id"` // -> loyaltyAccount.ID (многие к 1)
	TransactionType  string    `db:"transaction_type"`   // enum (withdraw, accrual)
	OrderID          string    `db:"order_id"`           // -> order.ID (многие к 1)
	Amount           int       `db:"amount"`             // int
}

// Преобразует DBLoyaltyAccountTransaction в доменную модель LoyaltyAccountTransaction.
func (dlat *LoyaltyAccountTransaction) ToDomain() (*domain.LoyaltyAccountTransaction, error) {
	id, err := uuid.Parse(dlat.ID)
	if err != nil {
		return nil, err
	}

	loyaltyAccountID, err := uuid.Parse(dlat.LoyaltyAccountID)
	if err != nil {
		return nil, err
	}

	orderID, err := uuid.Parse(dlat.OrderID)
	if err != nil {
		return nil, err
	}

	return domain.NewLoyaltyAccountTransaction(
		id,
		dlat.CreatedAt,
		loyaltyAccountID,
		dlat.TransactionType,
		orderID,
		domain.CurrencyUnit(dlat.Amount),
	)
}

// Преобразует доменную модель LoyaltyAccountTransaction в DBLoyaltyAccountTransaction.
func DBLoyaltyAccountTransactionFromDomain(lat *domain.LoyaltyAccountTransaction) *LoyaltyAccountTransaction {
	return &LoyaltyAccountTransaction{
		ID:               lat.GetID().String(),
		CreatedAt:        lat.GetCreatedAt(),
		LoyaltyAccountID: lat.GetLoyaltyAccountID().String(),
		TransactionType:  lat.GetTransactionType(),
		OrderID:          lat.GetOrderID().String(),
	}
}
