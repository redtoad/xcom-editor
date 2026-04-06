package internal

import (
	"strings"
)

// NullString is a null byte terminated string. When read from binary data via
// restruct, the full fixed-size byte field (including null bytes) is stored.
// Use String() to get the value truncated at the first null byte.
type NullString string

// String returns the string value, truncated at the first null byte.
func (s NullString) String() string {
	if nul := strings.IndexByte(string(s), 0x0); nul >= 0 {
		return string(s[:nul])
	}
	return string(s)
}

// NewNullString returns a NullString of the given fixed size, padded with null bytes.
func NewNullString(value string, size int) NullString {
	buf := make([]byte, size)
	copy(buf, value)
	return NullString(buf)
}
