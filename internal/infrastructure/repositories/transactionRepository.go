package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TransactionImplRepository struct {
	db *sqlx.DB
}

func NewTransactionImplRepository(db *DatabaseLoyalty) *TransactionImplRepository {
	return &TransactionImplRepository{db: db.DB}
}

func (r *TransactionImplRepository) Create(ctx context.Context, transaction dbModels.Transaction) (string, error) {
	query := `INSERT INTO transactions (id, created_at, type, order_id) VALUES (:id, :created_at, :type, :order_id)`
	_, err := r.db.NamedExecContext(ctx, query, transaction)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create loyalty account transaction")
	}
	return transaction.ID, nil
}

func (r *TransactionImplRepository) FindByID(ctx context.Context, id string) (*dbModels.Transaction, error) {
	query := `SELECT * FROM transactions WHERE id = $1`
	var transaction dbModels.Transaction
	err := r.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account transaction with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transaction by ID")
	}
	return &transaction, nil
}

func (r *TransactionImplRepository) FindByOrderID(ctx context.Context, orderID string) ([]dbModels.Transaction, error) {
	query := `SELECT * FROM transactions WHERE order_id = $1`
	var transactions []dbModels.Transaction
	err := r.db.SelectContext(ctx, &transactions, query, orderID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transactions by order ID")
	}
	return transactions, nil
}

func (r *TransactionImplRepository) CalcBalanceAndWithdraw(ctx context.Context, userID string) (*dbModels.BalanceCalcDTO, error) {
	query := `SELECT
            SUM(CASE WHEN type = 'ACCRUAL' THEN amount ELSE -amount END) AS balance,
            SUM(CASE WHEN type = 'WITHDRAW' THEN amount ELSE 0 END) AS withdrawTotal
        FROM transactions lat
        JOIN orders o ON lat.order_id = o.id
        WHERE o.user_id = $1`
	var calc dbModels.BalanceCalcDTO
	err := r.db.GetContext(ctx, &calc, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account transaction with user_id %s not found", userID)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transaction by ID")
	}
	return &calc, nil
}
