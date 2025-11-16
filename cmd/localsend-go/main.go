package main

import (
	"fmt"
	"os"

	"codeberg.org/ilius/localsend-go/pkg/config"
	"codeberg.org/ilius/localsend-go/pkg/discovery"
	"codeberg.org/ilius/localsend-go/pkg/go-clipboard"
	"codeberg.org/ilius/localsend-go/pkg/logging"
	"codeberg.org/ilius/localsend-go/pkg/send"
	"codeberg.org/ilius/localsend-go/pkg/server"
)

func main() {
	noColor := os.Getenv("NO_COLOR") != ""
	logger := logging.SetupLogger(noColor, logging.DefaultLevel)

	defer func() {
		r := recover()
		if r != nil {
			logger.Error(fmt.Sprintf("%v", r))
		}
	}()

	_flags := parseFlags()

	conf := config.Init(logger)
	logger = logging.SetupLoggerAfterConfigLoad(logger, conf, noColor)

	if conf.Receive.Clipboard {
		clipboard.Init()
	}

	// Enable broadcast and monitoring functions
	discovery.New(conf, logger).Start()

	if _flags.ReceiveMode {
		srv := server.New(conf, logger)
		srv.StartHttpServer()
		logger.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	} else {
		sender := send.New(conf, logger)
		err := sender.SendFile(_flags.ToDevice, _flags.FilePath)
		if err != nil {
			logger.Error("Send failed", "err", err)
		}
	}
}
