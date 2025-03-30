package domain

import (
	"errors"
)

var (
	ErrInvalidBalance       = errors.New("balance cannot be negative")
	ErrInvalidBonusCalcTotal = errors.New("bonus calculation total cannot be negative")
	ErrInvalidWithdrawTotal = errors.New("withdraw total cannot be negative")
)

// Представляет доменную модель состояния счета лояльности.
type BonusAccountState struct {
	balance        CurrencyUnit
	bonusCalcTotal CurrencyUnit
	withdrawTotal  CurrencyUnit
}

// NewOrder создает состояние бонусного аккаунта, если данные корректны.
func NewBonusAccountState(balance, bonusCalcTotal, withdrawTotal CurrencyUnit) (
	*BonusAccountState, error) {
	if balance < 0 {
		return nil, ErrInvalidBalance
	}
	if bonusCalcTotal < 0 {
		return nil, ErrInvalidWithdrawTotal
	}
	if withdrawTotal < 0 {
		return nil, ErrInvalidWithdrawTotal
	}

	return &BonusAccountState{
		balance:        balance,
		bonusCalcTotal: bonusCalcTotal,
		withdrawTotal:  withdrawTotal,
	}, nil
}

func (ba *BonusAccountState) GetBalance() CurrencyUnit {
	return ba.balance
}

func (ba *BonusAccountState) GetBonusCalcTotal() CurrencyUnit {
	return ba.bonusCalcTotal
}

func (ba *BonusAccountState) GetWithdrawTotal() CurrencyUnit {
	return ba.withdrawTotal
}
