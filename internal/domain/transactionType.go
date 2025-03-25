package domain

import (
	"fmt"
	"strings"
)

type TransactionType int

const (
	TransactionTypeAccrual TransactionType = iota
	TransactionTypeWithdraw
)

// Строковые представления для статусов
var tranStatusToString = map[TransactionType]string{
	TransactionTypeAccrual:  "ACCRUAL",
	TransactionTypeWithdraw: "WITHDRAW",
}

var stringToTranStatus = make(map[string]TransactionType)

func init() {
	// Инициализация обратного отображения
	for k, v := range tranStatusToString {
		stringToTranStatus[v] = k
	}
}

func (s TransactionType) String() string {
	if str, ok := tranStatusToString[s]; ok {
		return str
	}
	return "UNKNOWN"
}

func ParseTransactionType(str string) (TransactionType, error) {
	// Приводим строку к верхнему регистру для удобства сравнения
	str = strings.ToUpper(str)
	if status, ok := stringToTranStatus[str]; ok {
		return status, nil
	}
	return -1, fmt.Errorf("invalid status: %s", str)
}
