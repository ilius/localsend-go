package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// uploadHandler handles file upload requests
func (s *serverImp) uploadHandler(w http.ResponseWriter, r *http.Request) {
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
		s.log.Error("Failed to close file", "err", err)
		http.Error(w, fmt.Sprintf("Could not save file: %v", err), http.StatusInternalServerError)
		return
	}
	s.changeFileOwnerGroup(filePath)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func (s *serverImp) changeFileOwnerGroup(filePath string) {
	if s.conf.Receive.SaveUserID > 0 || s.conf.Receive.SaveGroupID > 0 {
		s.log.Debug("Changing file ownership and group")
		err := os.Chown(filePath, s.conf.Receive.SaveUserID, s.conf.Receive.SaveGroupID)
		if err != nil {
			s.log.Error("Failed to change ownership of file", "err", err)
		}
	}
}
