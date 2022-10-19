package tools

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUUIDv4(t *testing.T) {
	u, err := GenerateUUIDv4()
	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(uuid.UUID{}, u))
}
