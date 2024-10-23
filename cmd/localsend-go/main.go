package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"localsend_cli/pkg/config"
	"localsend_cli/pkg/discovery"
	"localsend_cli/pkg/handlers"
	"localsend_cli/pkg/send"
	"localsend_cli/pkg/server"
	"localsend_cli/static"
)

const (
	cmd_send    = "send"
	cmd_receive = "receive"
)

func main() {
	mode := flag.String(
		"mode",
		cmd_send,
		"Mode of operation: "+cmd_send+" or "+cmd_receive,
	)
	filePath := flag.String("file", "", "Path to the file to upload")
	toDevice := flag.String("to", "", "Send file to Device ip,Write device receiver ip here")
	flag.Parse()

	switch *mode {
	case cmd_send:
		if *filePath == "" {
			os.Stderr.WriteString("Send mode requires -file FILE_PATH\n")
			flag.Usage()
			os.Exit(1)
		}
		if *toDevice == "" {
			os.Stderr.WriteString("Send mode requires -to DEVICE_IP\n")
			flag.Usage()
			os.Exit(1)
		}
	case cmd_receive:
	default:
		flag.Usage()
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
