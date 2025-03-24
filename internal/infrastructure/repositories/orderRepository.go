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

type OrderImplRepository struct {
	db *sqlx.DB
}

func NewOrderImplRepository(db *DatabaseLoyalty) *OrderImplRepository {
	return &OrderImplRepository{db: db.DB}
}

func (r *OrderImplRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `INSERT INTO orders (id, number, uploaded_at, user_id, status) VALUES (:id, :number, :uploaded_at, :user_id, :status)`
	dbOrder := dbModels.DBOrderFromDomain(order)
	_, err := r.db.NamedExecContext(ctx, query, dbOrder)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to create order")
	}
	return nil
}

func (r *OrderImplRepository) FindByID(ctx context.Context, id string) (*domain.Order, error) {
	query := `SELECT * FROM orders WHERE id = $1`
	var order dbModels.Order
	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "order with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find order by ID")
	}
	return order.ToDomain()
}

func (r *OrderImplRepository) FindByNumber(ctx context.Context, number string) (*domain.Order, error) {
	query := `SELECT * FROM orders WHERE number = $1`
	var order dbModels.Order
	err := r.db.GetContext(ctx, &order, query, number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find order by ID")
	}
	return order.ToDomain()
}

func (r *OrderImplRepository) Update(ctx context.Context, order *domain.Order) error {
	query := `UPDATE orders SET number = :number, uploaded_at = :uploaded_at, user_id = :user_id, status = :status WHERE id = :id`
	dbOrder := dbModels.DBOrderFromDomain(order)
	_, err := r.db.NamedExecContext(ctx, query, dbOrder)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to update order")
	}
	return nil
}

func convertOrders(dbOrders []dbModels.Order) ([]domain.Order, error) {
	domOrders := make([]domain.Order, len(dbOrders))
	for i, dbOrder := range dbOrders {
		domOrder, err := dbOrder.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed convert to domain order: %w", err)
		}

		domOrders[i] = *domOrder
	}
	return domOrders, nil
}

// DB repo implements

func (r *OrderImplRepository) FindByUserID(ctx context.Context, userID string) ([]dbModels.OrderGetDTO, error) {
	query := `SELECT o.number, o.uploaded_at, o.status, b.accrual
	FROM orders as o
	LEFT JOIN bonus_calculations as b ON o.id = b.order_id
	WHERE o.user_id = $1`
	var orders []dbModels.OrderGetDTO
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, errors.Wrap(ErrDatabase, "failed to find orders by user ID")
	}

	return orders, nil
}
