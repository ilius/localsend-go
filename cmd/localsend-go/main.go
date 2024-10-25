package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/discovery"
	"github.com/ilius/localsend-go/pkg/handlers"
	"github.com/ilius/localsend-go/pkg/send"
	"github.com/ilius/localsend-go/pkg/server"
	"github.com/ilius/localsend-go/pkg/slogcolor"
	"github.com/ilius/localsend-go/pkg/static"
)

const (
	cmd_send    = "send"
	cmd_receive = "receive"
)

func main() {
	{
		handler := slogcolor.NewHandler(os.Stdout, &slogcolor.Options{
			Level:         slog.LevelInfo,
			TimeFormat:    time.DateTime,
			SrcFileMode:   slogcolor.ShortFile,
			SrcFileLength: 0,
			// MsgPrefix:     color.HiWhiteString("| "),
			MsgLength: 0,
			MsgColor:  color.New(),
			NoColor:   false,
		})
		logger := slog.New(handler)
		slog.SetDefault(logger)
		defer func() {
			r := recover()
			logger.Error("Panic", "r", r)
		}()
	}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	mode := flagSet.String(
		"mode",
		cmd_send,
		"Mode of operation: "+cmd_send+" or "+cmd_receive,
	)
	filePath := flagSet.String("file", "", "Path to the file to upload")
	toDevice := flagSet.String("to", "", "Send file to Device ip,Write device receiver ip here")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	switch *mode {
	case cmd_send:
		if *filePath == "" {
			os.Stderr.WriteString("Send mode requires -file FILE_PATH\n")
			flagSet.Usage()
			os.Exit(1)
		}
		if *toDevice == "" {
			os.Stderr.WriteString("Send mode requires -to DEVICE_IP\n")
			flagSet.Usage()
			os.Exit(1)
		}
	case cmd_receive:
	default:
		flagSet.Usage()
		os.Exit(1)
	}

	config.Init()

	// Enable broadcast and monitoring functions
	go discovery.ListenForBroadcasts()
	go discovery.StartBroadcast()
	go discovery.StartHTTPBroadcast() // Start HTTP Broadcast

	// Start HTTP Server
	httpServer := server.New()
	if config.ConfigData.Functions.HttpFileServer {
		// If you enable the http file server, enable the following routes
		httpServer.HandleFunc("/", handlers.IndexFileHandler)
		httpServer.HandleFunc("/uploads/", handlers.FileServerHandler)
		httpServer.Handle(
			"/static/",
			http.StripPrefix("/static/", http.FileServer(http.FS(static.EmbeddedStaticFiles))),
		)
	}

	// Send and receive part
	if config.ConfigData.Functions.LocalSendServer {
		httpServer.HandleFunc("/api/localsend/v2/prepare-upload", handlers.PrepareUploadAPIHandler)
		httpServer.HandleFunc("/api/localsend/v2/upload", handlers.UploadAPIHandler)
		httpServer.HandleFunc("/api/localsend/v2/info", handlers.GetInfoHandler)
		httpServer.HandleFunc("/send", handlers.UploadHandler)
		httpServer.HandleFunc("/receive", handlers.DownloadHandler)
	}

	go func() {
		slog.Info("Server started at :53317")
		if err := http.ListenAndServe(":53317", httpServer); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	switch *mode {
	case cmd_send:
		err := send.SendFile(*toDevice, *filePath)
		if err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	case cmd_receive:
		slog.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	}
}
