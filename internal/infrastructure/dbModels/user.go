package dbModels

import (
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
)

// Представляет таблицу пользователей.
type User struct {
	ID       string `db:"id"`       // uuid
	Login    string `db:"login"`    // string
	Password string `db:"password"` // string
}

// Преобразует DBUser в доменную модель User.
func (dbu *User) ToDomain() (*domain.User, error) {
	id, err := domain.ParseUniqueID(dbu.ID)
	if err != nil {
		return nil, fmt.Errorf("failed ParseUniqueID: %w", err)
	}

	return domain.NewUser(id, dbu.Login, dbu.Password)
}

// Преобразует доменную модель User в DBUser.
func DbUserFromDomain(user *domain.User) *User {
	return &User{
		ID:       user.GetID().String(),
		Login:    user.GetLogin(),
		Password: user.GetPassword(),
	}
}
