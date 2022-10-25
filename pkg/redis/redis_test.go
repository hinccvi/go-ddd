package redis

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	var flagMode = flag.String("mode", "local", "environment")

	flag.Parse()

	cfg, err := config.Load(*flagMode)
	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(config.Config{}, cfg))

	rds, err := Connect(context.TODO(), cfg)
	assert.Nil(t, err)
	assert.NotNil(t, rds)
}
