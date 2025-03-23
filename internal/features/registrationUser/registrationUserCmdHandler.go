package registrationUser

import (
	"context"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/domain/domainServices"
)

type RegistrationUserCmdHandlerImpl struct {
	repo        domain.UserRepository
	passService domainServices.PasswordService
}

func NewRegistrationUserCmdHandler(repo domain.UserRepository,
	passService domainServices.PasswordService) *RegistrationUserCmdHandlerImpl {
	return &RegistrationUserCmdHandlerImpl{repo: repo, passService: passService}
}

func (s *RegistrationUserCmdHandlerImpl) Execute(ctx context.Context, login string, password string) error {
	// Проверяем, что пользователь с таким логином не существует
	exists, err := s.repo.ExistsByLogin(ctx, login)
	if err != nil {
		return fmt.Errorf("failed check ExistsByLogin: %w", err)
	}
	if exists {
		return domain.ErrLoginAlreadyUsed
	}

	hashedPassword, err := s.passService.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed HashPassword: %w", err)
	}

	// Создаем нового пользователя
	user, err := domain.NewUser(domain.GenerateUniqueID(), login, hashedPassword)
	if err != nil {
		return err
	}

	// Сохраняем пользователя в репозитории
	return s.repo.Create(ctx, user)
}
