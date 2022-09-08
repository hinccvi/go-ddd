package tools

import (
	"golang.org/x/crypto/bcrypt"
)

func Bcrypt(password string) (string, error) {
	// 12 bcrypt cost produce around ~300ms delay,
	// and this is the max delay that average users can tolerate
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost+2)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func BcryptCompare(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
