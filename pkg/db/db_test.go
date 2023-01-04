package db

import (
	"context"
	"flag"
	"reflect"
	"testing"

	"github.com/hinccvi/go-ddd/internal/config"
	"github.com/stretchr/testify/assert"
)

//nolint:gochecknoglobals // environment flag that only used in main
var flagMode = flag.String("env", "local", "environment")

func TestConnect(t *testing.T) {
	flag.Parse()

	cfg, err := config.Load(*flagMode)
	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(config.Config{}, cfg))

	db, err := Connect(context.TODO(), &cfg)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestConnect_WhenConfigIsEmpty(t *testing.T) {
	db, err := Connect(context.TODO(), &config.Config{})
	assert.NotNil(t, err)
	assert.Nil(t, db)
}

func TestConnect_WhenInvalidDSN(t *testing.T) {
	cfg, err := config.Load(*flagMode)
	cfg.Dsn = "xxx"

	assert.Nil(t, err)

	db, err := Connect(context.TODO(), &cfg)
	assert.NotNil(t, err)
	assert.Nil(t, db)
}
