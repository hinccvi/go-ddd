package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew_AccessLog(t *testing.T) {
	tests := []struct {
		env     string
		logType Type
	}{
		{env: "local", logType: AccessLog},
		{env: "dev", logType: AccessLog},
		{env: "qa", logType: AccessLog},
		{env: "prod", logType: AccessLog},
		{env: "local", logType: SQLLog},
		{env: "dev", logType: SQLLog},
		{env: "qa", logType: SQLLog},
		{env: "prod", logType: SQLLog},
		{env: "local", logType: ErrorLog},
		{env: "dev", logType: ErrorLog},
		{env: "qa", logType: ErrorLog},
		{env: "prod", logType: ErrorLog},
	}

	for _, test := range tests {
		assert.NotNil(t, New(test.env, test.logType))
	}
}

func TestNew_WithZap(t *testing.T) {
	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)
	assert.NotNil(t, l)
}

func TestNew_With(t *testing.T) {
	tests := []struct {
		arg interface{}
	}{
		{arg: nil},
		{arg: make([]interface{}, 0)},
		{arg: []interface{}{"key", "value"}},
		{arg: []interface{}{"count", 12}},
	}

	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)

	for _, test := range tests {
		assert.NotNil(t, l.With(context.TODO(), test.arg))
	}
}

func TestEncoder(t *testing.T) {
	assert.NotNil(t, encoder("local"))
	assert.NotNil(t, encoder("dev"))
	assert.NotNil(t, encoder("qa"))
	assert.NotNil(t, encoder("prod"))
}

func TestWriteSyncer(t *testing.T) {
	ws := newWriteSyncer(accessLogFileName, accessLogMaxSize, accessLogMaxBackup, accessLogMaxAge)
	assert.NotNil(t, ws)
}

func TestNewForTest(t *testing.T) {
	logger, entries := NewForTest()
	assert.NotNil(t, logger)
	assert.NotNil(t, entries)
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 1")
	assert.Equal(t, 1, entries.Len())
	logger.Info("msg 2")
	logger.Info("msg 3")
	assert.Equal(t, 3, entries.Len())
	entries.TakeAll()
	assert.Equal(t, 0, entries.Len())
	logger.Info("msg 4")
	assert.Equal(t, 1, entries.Len())
	assert.NotNil(t, logger)
	assert.NotNil(t, entries)
}
