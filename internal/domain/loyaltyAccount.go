package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidBalance       = errors.New("balance cannot be negative")
	ErrInvalidWithdrawTotal = errors.New("withdraw total cannot be negative")
)

// Представляет доменную модель счета лояльности.
type LoyaltyAccount struct {
	id            uuid.UUID
	userID        uuid.UUID
	balance       CurrencyUnit
	withdrawTotal CurrencyUnit
}

// Создает новый счет лояльности, если данные корректны.
func NewLoyaltyAccount(id, userID uuid.UUID, balance, withdrawTotal CurrencyUnit) (*LoyaltyAccount, error) {
	if balance < 0 {
		return nil, ErrInvalidBalance
	}
	if withdrawTotal < 0 {
		return nil, ErrInvalidWithdrawTotal
	}

	return &LoyaltyAccount{
		id:            id,
		userID:        userID,
		balance:       balance,
		withdrawTotal: withdrawTotal,
	}, nil
}

// Getters
func (la *LoyaltyAccount) GetID() uuid.UUID {
	return la.id
}

func (la *LoyaltyAccount) GetUserID() uuid.UUID {
	return la.userID
}

func (la *LoyaltyAccount) GetBalance() CurrencyUnit {
	return la.balance
}

func (la *LoyaltyAccount) GetWithdrawTotal() CurrencyUnit {
	return la.withdrawTotal
}

// Setters
func (la *LoyaltyAccount) SetBalance(balance CurrencyUnit) error {
	if balance < 0 {
		return ErrInvalidBalance
	}
	la.balance = balance
	return nil
}

func (la *LoyaltyAccount) SetWithdrawTotal(withdrawTotal CurrencyUnit) error {
	if withdrawTotal < 0 {
		return ErrInvalidWithdrawTotal
	}
	la.withdrawTotal = withdrawTotal
	return nil
}
