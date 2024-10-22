package main

import (
	"flag"
	"fmt"
	"localsend_cli/internal/config"
	"localsend_cli/internal/discovery"
	"localsend_cli/internal/handlers"
	"localsend_cli/internal/pkg/server"
	"localsend_cli/static"
	"log"
	"net/http"
	"os"
)

func main() {
	mode := flag.String("mode", "send", "Mode of operation: send or receive")
	filePath := flag.String("file", "", "Path to the file to upload")
	toDevice := flag.String("to", "", "Send file to Device ip,Write device receiver ip here")
	flag.Parse()

	// Enable broadcast and monitoring functions
	go discovery.ListenForBroadcasts()
	go discovery.StartBroadcast()
	go discovery.StartHTTPBroadcast() // 启动HTTP广播

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

		httpServer.HandleFunc("/api/localsend/v2/prepare-upload", handlers.PrepareReceive)
		httpServer.HandleFunc("/api/localsend/v2/upload", handlers.ReceiveHandler)
		httpServer.HandleFunc("/api/localsend/v2/info", handlers.GetInfoHandler)
		httpServer.HandleFunc("/send", handlers.NormalSendHandler)       // Upload Handler
		httpServer.HandleFunc("/receive", handlers.NormalReceiveHandler) // Download Handler

	}
	go func() {
		log.Println("Server started at :53317")
		if err := http.ListenAndServe(":53317", httpServer); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	switch *mode {
	case "send":
		if *filePath == "" {
			fmt.Println("Send mode requires a file path")
			flag.Usage()
			os.Exit(1)
		}
		if *toDevice == "" {
			fmt.Println("Send mode requires a toDevice")
			flag.Usage()
			os.Exit(1)
		}
		err := handlers.SendFile(*toDevice, *filePath)
		if err != nil {
			log.Fatalf("Send failed: %v", err)
		}
		// if err := sendFile(*filePath); err != nil {
		// 	log.Fatalf("Send failed: %v", err)
		// }
	case "receive":
		fmt.Println("Waiting to receive files...")
		select {} // Blocking program waiting to receive file
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func sendFile(filePath string) error {
	fmt.Println("Sending file:", filePath)
	// Logic for uploading files
	return nil
}
