package dbModels

// User представляет таблицу пользователей.
type User struct {
	ID            string `db:"id"`              // uuid
	UserAccountID string `db:"user_account_id"` // -> userAccount.ID (1 к 1)
}
