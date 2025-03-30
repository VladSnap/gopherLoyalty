package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidJWT = errors.New("JWT token not valided")
var ErrInvalidJWTClaims = errors.New("JWT token claims invalid")

// Ключ для подписи JWT. Для упрощения храним тут, но в продакшен реди это надо хранить в хранилише секретов
const SecretKey = ":gn8(&^g80apoOIHf09yDFP(7GSFAP(una[wiof[9*EHFgoafhjlk';l;ojfsdiuf]*ioubHA&^uianfio"

// Структура для хранения данных в JWT (например, user_id)
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTTokenService interface {
	CreateToken(userID uuid.UUID) (string, error)
	GetAndValidateToken(tokenString string) (*Claims, error)
}

// Сервис для работы с JWT.
type JWTTokenServiceImpl struct{}

// Создает новый экземпляр сервиса.
func NewJWTTokenService() *JWTTokenServiceImpl {
	return &JWTTokenServiceImpl{}
}

// Создает JWT токен с указанным user_id.
func (s *JWTTokenServiceImpl) CreateToken(userID uuid.UUID) (string, error) {
	// Создаем claims с user_id и временем действия токена.
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Токен действителен 24 часа.
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user-authentication",
		},
	}

	// Создаем токен с использованием HMAC алгоритма подписи.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// Проверяет валидность токена и возвращает данные пользователя.
func (s *JWTTokenServiceImpl) GetAndValidateToken(tokenString string) (*Claims, error) {
	// Парсим токен с проверкой подписи.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется правильный алгоритм подписи.
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, errors.Join(ErrInvalidJWT, fmt.Errorf("invalid token: %w", err))
	}

	// Извлекаем claims из токена.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidJWTClaims
}
