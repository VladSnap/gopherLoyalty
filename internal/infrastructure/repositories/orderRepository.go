package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type OrderImplRepository struct {
	db *sqlx.DB
}

func NewOrderImplRepository(db *sqlx.DB) *OrderImplRepository {
	return &OrderImplRepository{db: db}
}

func (r *OrderImplRepository) Create(ctx context.Context, order dbModels.Order) (string, error) {
	query := `INSERT INTO orders (id, number, uploaded_at, user_id, status) VALUES (:id, :number, :uploaded_at, :user_id, :status)`
	_, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create order")
	}
	return order.ID, nil
}

func (r *OrderImplRepository) FindByID(ctx context.Context, id string) (*dbModels.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1`
	var order dbModels.Order
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "order with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find order by ID")
	}
	return &order, nil
}

func (r *OrderImplRepository) FindByUserID(ctx context.Context, userID string) ([]dbModels.Order, error) {
	query := `SELECT * FROM orders WHERE user_id = $1`
	var orders []dbModels.Order
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find orders by user ID")
	}
	return orders, nil
}

func (r *OrderImplRepository) Update(ctx context.Context, order dbModels.Order) error {
	query := `UPDATE orders SET number = :number, uploaded_at = :uploaded_at, user_id = :user_id, status = :status WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, order)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to update order")
	}
	return nil
}
