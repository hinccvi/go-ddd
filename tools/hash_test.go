package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		cypherText, err := Bcrypt(test.plainText)
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
		cypherText, err := Bcrypt(test.plainText)
		if assert.NoError(t, err) {
			assert.Nil(t, BcryptCompare(test.plainText, cypherText))
			assert.NotEmpty(t, cypherText)
		}
	}
}
