package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ilius/localsend-go/pkg/config"
	"github.com/ilius/localsend-go/pkg/static"
)

type serverImp struct {
	mux  *http.ServeMux
	conf *config.Config
	log  *slog.Logger
}

func New(conf *config.Config, logger *slog.Logger) *serverImp {
	return &serverImp{
		mux:  http.NewServeMux(),
		conf: conf,
		log:  logger,
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
		s.log.Info("Server starting on :53317")
		if err := http.ListenAndServe(":53317", s.mux); err != nil {
			panic(fmt.Sprintf("Server failed: %v", err))
		}
	}()
}

// If you enable the http file server, enable the following routes
func (s *serverImp) addHttpFileServerRoutes() {
	mux := s.mux
	mux.HandleFunc("/", s.indexFileHandler)
	mux.HandleFunc("/uploads/", s.fileServerHandler)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.FS(static.EmbeddedStaticFiles))),
	)
}

func (s *serverImp) addLocalSendServerRoutes() {
	mux := s.mux
	mux.HandleFunc("/api/localsend/v2/prepare-upload", s.prepareUploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/upload", s.uploadAPIHandler)
	mux.HandleFunc("/api/localsend/v2/info", s.getInfoHandler)
	mux.HandleFunc("/send", s.uploadHandler)
	mux.HandleFunc("/receive", s.downloadHandler)
}
