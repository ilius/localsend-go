package logging

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/go-color"
	"github.com/ilius/localsend-go/pkg/slogcolor"
)

const DefaultLevel = slog.LevelInfo

func SetupLogger(noColor bool, level slog.Level) *slog.Logger {
	handler := slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
		Level:         level,
		TimeFormat:    time.DateTime,
		SrcFileMode:   slogcolor.ShortFile,
		SrcFileLength: 0,
		// MsgPrefix:     color.HiWhiteString("| "),
		MsgLength: 0,
		MsgColor:  color.New(),
		NoColor:   noColor,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

func parseLevel(levelStr string) (slog.Level, bool) {
	switch strings.ToLower(levelStr) {
	case "error":
		return slog.LevelError, true
	case "warn", "warning":
		return slog.LevelWarn, true
	case "info":
		return slog.LevelInfo, true
	case "debug":
		return slog.LevelDebug, true
	}
	return slog.LevelInfo, false
}

func SetupLoggerAfterConfigLoad(logger *slog.Logger, conf *config.Config, noColor bool) *slog.Logger {
	recreateLogger := false
	level := DefaultLevel
	if !noColor && conf.Logging.NoColor {
		noColor = true
		recreateLogger = true
	}
	if conf.Logging.Level != "" {
		configLevel, ok := parseLevel(conf.Logging.Level)
		if ok {
			if configLevel != DefaultLevel {
				level = configLevel
				recreateLogger = true
			}
		} else {
			slog.Error("invalid log level name", "level", conf.Logging.Level)
		}
	}
	if recreateLogger {
		slog.Info("Re-creating logger after loading config")
		return SetupLogger(noColor, level)
	}
	return logger
}
