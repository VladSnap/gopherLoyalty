package dbModels

import "time"

// Order представляет таблицу заказов.
type Order struct {
	ID                 string    `db:"id"`                             // uuid
	Number             string    `db:"number"`                         // string
	UploadedAt         time.Time `db:"uploaded_at"`                    // datetime
	UserID             string    `db:"user_id"`                        // -> user.ID (многие к 1)
	BonusCalculationID string    `db:"bonus_calculation_id,omitempty"` // -> bonusCalculation.ID (1 к 1)
	Status             string    `db:"status"`                         // enum (NEW, INVALID, PROCESSING, PROCESSED)
}
