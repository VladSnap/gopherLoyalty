package loginuser

import (
	"context"
	"errors"
	"fmt"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/domainservices"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/repositories"
)

type LoginUserCmdHandlerImpl struct {
	repo        domain.UserRepository
	passService domainservices.PasswordService
}

func NewRegistrationUserCmdHandler(repo domain.UserRepository,
	passService domainservices.PasswordService) *LoginUserCmdHandlerImpl {
	return &LoginUserCmdHandlerImpl{repo: repo, passService: passService}
}

func (s *LoginUserCmdHandlerImpl) Execute(ctx context.Context, login string, password string) (*domain.User, error) {
	// Ищем юзера по логину, он уникальный на всю таблицу.
	user, err := s.repo.FindByLogin(ctx, login)
	switch {
	case errors.Is(err, repositories.ErrNotFound):
		return nil, domain.ErrLoginOrPassword
	case err != nil:
		return nil, fmt.Errorf("failed check ExistsByLogin: %w", err)
	case user == nil:
		return nil, domain.ErrLoginOrPassword
	}

	// Если нашли юзера, то проверяем его пароль
	valid, err := s.passService.VerifyPassword(user.GetPassword(), password)
	if err != nil || !valid {
		return nil, errors.Join(domain.ErrLoginOrPassword, fmt.Errorf("failed VerifyPassword: %w", err))
	}

	return user, nil
}
