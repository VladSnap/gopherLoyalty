package domain

import (
	"fmt"
	"strings"
)

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusInvalid
	OrderStatusProcessing
	OrderStatusProcessed
)

// Строковые представления для статусов
var statusToString = map[OrderStatus]string{
	OrderStatusNew:        "NEW",
	OrderStatusInvalid:    "INVALID",
	OrderStatusProcessing: "PROCESSING",
	OrderStatusProcessed:  "PROCESSED",
}

var stringToStatus = make(map[string]OrderStatus)

func init() {
	// Инициализация обратного отображения
	for k, v := range statusToString {
		stringToStatus[v] = k
	}
}

func (s OrderStatus) String() string {
	if str, ok := statusToString[s]; ok {
		return str
	}
	return "UNKNOWN"
}

func ParseOrderStatus(str string) (OrderStatus, error) {
	// Приводим строку к верхнему регистру для удобства сравнения
	str = strings.ToUpper(str)
	if status, ok := stringToStatus[str]; ok {
		return status, nil
	}
	return -1, fmt.Errorf("invalid status: %s", str)
}
