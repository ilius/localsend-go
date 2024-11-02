package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ilius/localsend-go/pkg/go-clipboard"
	"github.com/ilius/localsend-go/pkg/models"
)

var (
	sessionIDCounter = &atomic.Int64{}
	fileNames        = make(map[string]string) // To save the file name
	fileNamesRWMutex sync.RWMutex
	uploadCount      = &atomic.Int64{}
)

func (s *serverImp) prepareUploadAPIHandler(w http.ResponseWriter, r *http.Request) {
	var req models.PrepareReceiveRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	s.log.Info("Received request:", "req", req)

	sessionID := fmt.Sprintf("session-%d", sessionIDCounter.Add(1))

	files := make(map[string]string)
	for fileID, fileInfo := range req.Files {
		token := fmt.Sprintf("token-%s", fileID)
		files[fileID] = token

		// Save file name
		fileNamesRWMutex.Lock()
		fileNames[fileID] = fileInfo.FileName
		fileNamesRWMutex.Unlock()

		if strings.HasSuffix(fileInfo.FileName, ".txt") {
			s.log.Info("TXT file content preview", "preview", fileInfo.Preview)
			if s.conf.Receive.Clipboard {
				err := clipboard.WriteAll(fileInfo.Preview)
				if err != nil {
					s.log.Error("Error copying to clipboard", "err", err)
				}
			}
		}
	}

	resp := models.PrepareReceiveResponse{
		SessionID: sessionID,
		Files:     files,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *serverImp) uploadAPIHandler(w http.ResponseWriter, r *http.Request) {
	if s.receiveIpBlocked(r) {
		s.log.Warn("uploadAPIHandler blocked IP that is not allowd", "ip", r.RemoteAddr)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := r.URL.Query().Get("sessionId")
	fileID := r.URL.Query().Get("fileId")
	token := r.URL.Query().Get("token")

	// Verify request parameters
	if sessionID == "" || fileID == "" || token == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	// Use fileID to get the file name
	fileNamesRWMutex.RLock()
	fileName, ok := fileNames[fileID]
	fileNamesRWMutex.RUnlock()
	if !ok {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	// Generate file paths, preserving file extensions
	filePath := filepath.Join(s.conf.Receive.Directory, fileName)
	// Create the folder if it does not exist
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		s.log.Error("Error creating directory", "err", err)
		return
	}

	if s.conf.Receive.ExitAfterFileCount > 0 {
		defer s.checkExitAfterFileCount()
	}

	// Create a file
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		s.log.Error("Error creating file", "err", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 2*1024*1024) // 2MB buffer
	size := 0
	maxSize := s.conf.Receive.MaxFileSize
	for {
		n, err := r.Body.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			s.log.Error("Error reading file", "err", err)
			return
		}
		if n == 0 {
			break
		}
		size += n
		if size > maxSize {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			s.log.Error("Max file size reached", "size", size, "maxSize", maxSize)
		}

		_, err = file.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Failed to write file", http.StatusInternalServerError)
			s.log.Error("Error writing file", "err", err)
			return
		}
	}
	s.changeFileOwnerGroup(filePath)

	s.log.Info("Saved file", "filePath", filePath)
	w.WriteHeader(http.StatusOK)
}

func (s *serverImp) checkExitAfterFileCount() {
	count := int(uploadCount.Add(1))
	if count < s.conf.Receive.ExitAfterFileCount {
		return
	}
	s.log.Info("Exiting due to max recieved file count reached")
	go func() {
		time.Sleep(100 * time.Millisecond)
		os.Exit(0)
	}()
}
