package dbModels

// UserAccount представляет таблицу учетных записей пользователей.
type UserAccount struct {
	ID       string `db:"id"`       // uuid
	UserID   string `db:"user_id"`  // -> user.ID (1 к 1)
	Login    string `db:"login"`    // string
	Password string `db:"password"` // string
}
