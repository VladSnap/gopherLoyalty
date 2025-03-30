package repositories

import "github.com/pkg/errors"

// Определение пользовательских ошибок
var (
	ErrNotFound       = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input data")
	ErrDatabase       = errors.New("database error")
	ErrDuplicateEntry = errors.New("duplicate entry")
)
