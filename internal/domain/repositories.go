package domain

import (
	"context"
)

// UserAccountRepository определяет методы для работы с таблицей user_accounts.
type UserAccountRepository interface {
	Create(ctx context.Context, user User) (string, error)
	FindByID(ctx context.Context, id string) (*User, error)
	FindByLoginAndPassword(ctx context.Context, login, password string) (*User, error)
}

// OrderRepository определяет методы для работы с таблицей orders.
type OrderRepository interface {
	Create(ctx context.Context, order Order) (string, error)
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
	Update(ctx context.Context, order Order) error
}

// BonusCalculationRepository определяет методы для работы с таблицей bonus_calculations.
type BonusCalculationRepository interface {
	Create(ctx context.Context, bonusCalculation BonusCalculation) (string, error)
	FindByID(ctx context.Context, id string) (*BonusCalculation, error)
	FindByOrderID(ctx context.Context, orderID string) (*BonusCalculation, error)
}

// LoyaltyAccountRepository определяет методы для работы с таблицей loyalty_accounts.
type LoyaltyAccountRepository interface {
	Create(ctx context.Context, loyaltyAccount LoyaltyAccount) (string, error)
	FindByUserID(ctx context.Context, userID string) (*LoyaltyAccount, error)
	Update(ctx context.Context, loyaltyAccount LoyaltyAccount) error
}

// LoyaltyAccountTransactionRepository определяет методы для работы с таблицей loyalty_account_transactions.
type LoyaltyAccountTransactionRepository interface {
	Create(ctx context.Context, transaction LoyaltyAccountTransaction) (string, error)
	FindByID(ctx context.Context, id string) (*LoyaltyAccountTransaction, error)
	FindByLoyaltyAccountID(ctx context.Context, loyaltyAccountID string) ([]LoyaltyAccountTransaction, error)
	FindByOrderID(ctx context.Context, orderID string) ([]LoyaltyAccountTransaction, error)
}
