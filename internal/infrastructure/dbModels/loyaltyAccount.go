package dbModels

// LoyaltyAccount представляет таблицу счетов лояльности.
type LoyaltyAccount struct {
	ID            string `db:"id"`             // uuid
	UserID        string `db:"user_id"`        // -> user.ID (1 к 1)
	Balance       int    `db:"balance"`        // Хранит актуальное состояние счета
	WithdrawTotal int    `db:"withdraw_total"` // Хранит сколько всего было списано баллов
}
