package tools

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// 12 bcrypt cost produce around ~300ms delay,
	// and this is the max delay that average users can tolerate.
	bcryptCost = bcrypt.DefaultCost + 2
)

func Bcrypt(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func BcryptCompare(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
