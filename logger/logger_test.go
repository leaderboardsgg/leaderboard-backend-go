package logger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestInitLoggerSetsInfoLevelWhenDebugOff(t *testing.T) {
	debugOffConfig := LogConfig{
		Debug:   false,
		Console: true,
	}
	InitLogger(debugOffConfig)
	assert.Equal(t, zerolog.InfoLevel, Logger.GetLevel())
}

func TestInitLoggerSetsDebugLevelWhenDebugOn(t *testing.T) {
	debugOnConfig := LogConfig{
		Debug:   true,
		Console: true,
	}
	InitLogger(debugOnConfig)
	assert.Equal(t, zerolog.DebugLevel, Logger.GetLevel())
}
