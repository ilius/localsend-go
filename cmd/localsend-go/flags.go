package main

import (
	"flag"
	"os"
)

const (
	cmd_send    = "send"
	cmd_receive = "receive"
)

type Flags struct {
	ReceiveMode bool
	FilePath    string
	ToDevice    string
}

func parseFlags() *Flags {
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	mode := flagSet.String(
		"mode",
		cmd_send,
		"Mode of operation: "+cmd_send+" or "+cmd_receive,
	)
	filePath := flagSet.String("file", "", "Path to the file to upload")
	toDevice := flagSet.String("to", "", "Send file to Device ip,Write device receiver ip here")
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	switch *mode {
	case cmd_send:
		if *filePath == "" {
			os.Stderr.WriteString("Send mode requires -file FILE_PATH\n")
			flagSet.Usage()
			os.Exit(1)
		}
		if *toDevice == "" {
			os.Stderr.WriteString("Send mode requires -to DEVICE_IP\n")
			flagSet.Usage()
			os.Exit(1)
		}
		return &Flags{
			FilePath: *filePath,
			ToDevice: *toDevice,
		}
	case cmd_receive:
		return &Flags{
			ReceiveMode: true,
		}
	default:
		flagSet.Usage()
		os.Exit(1)
	}

	return nil
}
