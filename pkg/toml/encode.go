package toml

import (
	"reflect"
	"strings"
)

var dblQuotedReplacer = strings.NewReplacer(
	"\"", "\\\"",
	"\\", "\\\\",
	"\x00", `\u0000`,
	"\x01", `\u0001`,
	"\x02", `\u0002`,
	"\x03", `\u0003`,
	"\x04", `\u0004`,
	"\x05", `\u0005`,
	"\x06", `\u0006`,
	"\x07", `\u0007`,
	"\b", `\b`,
	"\t", `\t`,
	"\n", `\n`,
	"\x0b", `\u000b`,
	"\f", `\f`,
	"\r", `\r`,
	"\x0e", `\u000e`,
	"\x0f", `\u000f`,
	"\x10", `\u0010`,
	"\x11", `\u0011`,
	"\x12", `\u0012`,
	"\x13", `\u0013`,
	"\x14", `\u0014`,
	"\x15", `\u0015`,
	"\x16", `\u0016`,
	"\x17", `\u0017`,
	"\x18", `\u0018`,
	"\x19", `\u0019`,
	"\x1a", `\u001a`,
	"\x1b", `\u001b`,
	"\x1c", `\u001c`,
	"\x1d", `\u001d`,
	"\x1e", `\u001e`,
	"\x1f", `\u001f`,
	"\x7f", `\u007f`,
)

// Marshaler is the interface implemented by types that can marshal themselves
// into valid TOML.
type Marshaler interface {
	MarshalTOML() ([]byte, error)
}

type tagOptions struct {
	skip      bool // "-"
	name      string
	omitempty bool
	omitzero  bool
}

func getOptions(tag reflect.StructTag) tagOptions {
	t := tag.Get("toml")
	if t == "-" {
		return tagOptions{skip: true}
	}
	var opts tagOptions
	parts := strings.Split(t, ",")
	opts.name = parts[0]
	for _, s := range parts[1:] {
		switch s {
		case "omitempty":
			opts.omitempty = true
		case "omitzero":
			opts.omitzero = true
		}
	}
	return opts
}
