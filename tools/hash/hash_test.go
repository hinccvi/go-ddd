package tools

import (
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestBcrypt_WhenFail(t *testing.T) {
	tests := []struct {
		plainText string
		cost      int
	}{
		{plainText: "secret", cost: 50},
		{plainText: "MyNumberIs1995", cost: 100},
		{plainText: "specialSymbols!@#", cost: 200},
		{plainText: "hello world", cost: 200},
	}

	for _, test := range tests {
		cypherText, err := Bcrypt(test.plainText, test.cost)
		assert.NotNil(t, err)
		assert.Empty(t, cypherText)
	}
}

func TestBcryptCompare_WhenFail(t *testing.T) {
	tests := []struct {
		plainText  string
		cypherText string
	}{
		{plainText: "", cypherText: "secret"},
		{plainText: "secret", cypherText: ""},
	}

	for _, test := range tests {
		assert.NotNil(t, BcryptCompare(test.plainText, test.cypherText))
	}
}

func TestBcrypt_WhenSuccess(t *testing.T) {
	tests := []struct {
		plainText string
		err       error
	}{
		{plainText: "secret", err: nil},
		{plainText: "MyNumberIs1995", err: nil},
		{plainText: "specialSymbols!@#", err: nil},
		{plainText: "hello world", err: nil},
	}

	for _, test := range tests {
		cypherText, err := Bcrypt(test.plainText, constants.BcryptCost)
		assert.Equal(t, test.err, err)
		assert.NotEmpty(t, cypherText)
	}
}

func TestBcryptCompare_WhenSuccess(t *testing.T) {
	tests := []struct {
		plainText string
	}{
		{plainText: "secret"},
		{plainText: "MyNumberIs1995"},
		{plainText: "specialSymbols!@#"},
		{plainText: "hello world"},
	}

	for _, test := range tests {
		cypherText, err := Bcrypt(test.plainText, constants.BcryptCost)
		if assert.NoError(t, err) {
			assert.Nil(t, BcryptCompare(test.plainText, cypherText))
			assert.NotEmpty(t, cypherText)
		}
	}
}
