package files

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-restruct/restruct"
)

var DefaultByteOrder = binary.LittleEndian

// LoadDATFile loads binary file found at path and populates obj instance.
func LoadDATFile(path string, obj interface{}) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", path, err)
	}
	if err = restruct.Unpack(buf, DefaultByteOrder, obj); err != nil {
		return fmt.Errorf("could not unpack data: %w", err)
	}
	return nil
}
