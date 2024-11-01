// Copyright 2013 @atotto. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows && !darwin && !plan9
// +build !windows,!darwin,!plan9

package clipboard

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
)

const (
	xsel               = "xsel"
	xclip              = "xclip"
	wlcopy             = "wl-copy"
	wlpaste            = "wl-paste"
	termuxClipboardGet = "termux-clipboard-get"
	termuxClipboardSet = "termux-clipboard-set"
)

var (
	xselPasteArgs = []string{"--output", "--clipboard"}
	xselCopyArgs  = []string{"--input", "--clipboard"}

	xclipPasteArgs = []string{"-out", "-selection", "clipboard"}
	xclipCopyArgs  = []string{"-in", "-selection", "clipboard"}

	wlpasteArgs = []string{"--no-newline"}
	wlcopyArgs  = []string{}

	termuxPasteArgs = []string{}
	termuxCopyArgs  = []string{}

	missingCommands = errors.New(
		"No clipboard utilities available. Please install xsel, xclip, wl-clipboard or Termux:API add-on for termux-clipboard-get/set.",
	)
)

type commandInfo struct {
	pasteCmdArgs []string
	copyCmdArgs  []string

	unsupported bool
}

var cmd *commandInfo

func initPlatform() {
	cmd = findClipboardUtility()
}

func findClipboardUtility() *commandInfo {
	c := cmd
	if c != nil {
		return c
	}
	c = &commandInfo{}
	cmd = c

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		c.pasteCmdArgs = wlpasteArgs
		c.copyCmdArgs = wlcopyArgs

		if p, err := exec.LookPath(wlcopy); err == nil {
			slog.Info("found wlcopy executable", "path", p)
			if p, err := exec.LookPath(wlpaste); err == nil {
				slog.Info("found wlpaste executable", "path", p)
				return c
			}
		}
	}

	if xclipPath, err := exec.LookPath(xclip); err == nil {
		slog.Info("found xclip executable", "path", xclipPath)
		c.pasteCmdArgs = append([]string{xclipPath}, xclipPasteArgs...)
		c.copyCmdArgs = append([]string{xclipPath}, xclipCopyArgs...)
		return c
	}

	if xselPath, err := exec.LookPath(xsel); err == nil {
		slog.Info("found xsel executable", "path", xselPath)
		c.pasteCmdArgs = append([]string{xselPath}, xselPasteArgs...)
		c.copyCmdArgs = append([]string{xselPath}, xselCopyArgs...)
		return c
	}

	if p, err := exec.LookPath(termuxClipboardSet); err == nil {
		slog.Info("found termuxClipboardSet executable", "path", p)
		c.pasteCmdArgs = append([]string{p}, termuxPasteArgs...)
		c.copyCmdArgs = append([]string{p}, termuxCopyArgs...)
		if _, err := exec.LookPath(termuxClipboardGet); err == nil {
			return c
		}
	}
	c.unsupported = true

	return c
}

func getPasteCommand(c *commandInfo) *exec.Cmd {
	return exec.Command(c.pasteCmdArgs[0], c.pasteCmdArgs[1:]...)
}

func getCopyCommand(c *commandInfo) *exec.Cmd {
	return exec.Command(c.copyCmdArgs[0], c.copyCmdArgs[1:]...)
}

func readAll() (string, error) {
	c := findClipboardUtility()
	if c.unsupported {
		return "", missingCommands
	}

	pasteCmd := getPasteCommand(c)
	out, err := pasteCmd.Output()
	if err != nil {
		return "", err
	}
	result := string(out)
	return result, nil
}

func writeAll(text string) error {
	c := findClipboardUtility()
	if c.unsupported {
		return missingCommands
	}

	copyCmd := getCopyCommand(c)
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}
