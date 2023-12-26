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

// SaveDATFile saves a single data to its original location on disk file
// if the content has changed. name is the path inside the game directory.
func SaveDATFile(path string, obj interface{}) error {

	saveData, err := restruct.Pack(binary.LittleEndian, obj)
	if err != nil {
		return fmt.Errorf("could not pack data for file %s: %w", path, err)
	}

	original, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("could not read file %s: %w", path, err)
	}

	if bytes.Equal(saveData, original) {
		return nil
	}
	if err = os.WriteFile(path, saveData, os.ModePerm); err != nil {
		return fmt.Errorf("could not save file %s: %w", path, err)
	}
	return nil
}
