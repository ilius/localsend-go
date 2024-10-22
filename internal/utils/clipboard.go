package utils

import (
	"log/slog"
	"os/exec"
	"strings"
)

func WriteToClipBoard(text string) {
	os := OSType()
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
