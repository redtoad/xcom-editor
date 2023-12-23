package internal

import (
	"encoding/binary"
	"strings"
)

// NullString is a null byte terminated string.
type NullString string

// Unpack implements the restruct.Unpacker interface.
func (s *NullString) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	str := string(buf)
	nul := strings.IndexByte(str, 0x0)
	*s = NullString(str[0:nul])
	return []byte{}, nil
}
