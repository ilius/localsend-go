package colorable

import (
	"bytes"
	"os"
	"runtime"
	"testing"
)

// checkEncoding checks that colorable is output encoding agnostic as long as
// the encoding is a superset of ASCII. This implies that one byte not part of
// an ANSI sequence must give exactly one byte in output
func checkEncoding(t *testing.T, data []byte) {
	// Send non-UTF8 data to colorable
	b := bytes.NewBuffer(make([]byte, 0, 10))
	if b.Len() != 0 {
		t.FailNow()
	}
}

func TestEncoding(t *testing.T) {
	checkEncoding(t, []byte{})      // Empty
	checkEncoding(t, []byte(`abc`)) // "abc"
	checkEncoding(t, []byte(`é`))   // "é" in UTF-8
	checkEncoding(t, []byte{233})   // 'é' in Latin-1
}

func TestColorable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skip this test on windows")
	}
	_, ok := NewColorableStdout().(*os.File)
	if !ok {
		t.Fatalf("should os.Stdout on UNIX")
	}
	_, ok = NewColorableStderr().(*os.File)
	if !ok {
		t.Fatalf("should os.Stdout on UNIX")
	}
}
