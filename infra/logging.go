package infra

import (
	"log/slog"
	"os"
)

func InitLogging(conf *Config) *slog.LevelVar {

	// this logger can be used to change runtime setting through web
	loggingLevel := new(slog.LevelVar)
	loggingLevel.Set(slog.LevelDebug) // default to debug
	level, err := StringToLevel(conf.LogLevel)
	if err == nil {
		loggingLevel.Set(level)
	}

	// set default logger level - TODO make this configurable
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggingLevel, // Set the default log level to DEBUG
	}))
	// Set the logger as the default
	slog.SetDefault(logger)

	return loggingLevel
}
