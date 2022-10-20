package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_Ok(t *testing.T) {
	tests := []struct {
		Env string
	}{
		{"local"},
		{"dev"},
		{"qa"},
		{"prod"},
	}

	for _, test := range tests {
		cfg, err := Load(test.Env)
		assert.NotNil(t, cfg)
		assert.Nil(t, err)
	}
}

func TestLoad_Fail(t *testing.T) {
	cfg, err := Load("test")
	assert.Equal(t, Config{}, cfg)
	assert.NotNil(t, err)
}
