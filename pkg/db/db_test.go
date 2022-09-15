package db

import (
	"flag"
	"reflect"
	"testing"

	"github.com/hinccvi/Golang-Project-Structure-Conventional/internal/config"
	"github.com/hinccvi/Golang-Project-Structure-Conventional/pkg/log"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	var flagMode = flag.String("mode", "local", "environment")

	flag.Parse()

	cfg, err := config.Load(*flagMode)
	assert.Nil(t, err)
	assert.False(t, reflect.DeepEqual(config.Config{}, cfg))

	zap := log.New(*flagMode, log.ApiLog)

	db, err := Connect(&cfg, zap)
	assert.NotNil(t, db)
	assert.Nil(t, err)
}
