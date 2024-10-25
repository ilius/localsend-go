package toml

import (
	"strings"
)

// MetaData allows access to meta information about TOML data that's not
// accessible otherwise.
//
// It allows checking if a key is defined in the TOML data, whether any keys
// were undecoded, and the TOML type of a key.
type MetaData struct {
	context Key // Used only during decoding.

	keyInfo map[string]keyInfo
	mapping map[string]any
	keys    []Key
	decoded map[string]struct{}
	data    []byte // Input file; for errors.
}

// Key represents any TOML key, including key groups. Use [MetaData.Keys] to get
// values of this type.
type Key []string

func (k Key) String() string {
	// This is called quite often, so it's a bit funky to make it faster.
	var b strings.Builder
	b.Grow(len(k) * 25)
outer:
	for i, kk := range k {
		if i > 0 {
			b.WriteByte('.')
		}
		if kk == "" {
			b.WriteString(`""`)
		} else {
			for _, r := range kk {
				// "Inline" isBareKeyChar
				if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
					b.WriteByte('"')
					b.WriteString(dblQuotedReplacer.Replace(kk))
					b.WriteByte('"')
					continue outer
				}
			}
			b.WriteString(kk)
		}
	}
	return b.String()
}

// Like append(), but only increase the cap by 1.
func (k Key) add(piece string) Key {
	if cap(k) > len(k) {
		return append(k, piece)
	}
	newKey := make(Key, len(k)+1)
	copy(newKey, k)
	newKey[len(k)] = piece
	return newKey
}

func (k Key) parent() Key  { return k[:len(k)-1] } // all except the last piece.
func (k Key) last() string { return k[len(k)-1] }  // last piece of this key.
