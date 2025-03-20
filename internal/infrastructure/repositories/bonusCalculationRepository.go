package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BonusCalculationImplRepository struct {
	db *sqlx.DB
}

func NewBonusCalculationImplRepository(db *sqlx.DB) *BonusCalculationImplRepository {
	return &BonusCalculationImplRepository{db: db}
}

func (r *BonusCalculationImplRepository) Create(ctx context.Context, bonusCalculation dbModels.BonusCalculation) (string, error) {
	query := `INSERT INTO bonus_calculations (id, order_id, loyalty_status, accrual) VALUES (:id, :order_id, :loyalty_status, :accrual)`
	_, err := r.db.NamedExecContext(ctx, query, bonusCalculation)
	if err != nil {
		return "", errors.Wrap(ErrDatabase, "failed to create bonus calculation")
	}
	return bonusCalculation.ID, nil
}

func (r *BonusCalculationImplRepository) FindByID(ctx context.Context, id string) (*dbModels.BonusCalculation, error) {
	query := `SELECT * FROM bonus_calculations WHERE id = $1`
	var bonusCalculation dbModels.BonusCalculation
	err := r.db.GetContext(ctx, &bonusCalculation, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "bonus calculation with id %s not found", id)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find bonus calculation by ID")
	}
	return &bonusCalculation, nil
}

func (r *BonusCalculationImplRepository) FindByOrderID(ctx context.Context, orderID string) (*dbModels.BonusCalculation, error) {
	query := `SELECT * FROM bonus_calculations WHERE order_id = $1`
	var bonusCalculation dbModels.BonusCalculation
	err := r.db.GetContext(ctx, &bonusCalculation, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "bonus calculation for order ID %s not found", orderID)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find bonus calculation by order ID")
	}
	return &bonusCalculation, nil
}
