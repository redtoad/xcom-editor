package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestIglobFile_Unpack(t *testing.T) {
	buf := MustLoadTestFile(t, "PARTIAL", "IGLOB.DAT")
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
