package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CheckOSType() string {
	return runtime.GOOS
}

func WriteToClipBoard(text string) {
	os := CheckOSType()
	switch os {
	case "linux":
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err != nil {
			slog.Error("Error copying to clipboard on Linux", "err", err)
		} else {
			slog.Info("Text copied to clipboard on Linux!")
		}
	case "windows":
		cmd := exec.Command("cmd", "/c", "echo "+text+" | clip")
		err := cmd.Run()
		if err != nil {
			slog.Error("Error copying to clipboard on Windows", "err", err)
		} else {
			slog.Info("Text copied to clipboard on Windows!")
		}
	case "darwin":
		cmd := exec.Command("pbcopy")
		cmd.Stdin = strings.NewReader(text)
		err := cmd.Run()
		if err != nil {
			slog.Error("Error copying to clipboard on MacOS", "err", err)
		} else {
			slog.Info("Text copied to clipboard on MacOS!")
		}
	default:
		slog.Error("WriteToClipBoard: Unsupported OS", "os", os)
	}
}

func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
