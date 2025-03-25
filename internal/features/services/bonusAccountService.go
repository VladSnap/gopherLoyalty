package services

import (
	"context"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/google/uuid"
)

type BonusAccountService interface {
	GetBonusAccount(ctx context.Context, userID uuid.UUID) (*domain.BonusAccount, error)
	GetBonusAccountState(ctx context.Context, userID uuid.UUID) (*domain.BonusAccountState, error)
}
