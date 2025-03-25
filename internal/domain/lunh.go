package domain

import (
	"strconv"
	"strings"
)

// Функция для проверки валидности числа по алгоритму Луна
func IsValidLuhn(number string) bool {
	// Удаляем пробелы и тире из строки, если они есть
	number = strings.ReplaceAll(number, " ", "")
	number = strings.ReplaceAll(number, "-", "")

	// Проверяем, что строка содержит только цифры
	for _, char := range number {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Переменная для хранения суммы цифр
	sum := 0
	// Флаг для определения, нужно ли удваивать текущую цифру
	double := false

	// Проходим по цифрам числа справа налево
	for i := len(number) - 1; i >= 0; i-- {
		// Преобразуем символ в цифру
		digit, _ := strconv.Atoi(string(number[i]))

		// Если флаг double установлен, удваиваем цифру
		if double {
			digit *= 2
			// Если результат больше 9, вычитаем 9
			if digit > 9 {
				digit -= 9
			}
		}

		// Добавляем цифру к общей сумме
		sum += digit

		// Инвертируем флаг double
		double = !double
	}

	// Если сумма делится на 10 без остатка, число валидно
	return sum%10 == 0
}
