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
	Create(ctx context.Context, bonusCalculation *BonusCalculation) error
	FindByOrderID(ctx context.Context, orderID string) (*BonusCalculation, error)
	FindByUserID(ctx context.Context, userID string) ([]BonusCalculation, error)
	CalcTotal(ctx context.Context, userID string) (CurrencyUnit, error)
}

// WithdrawRepository определяет методы для работы с таблицей withdraws.
type WithdrawRepository interface {
	Create(ctx context.Context, withdraw *Withdraw) error
	FindByID(ctx context.Context, id string) (*Withdraw, error)
	FindByUserID(ctx context.Context, userID string) ([]Withdraw, error)
	CalcTotal(ctx context.Context, userID string) (CurrencyUnit, error)
}
