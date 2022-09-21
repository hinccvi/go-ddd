package constants

import "golang.org/x/crypto/bcrypt"

const (
	// 12 bcrypt cost produce around ~300ms delay,
	// and this is the max delay that average users can tolerate.
	BcryptCost = bcrypt.DefaultCost + 2
)
