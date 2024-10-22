package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"localsend_cli/internal/config"
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

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("sessionId")
	fileID := r.URL.Query().Get("fileId")
	token := r.URL.Query().Get("token")

	// Verify request parameters
	if sessionID == "" || fileID == "" || token == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Use fileID to get the file name
	fileName, ok := fileNames[fileID]
	if !ok {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Generate file paths, preserving file extensions
	filePath := filepath.Join("uploads", fileName)
	// Create the folder if it does not exist
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		slog.Error("Error creating directory", "err", err)
		return
	}
	// Create a file
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		slog.Error("Error creating file", "err", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 2*1024*1024) // 2MB buffer
	for {
		n, err := r.Body.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			slog.Error("Error reading file", "err", err)
			return
		}
		if n == 0 {
			break
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			slog.Error("Error writing file", "err", err)
			return

		}
	}
	if config.ConfigData.Receive.SaveUserID > 0 || config.ConfigData.Receive.SaveGroupID > 0 {
		slog.Debug("Changing file ownership and group")
		err := os.Chown(filePath, config.ConfigData.Receive.SaveUserID, config.ConfigData.Receive.SaveGroupID)
		if err != nil {
			slog.Error("Failed to change ownership of file", "err", err)
		}
	}

	slog.Info("Saved file", "filePath", filePath)
	w.WriteHeader(http.StatusOK)
}

// DownloadHandler handles file download requests
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File parameter is required", http.StatusBadRequest)
		return
	}

	// Assuming the files are stored in the "uploads" directory
	filePath := filepath.Join("uploads", fileName)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not open file: %v", err), http.StatusNotFound)
		return
	}
	defer file.Close()

	// Setting the response header
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	// Write the file contents to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not write file to response: %v", err), http.StatusInternalServerError)
		return
	}
}
