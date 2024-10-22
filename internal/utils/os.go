package utils

import (
	"runtime"
)

func OSType() string {
	return runtime.GOOS
}
