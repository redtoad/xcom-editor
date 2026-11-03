package internal

import (
	"bytes"
)

// NullString is a null terminated string stored as raw bytes.
// It preserves the original bytes (including any garbage after the null terminator),
// which is required for byte-identical round trips. Restruct handles serialization
// natively via the struct tag (e.g. [25]byte).
type NullString []byte

// String returns the string value, truncated at the first null byte.
func (s NullString) String() string {
	if nul := bytes.IndexByte(s, 0x0); nul >= 0 {
		return string(s[:nul])
	}
	return string(s)
}

// NewNullString returns a NullString of the given fixed size, padded with null bytes.
// Values that do not fit are truncated so that the last byte is always a null byte.
func NewNullString(value string, size int) NullString {
	buf := make([]byte, size)
	copy(buf[:size-1], value)
	return NullString(buf)
}
