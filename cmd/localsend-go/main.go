package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/go-clipboard"
	"github.com/ilius/localsend-go/pkg/handlers"
	"github.com/ilius/localsend-go/pkg/logging"
	"github.com/ilius/localsend-go/pkg/send"
	"github.com/ilius/localsend-go/pkg/startup"
)

func main() {
	defer func() {
		r := recover()
		if r != nil {
			slog.Error(fmt.Sprintf("%v", r))
		}
	}()

	noColor := os.Getenv("NO_COLOLR") != ""
	logging.SetupLogger(noColor, logging.DefaultLevel)

	_flags := parseFlags()

	conf := config.Init()
	logging.SetupLoggerAfterConfigLoad(conf, noColor)
	handlers.SetConfig(conf)

	if conf.Receive.Clipboard {
		clipboard.Init()
	}

	startup.StartDiscovery(conf) // Enable broadcast and monitoring functions

	if _flags.ReceiveMode {
		startup.StartHttpServer(conf)
		slog.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	} else {
		err := send.SendFile(conf, _flags.ToDevice, _flags.FilePath)
		if err != nil {
			slog.Error("Send failed", "err", err)
		}
	}
}
