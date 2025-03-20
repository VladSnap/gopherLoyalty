package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserAccountImplRepository struct {
	db *sqlx.DB
}

func NewUserAccountImplRepository(db *sqlx.DB) *UserAccountImplRepository {
	return &UserAccountImplRepository{db: db}
}

func (r *UserAccountImplRepository) Create(ctx context.Context, userAccount dbModels.UserAccount) (string, error) {
	query := `INSERT INTO user_accounts (id, user_id, login, password) VALUES (:id, :user_id, :login, :password)`
	_, err := r.db.NamedExecContext(ctx, query, userAccount)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create user account")
	}
	return userAccount.ID, nil
}

func (r *UserAccountImplRepository) FindByID(ctx context.Context, id string) (*dbModels.UserAccount, error) {
	query := `SELECT * FROM user_accounts WHERE id = $1`
	var userAccount dbModels.UserAccount
	err := r.db.GetContext(ctx, &userAccount, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user account with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user account by ID")
	}
	return &userAccount, nil
}

func (r *UserAccountImplRepository) FindByLoginAndPassword(ctx context.Context, login, password string) (*dbModels.UserAccount, error) {
	query := `SELECT * FROM user_accounts WHERE login = $1 AND password = $2`
	var userAccount dbModels.UserAccount
	err := r.db.GetContext(ctx, &userAccount, query, login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user account with login %s not found", login)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user account by login and password")
	}
	return &userAccount, nil
}
