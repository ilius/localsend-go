package startup

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/discovery"
	"github.com/ilius/localsend-go/pkg/go-clipboard"
	"github.com/ilius/localsend-go/pkg/handlers"
	"github.com/ilius/localsend-go/pkg/server"
	"github.com/ilius/localsend-go/pkg/static"
)

func StartupServices(conf *config.Config, receiveMode bool) {
	if conf.Receive.Clipboard {
		clipboard.Init()
	}

	startDiscovery(conf) // Enable broadcast and monitoring functions

	mux := server.New()

	if receiveMode {
		if conf.Functions.HttpFileServer {
			addHttpFileServerRoutes(mux) // Start HTTP Server
		}
		if conf.Functions.LocalSendServer {
			addLocalSendServerRoutes(mux) // Send and receive part
		}
	}

	go func() {
		slog.Info("Server starting on :53317")
		if err := http.ListenAndServe(":53317", mux); err != nil {
			panic(fmt.Sprintf("Server failed: %v", err))
		}
	}()
}

// Enable broadcast and monitoring functions
func startDiscovery(conf *config.Config) {
	go discovery.ListenForBroadcasts()
	go discovery.StartBroadcast(conf)
	go discovery.StartHTTPBroadcast(conf)
}

// If you enable the http file server, enable the following routes
func addHttpFileServerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handlers.IndexFileHandler)
	mux.HandleFunc("/uploads/", handlers.FileServerHandler)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.EmbeddedStaticFiles))),
	)
}

func addLocalSendServerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/localsend/v2/prepare-upload", handlers.PrepareUploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/upload", handlers.UploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/info", handlers.GetInfoHandler)
	mux.HandleFunc("/send", handlers.UploadHandler)
	mux.HandleFunc("/receive", handlers.DownloadHandler)
}
