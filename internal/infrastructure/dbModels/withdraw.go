package dbModels

import (
	"time"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

// Представляет таблицу транзакций по счетам лояльности.
type Withdraw struct {
	ID          string    `db:"id"`           // uuid
	CreatedAt   time.Time `db:"created_at"`   // datetime
	OrderNumber string    `db:"order_number"` // -> order.ID (многие к 1)
	UserID      string    `db:"user_id"`      // -> user.id (многие к 1)
	Amount      int       `db:"amount"`       // int
}

// Преобразует DBWithdraw в доменную модель Withdraw.
func (w *Withdraw) ToDomain() (*domain.Withdraw, error) {
	id, err := uuid.Parse(w.ID)
	if err != nil {
		return nil, err
	}

	return domain.CreateWithdrawFromDB(
		id,
		w.CreatedAt,
		w.OrderNumber,
		w.Amount,
	)
}

// Преобразует доменную модель Withdraw в DBWithdraw.
func DBWithdrawFromDomain(w *domain.Withdraw) *Withdraw {
	return &Withdraw{
		ID:          w.GetID().String(),
		CreatedAt:   w.GetCreatedAt(),
		OrderNumber: w.GetOrderNumber(),
		UserID:      w.GetUserID().String(),
		Amount:      int(w.GetAmount()),
	}
}
