package geoscape_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestIglobFile_Unpack(t *testing.T) {
	buf := MustLoadFromBase64(testFile_IGLOB_DAT)
	var f geoscape.IglobFile
	err := restruct.Unpack(buf, binary.LittleEndian, &f)
	assert.NoError(t, err)

	assert.Equal(t, int32(7), f.Month)
	assert.Equal(t, int32(3), f.Weekday)
	assert.Equal(t, int32(18), f.Day)
	assert.Equal(t, int32(18), f.Hour)
	assert.Equal(t, int32(48), f.Minute)
	assert.Equal(t, int32(45), f.Second)
	assert.Equal(t, int32(0), f.HiddenCraft)
	assert.Equal(t, int32(0), f.AirborneXCraft)
	assert.Equal(t, int32(128), f.CurrentMIDI)
	// 60-byte DOS version: no difficulty field
	assert.Nil(t, f.Difficulty)
}

func TestIglobFile_RoundTrip(t *testing.T) {
	buf := MustLoadFromBase64(testFile_IGLOB_DAT)

	var f geoscape.IglobFile
	if err := restruct.Unpack(buf, binary.LittleEndian, &f); err != nil {
		t.Fatalf("could not unpack IglobFile: %v", err)
	}

	packed, err := restruct.Pack(binary.LittleEndian, &f)
	if err != nil {
		t.Fatalf("could not pack IglobFile: %v", err)
	}

	if !bytes.Equal(buf, packed) {
		t.Errorf("re-packed data does not match original")
		fmt.Println("Original:", hex.Dump(buf))
		fmt.Println("Re-packed:", hex.Dump(packed))
	}
}

// base64 -i GAME/GAME_1/IGLOB.DAT | pbcopy
const testFile_IGLOB_DAT = `BwAAAAMAAAASAAAAEgAAADAAAAAtAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACAAAAA`
