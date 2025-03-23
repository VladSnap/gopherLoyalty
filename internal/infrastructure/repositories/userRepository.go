package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserImplRepository struct {
	db *sqlx.DB
}

func NewUserImplRepository(db *sqlx.DB) *UserImplRepository {
	return &UserImplRepository{db: db}
}

func (r *UserImplRepository) Create(ctx context.Context, User dbModels.User) (string, error) {
	query := `INSERT INTO users (id, login, password) VALUES (:id, :login, :password)`
	_, err := r.db.NamedExecContext(ctx, query, User)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create user")
	}
	return User.ID, nil
}

func (r *UserImplRepository) FindByID(ctx context.Context, id string) (*dbModels.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var User dbModels.User
	err := r.db.GetContext(ctx, &User, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user by ID")
	}
	return &User, nil
}

func (r *UserImplRepository) FindByLoginAndPassword(ctx context.Context, login, password string) (*dbModels.User, error) {
	query := `SELECT * FROM users WHERE login = $1 AND password = $2`
	var User dbModels.User
	err := r.db.GetContext(ctx, &User, query, login, password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user with login %s not found", login)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user by login and password")
	}
	return &User, nil
}
