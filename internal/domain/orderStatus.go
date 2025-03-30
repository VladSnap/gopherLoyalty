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
var orderStatusToString = map[OrderStatus]string{
	OrderStatusNew:        "NEW",
	OrderStatusInvalid:    "INVALID",
	OrderStatusProcessing: "PROCESSING",
	OrderStatusProcessed:  "PROCESSED",
}

var stringToOrderStatus = make(map[string]OrderStatus)

func init() {
	// Инициализация обратного отображения
	for k, v := range orderStatusToString {
		stringToOrderStatus[v] = k
	}
}

func (s OrderStatus) String() string {
	if str, ok := orderStatusToString[s]; ok {
		return str
	}
	return "UNKNOWN"
}

func ParseOrderStatus(str string) (OrderStatus, error) {
	// Приводим строку к верхнему регистру для удобства сравнения
	str = strings.ToUpper(str)
	if status, ok := stringToOrderStatus[str]; ok {
		return status, nil
	}
	return -1, fmt.Errorf("invalid status: %s", str)
}
