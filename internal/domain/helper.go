package domain

import "github.com/google/uuid"

func GenerateUniqueID() uuid.UUID {
	return uuid.New()
}

func ParseUniqueID(uuidValue string) (uuid.UUID, error) {
	return uuid.Parse(uuidValue)
}
