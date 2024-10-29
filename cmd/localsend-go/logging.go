package main

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/go-color"
	"github.com/ilius/localsend-go/pkg/slogcolor"
)

const defaultLevel = slog.LevelInfo

func setupLogger(noColor bool, level slog.Level) {
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
	slog.SetDefault(slog.New(handler))
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

func setupLoggerAfterConfigLoad(noColor bool) {
	recreateLogger := false
	level := defaultLevel
	if !noColor && config.Global.Logging.NoColor {
		noColor = true
		recreateLogger = true
	}
	if config.Global.Logging.Level != "" {
		configLevel, ok := parseLevel(config.Global.Logging.Level)
		if ok {
			if configLevel != defaultLevel {
				level = configLevel
				recreateLogger = true
			}
		} else {
			slog.Error("invalid log level name", "level", config.Global.Logging.Level)
		}
	}
	if recreateLogger {
		slog.Info("Re-creating logger after loading config")
		setupLogger(noColor, level)
	}
}
