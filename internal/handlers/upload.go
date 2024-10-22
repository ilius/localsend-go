package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"localsend_cli/internal/config"
)

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
