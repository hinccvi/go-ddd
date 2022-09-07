package tools

import "github.com/google/uuid"

func GenerateUUIDv4() uuid.UUID {
	u, err := uuid.NewRandom()

	for err != nil {
		u, err = uuid.NewRandom()
	}

	return u
}
