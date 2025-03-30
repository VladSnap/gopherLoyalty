package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidLogin     = errors.New("User: login is required")
	ErrInvalidPassword  = errors.New("User: password is required")
	ErrLoginAlreadyUsed = errors.New("User: login is already used")
	ErrLoginOrPassword  = errors.New("User: login or password invalid")
)

type User struct {
	id       uuid.UUID
	login    string
	password string
}

// NewUser создает нового пользователя, если данные корректны.
func NewUser(login string, password string) (*User, error) {
	if login == "" {
		return nil, ErrInvalidLogin
	}
	if password == "" {
		return nil, ErrInvalidPassword
	}
	return &User{
		id:       GenerateUniqueID(),
		login:    login,
		password: password,
	}, nil
}

// CreateUserFromDB создает пользователя из БД, игнорирует валидацию.
func CreateUserFromDB(id uuid.UUID, login string, password string) *User {
	return &User{
		id:       id,
		login:    login,
		password: password,
	}
}

// GetID возвращает ID пользователя.
func (u *User) GetID() uuid.UUID {
	return u.id
}

// GetLogin возвращает логин пользователя.
func (u *User) GetLogin() string {
	return u.login
}

// GetPassword возвращает Password пользователя.
func (u *User) GetPassword() string {
	return u.password
}
