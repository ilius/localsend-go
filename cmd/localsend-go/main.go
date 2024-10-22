package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"localsend_cli/internal/config"
	"localsend_cli/internal/discovery"
	"localsend_cli/internal/handlers"
	"localsend_cli/internal/pkg/server"
	"localsend_cli/internal/send"
	"localsend_cli/static"
)

func main() {
	mode := flag.String("mode", "send", "Mode of operation: send or receive")
	filePath := flag.String("file", "", "Path to the file to upload")
	toDevice := flag.String("to", "", "Send file to Device ip,Write device receiver ip here")
	flag.Parse()

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
		httpServer.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(static.EmbeddedStaticFiles))))
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
	case "send":
		if *filePath == "" {
			slog.Info("Send mode requires a file path")
			flag.Usage()
			os.Exit(1)
		}
		if *toDevice == "" {
			slog.Info("Send mode requires a toDevice")
			flag.Usage()
			os.Exit(1)
		}
		err := send.SendFile(*toDevice, *filePath)
		if err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	case "receive":
		slog.Info("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	default:
		flag.Usage()
		os.Exit(1)
	}
}
