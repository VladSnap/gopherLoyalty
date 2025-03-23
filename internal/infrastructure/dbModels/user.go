package dbModels

// User представляет таблицу пользователей.
type User struct {
	ID       string `db:"id"`       // uuid
	Login    string `db:"login"`    // string
	Password string `db:"password"` // string
}
