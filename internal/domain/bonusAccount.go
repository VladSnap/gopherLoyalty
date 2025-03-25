package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
)

// Представляет доменную модель счета лояльности.
type BonusAccount struct {
	userID     uuid.UUID
	bonusCalcs []BonusCalculation
	withdraws  []Withdraw
}

// Создает новый счет лояльности, если данные корректны.
func NewBonusAccount(userID uuid.UUID, bonusCalcs []BonusCalculation, withdraws []Withdraw) (
	*BonusAccount, error) {
	return &BonusAccount{
		userID:     userID,
		bonusCalcs: bonusCalcs,
		withdraws:  withdraws,
	}, nil
}

func (ba *BonusAccount) GetUserID() uuid.UUID {
	return ba.userID
}

func (ba *BonusAccount) GetState() (*BonusAccountState, error) {
	bonusCalcTotal := ba.getBonusCalcTotal()
	withdrawTotal := ba.getWithdrawTotal()
	return NewBonusAccountState(bonusCalcTotal-withdrawTotal, bonusCalcTotal, withdrawTotal)
}

func (ba *BonusAccount) GetWithdrawals() []Withdraw {
	return ba.withdraws
}

func (ba *BonusAccount) AddWithdraw(orderNumber string, amount CurrencyUnit) (*Withdraw, error) {
	withdraw, err := NewWithdraw(orderNumber, ba.userID, amount)
	if err != nil {
		return nil, err
	}

	state, err := ba.GetState()
	if err != nil {
		return nil, err
	}
	// Проверим достаточно ли средств на балансе
	if state.balance < amount {
		return nil, ErrInsufficientBalance
	}

	// Тут можно проверить, имеется ли списание для переданного номера заказа

	ba.withdraws = append(ba.withdraws, *withdraw)
	return withdraw, nil
}

func (ba *BonusAccount) getBonusCalcTotal() CurrencyUnit {
	var total CurrencyUnit
	for _, b := range ba.bonusCalcs {
		if b.GetLoyaltyStatus() == LoyaltyStatusProcessed {
			total = total + b.GetAccrual()
		}
	}
	return total
}

func (ba *BonusAccount) getWithdrawTotal() CurrencyUnit {
	var total CurrencyUnit
	for _, w := range ba.withdraws {
		total = total + w.GetAmount()
	}
	return total
}
