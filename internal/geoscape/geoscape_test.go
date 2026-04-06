package geoscape_test

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"
	"unicode"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

// MustLoadFromHex loads a file content from a hex-encoded string.
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

func TestFileUnpackingAndPackingIsBinaryIdentical(t *testing.T) {

	tt := []struct {
		name string
		obj  any
		data string
	}{
		{
			name: "CRAFT.DAT",
			obj:  &geoscape.CraftFile{},
			data: testFile_CRAFT_DAT,
		},
		{
			name: "BASE.DAT",
			obj:  &geoscape.BaseFile{},
			data: testFile_BASE_DAT,
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
			name: "LIGLOB.DAT",
			obj:  &geoscape.LiglobFile{},
			data: testFile_LIGLOB_DAT,
		},
		{
			name: "SAVEINFO.DAT",
			obj:  &geoscape.SaveinfoFile{},
			data: testFile_SAVEINFO_DAT,
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

			if err := restruct.Unpack(buffer, binary.LittleEndian, tc.obj); err != nil {
				t.Fatalf("could not unpack %T: %v", tc.obj, err)
			}

			packed, err := restruct.Pack(binary.LittleEndian, tc.obj)
			if err != nil {
				t.Fatalf("could not pack %T: %v", tc.obj, err)
			}
			if !bytes.Equal(buffer, packed) {
				t.Errorf("re-packed data does not match original data")
				fmt.Println("Original:", hex.Dump(buffer))
				fmt.Println("Re-packed:", hex.Dump(packed))

				// dump contents of buffer and packed to files for manual inspection
				if err := os.WriteFile(fmt.Sprintf("%s_original.bin", tc.name), buffer, 0644); err != nil {
					t.Fatalf("could not write original data to file: %v", err)
				}
				if err := os.WriteFile(fmt.Sprintf("%s_repacked.bin", tc.name), packed, 0644); err != nil {
					t.Fatalf("could not write re-packed data to file: %v", err)
				}

			}
		})
	}
}
