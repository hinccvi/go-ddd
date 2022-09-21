package tools

import (
	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"golang.org/x/crypto/bcrypt"
)

func Bcrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), constants.BcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func BcryptCompare(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
