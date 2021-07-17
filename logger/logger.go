package logger

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

type LogConfig struct {
	Console bool
	Debug   bool
}

// InitLogger initializes the exported logger using
// the expected environment variables.
func InitLogger(config LogConfig) {
	var logLevel zerolog.Level
	if config.Debug {
		logLevel = zerolog.DebugLevel
	} else {
		logLevel = zerolog.InfoLevel
	}

	if config.Console {
		Logger = zerolog.New(zerolog.NewConsoleWriter()).
			Level(logLevel).
			With().
			Timestamp().
			Logger()
	} else {
		Logger = zerolog.New(defaultLumberjackConfig()).
			Level(logLevel).
			With().
			Timestamp().
			Logger()
	}
}

func defaultLumberjackConfig() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   "./logs/gql_server.log",
		MaxSize:    500, // Megabytes
		MaxBackups: 3,
		MaxAge:     28, // Days
	}
}
