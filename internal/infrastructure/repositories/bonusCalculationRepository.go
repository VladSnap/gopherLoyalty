package repositories

import (
	"context"
	"database/sql"

	"github.com/VladSnap/gopherLoyalty/internal/domain"
	"github.com/VladSnap/gopherLoyalty/internal/infrastructure/dbModels"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type BonusCalculationImplRepository struct {
	db *sqlx.DB
}

func NewBonusCalculationImplRepository(db *DatabaseLoyalty) *BonusCalculationImplRepository {
	return &BonusCalculationImplRepository{db: db.DB}
}

func (r *BonusCalculationImplRepository) Create(ctx context.Context, bonusCalculation *domain.BonusCalculation) error {
	dbBonusCalc := dbModels.DBBonusCalculationFromDomain(bonusCalculation)
	query := `INSERT INTO bonus_calculations (id, order_id, loyalty_status, accrual) VALUES (:id, :order_id, :loyalty_status, :accrual)`
	_, err := r.db.NamedExecContext(ctx, query, dbBonusCalc)
	if err != nil {
		return errors.Wrap(ErrDatabase, "failed to create bonus calculation")
	}
	return nil
}

func (r *BonusCalculationImplRepository) FindByOrderID(ctx context.Context, orderID string) (*domain.BonusCalculation, error) {
	query := `SELECT * FROM bonus_calculations WHERE order_id = $1`
	var bonusCalculation dbModels.BonusCalculation
	err := r.db.GetContext(ctx, &bonusCalculation, query, orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(ErrNotFound, "bonus calculation for order ID %s not found", orderID)
		}
		return nil, errors.Wrap(ErrDatabase, "failed to find bonus calculation by order ID")
	}
	return bonusCalculation.ToDomain()
}

func (r *BonusCalculationImplRepository) CalcTotal(ctx context.Context, userID string) (domain.CurrencyUnit, error) {
	var total sql.NullInt32
	query := `SELECT SUM(b.accrual) AS total
        FROM bonus_calculations b
        JOIN orders o ON b.order_id = o.id
        WHERE o.user_id = $1 and loyalty_status = 'PROCESSED'`
	err := r.db.GetContext(ctx, &total, query, userID)
	if err != nil {
		return domain.CurrencyUnit(0), errors.Wrap(ErrDatabase, "failed to calc total")
	}

	return domain.CurrencyUnit(total.Int32), nil
}
