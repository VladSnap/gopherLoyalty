package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbmodels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type UserImplRepository struct {
	db *sqlx.DB
}

func NewUserImplRepository(db *DatabaseLoyalty) *UserImplRepository {
	return &UserImplRepository{db: db.DB}
}

func (r *UserImplRepository) Create(ctx context.Context, user *domain.User) error {
	dbUser := dbmodels.DBUserFromDomain(user)
	query := `INSERT INTO users (id, login, password) VALUES (:id, :login, :password)`
	_, err := r.db.NamedExecContext(ctx, query, dbUser)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to create user")
	}
	return nil
}

func (r *UserImplRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var dbUser dbmodels.User
	err := r.db.GetContext(ctx, &dbUser, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user by ID")
	}
	return dbUser.ToDomain()
}

func (r *UserImplRepository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE login = $1`
	var dbUser dbmodels.User
	err := r.db.GetContext(ctx, &dbUser, query, login)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "user with login %s not found", login)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find user by login and password")
	}
	return dbUser.ToDomain()
}

func (r *UserImplRepository) ExistsByLogin(ctx context.Context, login string) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM users WHERE login = $1", login)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
