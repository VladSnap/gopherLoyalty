package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type LoyaltyAccountTransactionImplRepository struct {
	db *sqlx.DB
}

func NewLoyaltyAccountTransactionImplRepository(db *sqlx.DB) *LoyaltyAccountTransactionImplRepository {
	return &LoyaltyAccountTransactionImplRepository{db: db}
}

func (r *LoyaltyAccountTransactionImplRepository) Create(ctx context.Context, transaction dbModels.LoyaltyAccountTransaction) (string, error) {
	query := `INSERT INTO loyalty_account_transactions (id, created_at, loyalty_account_id, transaction_type, order_id) VALUES (:id, :created_at, :loyalty_account_id, :transaction_type, :order_id)`
	_, err := r.db.NamedExecContext(ctx, query, transaction)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create loyalty account transaction")
	}
	return transaction.ID, nil
}

func (r *LoyaltyAccountTransactionImplRepository) FindByID(ctx context.Context, id string) (*dbModels.LoyaltyAccountTransaction, error) {
	query := `SELECT * FROM loyalty_account_transactions WHERE id = $1`
	var transaction dbModels.LoyaltyAccountTransaction
	err := r.db.GetContext(ctx, &transaction, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "loyalty account transaction with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transaction by ID")
	}
	return &transaction, nil
}

func (r *LoyaltyAccountTransactionImplRepository) FindByLoyaltyAccountID(ctx context.Context, loyaltyAccountID string) ([]dbModels.LoyaltyAccountTransaction, error) {
	query := `SELECT * FROM loyalty_account_transactions WHERE loyalty_account_id = $1`
	var transactions []dbModels.LoyaltyAccountTransaction
	err := r.db.SelectContext(ctx, &transactions, query, loyaltyAccountID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transactions by loyalty account ID")
	}
	return transactions, nil
}

func (r *LoyaltyAccountTransactionImplRepository) FindByOrderID(ctx context.Context, orderID string) ([]dbModels.LoyaltyAccountTransaction, error) {
	query := `SELECT * FROM loyalty_account_transactions WHERE order_id = $1`
	var transactions []dbModels.LoyaltyAccountTransaction
	err := r.db.SelectContext(ctx, &transactions, query, orderID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find loyalty account transactions by order ID")
	}
	return transactions, nil
}
