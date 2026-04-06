package internal

import (
	"bytes"
)

// NullString is a null terminated string stored as raw bytes.
// It preserves the original bytes (including any garbage after the null terminator).
// Restruct handles serialization natively via the struct tag (e.g. [25]byte).
type NullString []byte

// String returns the string value without the null terminator.
func (s NullString) String() string {
	nul := bytes.IndexByte(s, 0x0)
	if nul < 0 {
		return string(s)
	}
	return string(s[:nul])
}

func NewNullString(value string, size int) NullString {
	buf := make([]byte, size)
	copy(buf, []byte(value))
	return NullString(buf)
}
