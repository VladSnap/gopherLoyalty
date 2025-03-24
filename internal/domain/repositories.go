package domain

import (
	"context"
)

// UserRepository определяет методы для работы с таблицей user.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	ExistsByLogin(ctx context.Context, login string) (bool, error)
}

// OrderRepository определяет методы для работы с таблицей orders.
type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	FindByID(ctx context.Context, id string) (*Order, error)
	FindByNumber(ctx context.Context, number string) (*Order, error)
	Update(ctx context.Context, order *Order) error
}

// BonusCalculationRepository определяет методы для работы с таблицей bonus_calculations.
type BonusCalculationRepository interface {
	Create(ctx context.Context, bonusCalculation BonusCalculation) (string, error)
	FindByID(ctx context.Context, id string) (*BonusCalculation, error)
	FindByOrderID(ctx context.Context, orderID string) (*BonusCalculation, error)
}

// TransactionRepository определяет методы для работы с таблицей transactions.
type TransactionRepository interface {
	Create(ctx context.Context, transaction Transaction) (string, error)
	FindByID(ctx context.Context, id string) (*Transaction, error)
	FindByOrderID(ctx context.Context, orderID string) ([]Transaction, error)
}
