package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New("local", ApiLog))
	assert.NotNil(t, New("dev", ApiLog))
	assert.NotNil(t, New("qa", ApiLog))
	assert.NotNil(t, New("prod", ApiLog))
}

func TestNewWithZap(t *testing.T) {
	zl, _ := zap.NewProduction()
	l := NewWithZap(zl)
	assert.NotNil(t, l)
}

func TestEncoder(t *testing.T) {
	assert.NotNil(t, encoder("local"))
	assert.NotNil(t, encoder("dev"))
	assert.NotNil(t, encoder("qa"))
	assert.NotNil(t, encoder("prod"))
}

func TestWriteSyncer(t *testing.T) {
	ws := newWriteSyncer(apiFileName, apiMaxSize, apiMaxBackup, apiMaxAge)
	assert.NotNil(t, ws)
}

func TestNewForTest(t *testing.T) {
	logger, entries := newForTest()
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
