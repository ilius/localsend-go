package server

import (
	"encoding/json"
	"net/http"

	"codeberg.org/ilius/localsend-go/pkg/discovery/shared"
)

func (s *serverImp) getInfoHandler(w http.ResponseWriter, _ *http.Request) {
	msg := shared.GetMesssage(s.conf)
	res, err := json.Marshal(msg)
	if err != nil {
		s.log.Error("json convert failed:", "err", err)
		http.Error(w, "json convert failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		s.log.Error("Error writing file:", "err", err)
		return
	}
}
