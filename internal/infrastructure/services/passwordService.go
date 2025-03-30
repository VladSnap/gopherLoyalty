package services

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"
)

var (
	ErrHashingFailed = errors.New("failed to hash password")
	ErrInvalidHash   = errors.New("invalid hash format")
)

// PasswordServiceImpl предоставляет методы для хэширования и проверки паролей.
type PasswordServiceImpl struct {
}

func NewPasswordServiceImpl() *PasswordServiceImpl {
	return &PasswordServiceImpl{}
}

// HashPassword хэширует пароль с использованием SHA-256.
func (ps *PasswordServiceImpl) HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", ErrHashingFailed
	}

	hash := sha256.Sum256(append(salt, []byte(password)...)) // Используем SHA-256
	encodedHash := base64.RawStdEncoding.EncodeToString(hash[:])
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)

	return encodedSalt + "$" + encodedHash, nil
}

// VerifyPassword проверяет, соответствует ли пароль хэшу.
func (ps *PasswordServiceImpl) VerifyPassword(storedHash, password string) (bool, error) {
	parts := strings.Split(storedHash, "$")
	if len(parts) != 2 {
		return false, ErrInvalidHash
	}

	encodedSalt, encodedHash := parts[0], parts[1]

	salt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false, ErrInvalidHash
	}

	hash := sha256.Sum256(append(salt, []byte(password)...)) // Используем SHA-256
	newEncodedHash := base64.RawStdEncoding.EncodeToString(hash[:])

	// Безопасное сравнение хэшей
	return subtle.ConstantTimeCompare([]byte(newEncodedHash), []byte(encodedHash)) == 1, nil
}
