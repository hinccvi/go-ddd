package tools

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTestStatusCode(t *testing.T) {
	tests := []struct {
		code          int
		generatedCode int
	}{
		{code: http.StatusNotFound, generatedCode: http.StatusNotFound},
		{code: http.StatusInternalServerError, generatedCode: http.StatusInternalServerError},
		{code: http.StatusNetworkAuthenticationRequired, generatedCode: http.StatusNetworkAuthenticationRequired},
		{code: 9999, generatedCode: http.StatusBadRequest},
	}

	for _, test := range tests {
		assert.Equal(t, test.generatedCode, generateStatusCode(test.code))
	}
}
