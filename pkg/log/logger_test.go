package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New("local", zap.InfoLevel))
	assert.NotNil(t, New("dev", zap.ErrorLevel))
	assert.NotNil(t, New("qa", zap.ErrorLevel))
	assert.NotNil(t, New("prod", zap.ErrorLevel))
}

func TestNewWithZap(t *testing.T) {
	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)
	assert.NotNil(t, l)
}

func TestEncoder(t *testing.T) {
	assert.NotNil(t, Encoder("local"))
	assert.NotNil(t, Encoder("dev"))
	assert.NotNil(t, Encoder("qa"))
	assert.NotNil(t, Encoder("prod"))
}

func TestWriteSyncer(t *testing.T) {
	ws := WriteSyncer()
	assert.NotNil(t, ws)
}

func TestNewForTest(t *testing.T) {
	logger, entries := NewForTest()
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
}
