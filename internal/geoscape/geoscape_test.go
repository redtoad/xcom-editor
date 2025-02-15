package geoscape_test

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"strings"
	"testing"
	"unicode"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

// MustLoadFromBase64 loads a file content from a hex-encoded string.
func MustLoadFromHex(str string) []byte {
	txt := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, str)
	buffer, err := hex.DecodeString(txt)
	if err != nil {
		panic(err)
	}
	return buffer
}

// MustLoadFromBase64 loads a file content from a base64-encoded string.
func MustLoadFromBase64(base64String string) []byte {
	buffer, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		panic(err)
	}
	return buffer
}

func TestFileUnpacking(t *testing.T) {

	tt := []struct {
		name string
		obj  restruct.Unpacker
		data string
	}{
		{
			name: "BASE.DAT",
			obj:  &geoscape.BaseFile{},
			data: testFile_BASE_DAT,
		},
		{
			name: "CRAFT.DAT",
			obj:  &geoscape.CraftFile{},
			data: testFile_CRAFT_DAT,
		},
		{
			name: "INTER.DAT",
			obj:  &geoscape.InterFile{},
			data: testFile_INTER_DAT,
		},
		{
			name: "LOC.DAT",
			obj:  &geoscape.LocFile{},
			data: testFile_LOC_DAT,
		},
		{
			name: "SOLDIER.DAT",
			obj:  &geoscape.SoldierFile{},
			data: testFile_SOLDIER_DAT,
		},
		{
			name: "TRANSFER.DAT",
			obj:  &geoscape.TransferFile{},
			data: testFile_TRANSFER_DAT,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			buffer := MustLoadFromBase64(tc.data)
			remaining, err := tc.obj.Unpack(buffer, binary.LittleEndian)
			if err != nil {
				t.Fatalf("could not unpack %T: %v", tc.obj, err)
			}
			if len(remaining) != 0 {
				t.Errorf("expected no remaining bytes, got %d", len(remaining))
			}
		})
	}
}
