package domain

import (
	"fmt"
	"strings"
)

type LoyaltyStatus int

const (
	LoyaltyStatusRegistered LoyaltyStatus = iota
	LoyaltyStatusInvalid
	LoyaltyStatusProcessing
	LoyaltyStatusProcessed
)

// Строковые представления для статусов
var LoyaltyStatusToString = map[LoyaltyStatus]string{
	LoyaltyStatusRegistered: "REGISTERED",
	LoyaltyStatusInvalid:    "INVALID",
	LoyaltyStatusProcessing: "PROCESSING",
	LoyaltyStatusProcessed:  "PROCESSED",
}

var stringToLoyaltyStatus = make(map[string]LoyaltyStatus)

func init() {
	// Инициализация обратного отображения
	for k, v := range LoyaltyStatusToString {
		stringToLoyaltyStatus[v] = k
	}
}

func (s LoyaltyStatus) String() string {
	if str, ok := LoyaltyStatusToString[s]; ok {
		return str
	}
	return "UNKNOWN"
}

func ParseLoyaltyStatus(str string) (LoyaltyStatus, error) {
	// Приводим строку к верхнему регистру для удобства сравнения
	str = strings.ToUpper(str)
	if status, ok := stringToLoyaltyStatus[str]; ok {
		return status, nil
	}
	return -1, fmt.Errorf("invalid status: %s", str)
}
