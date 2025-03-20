package dbModels

// BonusCalculation представляет таблицу расчета бонусов.
type BonusCalculation struct {
	ID            string `db:"id"`             // uuid
	OrderID       string `db:"order_id"`       // -> order.ID (1 к 1)
	LoyaltyStatus string `db:"loyalty_status"` // enum (NEW?, REGISTERED, INVALID, PROCESSING, PROCESSED)
	Accrual       int    `db:"accrual"`        // int
}
