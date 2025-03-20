package dbModels

import "time"

// LoyaltyAccountTransaction представляет таблицу транзакций по счетам лояльности.
type LoyaltyAccountTransaction struct {
	ID               string    `db:"id"`                 // uuid
	CreatedAt        time.Time `db:"created_at"`         // datetime
	LoyaltyAccountID string    `db:"loyalty_account_id"` // -> loyaltyAccount.ID (многие к 1)
	TransactionType  string    `db:"transaction_type"`   // enum (withdraw, accrual)
	OrderID          string    `db:"order_id"`           // -> order.ID (многие к 1)
}
