package tools

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrapRecursive(t *testing.T) {
	t.Run("success: no internal error", func(t *testing.T) {
		layer1Err := errors.New("layer 1 error")
		assert.Equal(t, layer1Err, UnwrapRecursive(layer1Err))
	})

	t.Run("success: 1 level nested error", func(t *testing.T) {
		layer2Err := errors.New("layer 2 error")
		layer1Err := fmt.Errorf("layer 1 error: %w", layer2Err)
		assert.Equal(t, layer2Err, UnwrapRecursive(layer1Err))
	})

	t.Run("success: 2 level nested error", func(t *testing.T) {
		layer3Err := errors.New("layer 3 error")
		layer2Err := fmt.Errorf("layer 2 error: %w", layer3Err)
		layer1Err := fmt.Errorf("layer 1 error: %w", layer2Err)
		assert.Equal(t, layer3Err, UnwrapRecursive(layer1Err))
	})
}
