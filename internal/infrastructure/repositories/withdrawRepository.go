package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type WithdrawImplRepository struct {
	db *sqlx.DB
}

func NewWithdrawImplRepository(db *DatabaseLoyalty) *WithdrawImplRepository {
	return &WithdrawImplRepository{db: db.DB}
}

func (r *WithdrawImplRepository) Create(ctx context.Context, transaction *domain.Withdraw) error {
	dbTran := dbModels.DBWithdrawFromDomain(transaction)
	query := `INSERT INTO withdraws (id, created_at, type, order_number, amount) VALUES (:id, :created_at, :type, :order_number, :amount)`
	_, err := r.db.NamedExecContext(ctx, query, dbTran)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to create loyalty account transaction")
	}
	return nil
}

func (r *WithdrawImplRepository) FindByID(ctx context.Context, id string) (*domain.Withdraw, error) {
	query := `SELECT * FROM withdraws WHERE id = $1`
	var Withdraw dbModels.Withdraw
	err := r.db.GetContext(ctx, &Withdraw, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account transaction with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transaction by ID")
	}
	return Withdraw.ToDomain()
}

func (r *WithdrawImplRepository) FindByUserID(ctx context.Context, userID string) ([]dbModels.Withdraw, error) {
	query := `SELECT w.order_number, w.amount, w.created_at
	FROM withdraws w
	WHERE o.user_id = $1`
	var withdraws []dbModels.Withdraw
	err := r.db.SelectContext(ctx, &withdraws, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find withdraws by order userID")
	}
	return withdraws, nil
}

func (r *WithdrawImplRepository) CalcTotal(ctx context.Context, userID string) (domain.CurrencyUnit, error) {
	var total sql.NullInt32
	query := `SELECT SUM(amount) AS total
        FROM withdraws
        WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return domain.CurrencyUnit(0), errors.Wrap(ErrDatabase, "failed to calc total")
	}
	return domain.CurrencyUnit(total.Int32), nil
}
