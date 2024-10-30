package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/handlers"
	"github.com/ilius/localsend-go/pkg/send"
	"github.com/ilius/localsend-go/pkg/startup"
)

func main() {
	defer func() {
		r := recover()
		slog.Error(fmt.Sprintf("%v", r))
	}()

	noColor := os.Getenv("NO_COLOLR") != ""
	setupLogger(noColor, defaultLevel)

	_flags := parseFlags()

	conf := config.Init()
	setupLoggerAfterConfigLoad(conf, noColor)
	handlers.SetConfig(conf)

	startup.StartupServices(conf)

	if _flags.ReceiveMode {
		slog.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	} else {
		err := send.SendFile(conf, _flags.ToDevice, _flags.FilePath)
		if err != nil {
			slog.Error("Send failed", "err", err)
		}
	}
}
