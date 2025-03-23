package dbModels

import (
	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Представляет таблицу счетов лояльности.
type LoyaltyAccount struct {
	ID            string `db:"id"`             // uuid
	UserID        string `db:"user_id"`        // -> user.ID (1 к 1)
	Balance       int    `db:"balance"`        // Хранит актуальное состояние счета
	WithdrawTotal int    `db:"withdraw_total"` // Хранит сколько всего было списано баллов
}

// Преобразует DBLoyaltyAccount в доменную модель LoyaltyAccount.
func (dla *LoyaltyAccount) ToDomain() (*domain.LoyaltyAccount, error) {
	id, err := uuid.Parse(dla.ID)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(dla.UserID)
	if err != nil {
		return nil, err
	}

	return domain.NewLoyaltyAccount(
		id,
		userID,
		dla.Balance,
		dla.WithdrawTotal,
	)
}

// Преобразует доменную модель LoyaltyAccount в DBLoyaltyAccount.
func DBLoyaltyAccountFromDomain(la *domain.LoyaltyAccount) *LoyaltyAccount {
	return &LoyaltyAccount{
		ID:            la.GetID().String(),
		UserID:        la.GetUserID().String(),
		Balance:       la.GetBalance(),
		WithdrawTotal: la.GetWithdrawTotal(),
	}
}
