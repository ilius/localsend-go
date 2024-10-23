package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"localsend_cli/pkg/config"
)

// UploadHandler handles file upload requests
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parsing multipart/form-data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, fmt.Sprintf("Could not parse multipart form: %v", err), http.StatusBadRequest)
		return
	}

	// Get File
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get uploaded file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create the upload directory if it does not exist
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, fmt.Sprintf("Could not create upload directory: %v", err), http.StatusInternalServerError)
		return
	}

	// Creating a target file
	filePath := filepath.Join(uploadDir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not create file: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the uploaded file content to the target file
	if _, err := io.Copy(dst, file); err != nil {
		dst.Close()
		http.Error(w, fmt.Sprintf("Could not save file: %v", err), http.StatusInternalServerError)
		return
	}

	err = dst.Close()
	if err != nil {
		slog.Error("Failed to close file", "err", err)
		http.Error(w, fmt.Sprintf("Could not save file: %v", err), http.StatusInternalServerError)
		return
	}
	changeFileOwnerGroup(filePath)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func changeFileOwnerGroup(filePath string) {
	if config.ConfigData.Receive.SaveUserID > 0 || config.ConfigData.Receive.SaveGroupID > 0 {
		slog.Debug("Changing file ownership and group")
		err := os.Chown(filePath, config.ConfigData.Receive.SaveUserID, config.ConfigData.Receive.SaveGroupID)
		if err != nil {
			slog.Error("Failed to change ownership of file", "err", err)
		}
	}
}
