package registrationuser

import (
	"context"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/domainservices"
)

type RegistrationUserCmdHandlerImpl struct {
	repo        domain.UserRepository
	passService domainservices.PasswordService
}

func NewRegistrationUserCmdHandler(repo domain.UserRepository,
	passService domainservices.PasswordService) *RegistrationUserCmdHandlerImpl {
	return &RegistrationUserCmdHandlerImpl{repo: repo, passService: passService}
}

func (s *RegistrationUserCmdHandlerImpl) Execute(ctx context.Context, login string, password string) (*domain.User, error) {
	// Проверяем, что пользователь с таким логином не существует
	exists, err := s.repo.ExistsByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("failed check ExistsByLogin: %w", err)
	}
	if exists {
		return nil, domain.ErrLoginAlreadyUsed
	}

	hashedPassword, err := s.passService.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed HashPassword: %w", err)
	}

	// Создаем нового пользователя
	user, err := domain.NewUser(login, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed create new user: %w", err)
	}

	// Сохраняем пользователя в репозитории
	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed save new user in DB: %w", err)
	}

	return user, nil
}
