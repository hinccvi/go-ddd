package db

import (
	"flag"
	"reflect"
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/stretchr/testify/assert"
)

//nolint:gochecknoglobals // environment flag that only used in main
var flagMode = flag.String("env", "local", "environment")

func TestConnect(t *testing.T) {
	flag.Parse()

	cfg, err := config.Load(*flagMode)
	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(config.Config{}, cfg))

	zap := log.New(*flagMode, log.AccessLog)

	db, err := Connect(&cfg, zap)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func TestConnect_WhenConfigIsEmpty(t *testing.T) {
	zap := log.New(*flagMode, log.AccessLog)

	db, err := Connect(&config.Config{}, zap)
	assert.NotNil(t, err)
	assert.Nil(t, db)
}

func TestConnect_WhenInvalidDSN(t *testing.T) {
	zap := log.New(*flagMode, log.AccessLog)

	cfg, err := config.Load(*flagMode)
	cfg.Dsn = "xxx"

	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(cfg, zap))

	db, err := Connect(&cfg, zap)
	assert.NotNil(t, err)
	assert.Nil(t, db)
}
