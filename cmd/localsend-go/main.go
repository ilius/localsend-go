package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/discovery"
	"github.com/ilius/localsend-go/pkg/go-clipboard"
	"github.com/ilius/localsend-go/pkg/logging"
	"github.com/ilius/localsend-go/pkg/send"
	"github.com/ilius/localsend-go/pkg/server"
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

	if conf.Receive.Clipboard {
		clipboard.Init()
	}

	discovery.Start(conf) // Enable broadcast and monitoring functions

	if _flags.ReceiveMode {
		srv := server.New(conf)
		srv.StartHttpServer()
		slog.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	} else {
		err := send.SendFile(conf, _flags.ToDevice, _flags.FilePath)
		if err != nil {
			slog.Error("Send failed", "err", err)
		}
	}
}
