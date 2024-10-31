package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadHandler handles file download requests
func (s *serverImp) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if fileName == "" {
		http.Error(w, "File parameter is required", http.StatusBadRequest)
		return
	}

	// Assuming the files are stored in the configured directory
	filePath := filepath.Join(s.conf.Send.Directory, fileName)
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
