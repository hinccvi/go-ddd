package tools

import "github.com/google/uuid"

func GenerateUUIDv4() (uuid.UUID, error) {
	u, err := uuid.NewRandom()

	return u, err
}
