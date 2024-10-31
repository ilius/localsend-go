package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/static"
)

type serverImp struct {
	*http.ServeMux
	conf *config.Config
}

func New(conf *config.Config) *serverImp {
	return &serverImp{
		ServeMux: http.NewServeMux(),
		conf:     conf,
	}
}

func (s *serverImp) StartHttpServer() {
	if s.conf.Functions.HttpFileServer {
		s.addHttpFileServerRoutes()
	}
	if s.conf.Functions.LocalSendServer {
		s.addLocalSendServerRoutes() // Send and receive part
	}
	go func() {
		slog.Info("Server starting on :53317")
		if err := http.ListenAndServe(":53317", s.ServeMux); err != nil {
			panic(fmt.Sprintf("Server failed: %v", err))
		}
	}()
}

// If you enable the http file server, enable the following routes
func (s *serverImp) addHttpFileServerRoutes() {
	mux := s.ServeMux
	mux.HandleFunc("/", s.IndexFileHandler)
	mux.HandleFunc("/uploads/", s.FileServerHandler)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.EmbeddedStaticFiles))),
	)
}

func (s *serverImp) addLocalSendServerRoutes() {
	mux := s.ServeMux
	mux.HandleFunc("/api/localsend/v2/prepare-upload", s.PrepareUploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/upload", s.UploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/info", s.GetInfoHandler)
	mux.HandleFunc("/send", s.UploadHandler)
	mux.HandleFunc("/receive", s.DownloadHandler)
}
