package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrderNumber          = errors.New("order number is required")
	ErrInvalidOrderNotFound        = errors.New("order number not found")
	ErrNotAuthorizeAccessOrder     = errors.New("order not authorize access to order")
	ErrInvalidUserID               = errors.New("user ID is required")
	ErrInvalidStatus               = errors.New("status is required")
	ErrAlreadyUploadedOrderCurrent = errors.New("order already uploaded by current user")
	ErrAlreadyUploadedOrderAnother = errors.New("order already uploaded by another user")
	ErrInvalidBonusCalc            = errors.New("bonus calculation id does not match the order")
	ErrOrderChangeStatus           = errors.New("the current order status is not correct for changing it")
)

// Order представляет доменную модель заказа.
type Order struct {
	id                 uuid.UUID
	number             string
	uploadedAt         time.Time
	userID             uuid.UUID
	bonusCalculationID *uuid.UUID // Опциональное поле
	status             OrderStatus
}

// NewOrder создает новый заказ, если данные корректны.
func NewOrder(number string, uploadedAt time.Time,
	userID uuid.UUID) (*Order, error) {
	if number == "" || !IsValidLuhn(number) {
		return nil, ErrInvalidOrderNumber
	}
	if userID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	return &Order{
		id:                 GenerateUniqueID(),
		number:             number,
		uploadedAt:         uploadedAt,
		userID:             userID,
		bonusCalculationID: nil,
		status:             OrderStatusNew,
	}, nil
}

// CreateOrderFromDB создает заказ из БД, игнорирует валидацию.
func CreateOrderFromDB(id uuid.UUID, number string, uploadedAt time.Time,
	userID uuid.UUID, status string, bonusCalculationID *uuid.UUID) (*Order, error) {
	oStatus, err := ParseOrderStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed ParseOrderStatus: %w", err)
	}
	return &Order{
		id:                 id,
		number:             number,
		uploadedAt:         uploadedAt,
		userID:             userID,
		bonusCalculationID: bonusCalculationID,
		status:             oStatus,
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

func (o *Order) GetUserID() uuid.UUID {
	return o.userID
}

func (o *Order) GetBonusCalculationID() *uuid.UUID {
	return o.bonusCalculationID
}

func (o *Order) GetStatus() OrderStatus {
	return o.status
}

func (o *Order) MarkProcessed(bonCalc BonusCalculation) error {
	if o.status != OrderStatusNew && o.status != OrderStatusProcessing {
		return ErrOrderChangeStatus
	}

	if bonCalc.orderID != o.id {
		return ErrInvalidBonusCalc
	}

	o.status = OrderStatusProcessed
	return nil
}

func (o *Order) MarkProcessing() error {
	if o.status != OrderStatusNew {
		return ErrOrderChangeStatus
	}

	o.status = OrderStatusProcessing
	return nil
}

func (o *Order) MarkInvalid() error {
	if o.status != OrderStatusNew && o.status != OrderStatusProcessing {
		return ErrOrderChangeStatus
	}

	o.status = OrderStatusInvalid
	return nil
}
