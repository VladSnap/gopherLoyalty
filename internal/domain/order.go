package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrderNumber = errors.New("order number is required")
	ErrInvalidUserID      = errors.New("user ID is required")
	ErrInvalidStatus      = errors.New("status is required")
)

// Order представляет доменную модель заказа.
type Order struct {
	id                 uuid.UUID
	number             string
	uploadedAt         time.Time
	userID             string
	bonusCalculationID *uuid.UUID // Опциональное поле
	status             string
}

// NewOrder создает новый заказ, если данные корректны.
func NewOrder(id uuid.UUID, number string, uploadedAt time.Time, userID, status string, bonusCalculationID *uuid.UUID) (*Order, error) {
	if number == "" {
		return nil, ErrInvalidOrderNumber
	}
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	if status == "" {
		return nil, ErrInvalidStatus
	}

	return &Order{
		id:                 id,
		number:             number,
		uploadedAt:         uploadedAt,
		userID:             userID,
		bonusCalculationID: bonusCalculationID,
		status:             status,
	}, nil
}

// Getters
func (o *Order) GetID() uuid.UUID {
	return o.id
}

func (o *Order) GetNumber() string {
	return o.number
}

func (o *Order) GetUploadedAt() time.Time {
	return o.uploadedAt
}

func (o *Order) GetUserID() string {
	return o.userID
}

func (o *Order) GetBonusCalculationID() *uuid.UUID {
	return o.bonusCalculationID
}

func (o *Order) GetStatus() string {
	return o.status
}

// Setters
func (o *Order) SetNumber(number string) error {
	if number == "" {
		return ErrInvalidOrderNumber
	}
	o.number = number
	return nil
}

func (o *Order) SetUploadedAt(uploadedAt time.Time) {
	o.uploadedAt = uploadedAt
}

func (o *Order) SetUserID(userID string) error {
	if userID == "" {
		return ErrInvalidUserID
	}
	o.userID = userID
	return nil
}

func (o *Order) SetBonusCalculationID(bonusCalculationID *uuid.UUID) {
	o.bonusCalculationID = bonusCalculationID
}

func (o *Order) SetStatus(status string) error {
	if status == "" {
		return ErrInvalidStatus
	}
	o.status = status
	return nil
}
