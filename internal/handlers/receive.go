package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"localsend_cli/internal/models"
	"localsend_cli/internal/utils"
)

var (
	sessionIDCounter = 0
	sessionMutex     sync.Mutex
	fileNames        = make(map[string]string) // To save the file name
)

func PrepareReceive(w http.ResponseWriter, r *http.Request) {
	var req models.PrepareReceiveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	slog.Info("Received request:", "req", req)

	sessionMutex.Lock()
	sessionIDCounter++
	sessionID := fmt.Sprintf("session-%d", sessionIDCounter)
	sessionMutex.Unlock()

	files := make(map[string]string)
	for fileID, fileInfo := range req.Files {
		token := fmt.Sprintf("token-%s", fileID)
		files[fileID] = token

		// Save file name
		fileNames[fileID] = fileInfo.FileName

		if strings.HasSuffix(fileInfo.FileName, ".txt") {
			slog.Info("TXT file content preview", "preview", string(fileInfo.Preview))
			utils.WriteToClipBoard(fileInfo.Preview)
		}
	}

	resp := models.PrepareReceiveResponse{
		SessionID: sessionID,
		Files:     files,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
