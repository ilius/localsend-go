package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilius/localsend-go/pkg/discovery/shared"
)

func (s *serverImp) GetInfoHandler(w http.ResponseWriter, r *http.Request) {
	msg := shared.GetMesssage(s.conf)
	res, err := json.Marshal(msg)
	if err != nil {
		slog.Error("json convert failed:", "err", err)
		http.Error(w, "json convert failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		slog.Error("Error writing file:", "err", err)
		return
	}
}
