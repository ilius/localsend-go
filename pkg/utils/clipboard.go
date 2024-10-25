package utils

import (
	"log/slog"

	"github.com/atotto/clipboard"
)

func WriteToClipBoard(text string) {
	err := clipboard.WriteAll(text)
	if err != nil {
		slog.Error("Error copying to clipboard", "err", err)
	}
}
