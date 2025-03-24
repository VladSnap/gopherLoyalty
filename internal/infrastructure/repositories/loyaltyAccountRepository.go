package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type LoyaltyAccountImplRepository struct {
	db *sqlx.DB
}

func NewLoyaltyAccountImplRepository(db *DatabaseLoyalty) *LoyaltyAccountImplRepository {
	return &LoyaltyAccountImplRepository{db: db.DB}
}

func (r *LoyaltyAccountImplRepository) Create(ctx context.Context, loyaltyAccount dbModels.LoyaltyAccount) (string, error) {
	query := `INSERT INTO loyalty_accounts (id, user_id, balance, withdraw_total) VALUES (:id, :user_id, :balance, :withdraw_total)`
	_, err := r.db.NamedExecContext(ctx, query, loyaltyAccount)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create loyalty account")
	}
	return loyaltyAccount.ID, nil
}

func (r *LoyaltyAccountImplRepository) FindByUserID(ctx context.Context, userID string) (*dbModels.LoyaltyAccount, error) {
	query := `SELECT * FROM loyalty_accounts WHERE user_id = $1`
	var loyaltyAccount dbModels.LoyaltyAccount
	err := r.db.GetContext(ctx, &loyaltyAccount, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account for user ID %s not found", userID)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account by user ID")
	}
	return &loyaltyAccount, nil
}

func (r *LoyaltyAccountImplRepository) Update(ctx context.Context, loyaltyAccount dbModels.LoyaltyAccount) error {
	query := `UPDATE loyalty_accounts SET balance = :balance, withdraw_total = :withdraw_total WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, loyaltyAccount)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to update loyalty account")
	}
	return nil
}
