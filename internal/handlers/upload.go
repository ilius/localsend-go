package handlers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"localsend_cli/internal/config"
	"localsend_cli/internal/discovery/shared"
	"localsend_cli/internal/models"
	"localsend_cli/internal/utils"
)

func SendFileToOtherDevicePrepare(ip, path string) (*models.PrepareReceiveResponse, error) {
	// Prepare metadata for all files
	files := make(map[string]models.FileInfo)
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			sha256Hash, err := utils.CalculateSHA256(filePath)
			if err != nil {
				return fmt.Errorf("error calculating SHA256 hash: %w", err)
			}
			fileMetadata := models.FileInfo{
				ID:       info.Name(), // Use the file name as ID
				FileName: info.Name(),
				Size:     info.Size(),
				FileType: filepath.Ext(filePath),
				SHA256:   sha256Hash,
			}
			files[fileMetadata.ID] = fileMetadata
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking the path: %w", err)
	}

	// Create and fill the PrepareReceiveRequest structure
	request := models.PrepareReceiveRequest{
		Info: models.Info{
			Alias:       shared.Messsage.Alias,
			Version:     shared.Messsage.Version,
			DeviceModel: shared.Messsage.DeviceModel,
			DeviceType:  shared.Messsage.DeviceType,
			Fingerprint: shared.Messsage.Fingerprint,
			Port:        shared.Messsage.Port,
			Protocol:    shared.Messsage.Protocol,
			Download:    shared.Messsage.Download,
		},
		Files: files,
	}

	// Encode the request structure as JSON
	requestJson, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error encoding request to JSON: %w", err)
	}

	// Sending a POST request
	url := fmt.Sprintf("https://%s:53317/api/localsend/v2/prepare-upload", ip)
	client := &http.Client{
		Timeout: 60 * time.Second, // Transmission timeout
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Ignore TLS
			},
		},
	}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return nil, fmt.Errorf("error sending POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 204:
			return nil, fmt.Errorf("finished (No file transfer needed)")
		case 400:
			return nil, fmt.Errorf("invalid body")
		case 403:
			return nil, fmt.Errorf("rejected")
		case 500:
			return nil, fmt.Errorf("unknown error by receiver")
		}
		return nil, fmt.Errorf("failed to send metadata: received status code %d", resp.StatusCode)
	}

	// Decode the response JSON into a PrepareReceiveResponse structure
	var prepareReceiveResponse models.PrepareReceiveResponse
	if err := json.NewDecoder(resp.Body).Decode(&prepareReceiveResponse); err != nil {
		return nil, fmt.Errorf("error decoding response JSON: %w", err)
	}

	return &prepareReceiveResponse, nil
}

func uploadFile(ip, sessionId, fileId, token, filePath string) error {
	// Open the file you want to send
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a request body with file content
	var requestBody bytes.Buffer
	if _, err := io.Copy(&requestBody, file); err != nil {
		return fmt.Errorf("error copying file content: %w", err)
	}

	// Constructing a URL for file upload
	uploadURL := fmt.Sprintf("https://%s:53317/api/localsend/v2/upload?sessionId=%s&fileId=%s&token=%s",
		ip, sessionId, fileId, token)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("POST", uploadURL, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating POST request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending file upload request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case 400:
			return fmt.Errorf("missing parameters")
		case 403:
			return fmt.Errorf("invalid token or IP address")
		case 409:
			return fmt.Errorf("blocked by another session")
		case 500:
			return fmt.Errorf("unknown error by receiver")
		}
		return fmt.Errorf("file upload failed: received status code %d", resp.StatusCode)
	}

	slog.Info("File uploaded successfully")
	return nil
}

func SendFile(ip, path string) error {
	response, err := SendFileToOtherDevicePrepare(ip, path)
	slog.Info("SendFile: got response", "response", response)
	if err != nil {
		return err
	}

	// Traversing directories and sub-files
	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Get fileId and token
			fileId := info.Name() // Use the file name as fileId
			token, ok := response.Files[fileId]
			if !ok {
				return fmt.Errorf("token not found for file: %s", fileId)
			}
			err = uploadFile(ip, response.SessionID, fileId, token, filePath)
			if err != nil {
				return fmt.Errorf("error uploading file: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path: %w", err)
	}

	return nil
}

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

	dst.Close()
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
