package repositories

import (
	"context"
	"database/sql"
	"fmt"

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

func (r *WithdrawImplRepository) Create(ctx context.Context, withdraw *domain.Withdraw) error {
	dbTran := dbModels.DBWithdrawFromDomain(withdraw)
	query := `INSERT INTO withdraws (id, created_at, order_number, user_id, amount) 
	VALUES (:id, :created_at, :order_number, :user_id, :amount)`
	_, err := r.db.NamedExecContext(ctx, query, dbTran)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to create loyalty account withdraw")
	}
	return nil
}

func (r *WithdrawImplRepository) FindByID(ctx context.Context, id string) (*domain.Withdraw, error) {
	query := `SELECT * FROM withdraws WHERE id = $1`
	var Withdraw dbModels.Withdraw
	err := r.db.GetContext(ctx, &Withdraw, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account withdraw with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account withdraw by ID")
	}
	return Withdraw.ToDomain()
}

func (r *WithdrawImplRepository) DBFindByUserID(ctx context.Context, userID string) ([]dbModels.Withdraw, error) {
	query := `SELECT * FROM withdraws WHERE user_id = $1`
	var withdraws []dbModels.Withdraw
	err := r.db.SelectContext(ctx, &withdraws, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find withdraws by userID")
	}
	return withdraws, nil
}

func (r *WithdrawImplRepository) FindByUserID(ctx context.Context, userID string) ([]domain.Withdraw, error) {
	query := `SELECT * FROM withdraws WHERE user_id = $1`
	var withdraws []dbModels.Withdraw
	err := r.db.SelectContext(ctx, &withdraws, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find withdraws by userID")
	}
	return convertToDomWithdraw(withdraws)
}

func (r *WithdrawImplRepository) CalcTotal(ctx context.Context, userID string) (domain.CurrencyUnit, error) {
	var total sql.NullInt32
	query := `SELECT SUM(amount) AS total
        FROM withdraws
        WHERE user_id = $1`
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return domain.CurrencyUnit(0), errors.Wrap(ErrDatabase, "failed calc withdraw total")
	}
	return domain.CurrencyUnit(total.Int32), nil
}

func convertToDomWithdraw(dbWs []dbModels.Withdraw) ([]domain.Withdraw, error) {
	domWs := make([]domain.Withdraw, len(dbWs))
	for _, d := range dbWs {
		domW, err := d.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("fail convertToDomWithdraw: %w", err)
		}

		domWs = append(domWs, *domW)
	}

	return domWs, nil
}
