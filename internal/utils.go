package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-restruct/restruct"
)

var DefaultByteOrder = binary.LittleEndian

// LoadDATFile loads binary file found at path and populates obj instance.
func LoadDATFile(path string, obj any) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not open file %s: %w", path, err)
	}
	if err = Unmarshall(buf, obj); err != nil {
		return fmt.Errorf("could not unpack data: %w", err)
	}
	return nil
}

// Unmarshall decodes binary data from a buffer into a struct.
func Unmarshall(buffer []byte, obj any) error {
	if err := restruct.Unpack(buffer, DefaultByteOrder, obj); err != nil {
		return fmt.Errorf("could not unpack data: %w", err)
	}
	return nil
}

// SaveDATFile saves a single data to its original location on disk file
// if the content has changed. name is the path inside the game directory.
func SaveDATFile(path string, obj any) error {

	saveData, err := Marshall(obj)
	if err != nil {
		return fmt.Errorf("could not marshal data for file %s: %w", path, err)
	}

	original, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}

	if bytes.Equal(saveData, original) {
		return nil
	}
	if err = os.WriteFile(path, saveData, 0644); err != nil {
		return fmt.Errorf("could not save file %s: %w", path, err)
	}
	return nil
}

// Marshall encodes a struct into binary data.
func Marshall(obj interface{}) ([]byte, error) {
	buffer, err := restruct.Pack(binary.LittleEndian, obj)
	if err != nil {
		return nil, fmt.Errorf("could not pack data: %w", err)
	}
	return buffer, nil
}
