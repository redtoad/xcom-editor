package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BASE_DAT_RoundTrip(t *testing.T) {
	// Create a buffer of 2336 bytes (8 bases * 292 bytes each)
	data := make([]byte, 8*292)

	// Set up first base with a name and active flag
	// Name: "XCOM BASE" (16 bytes, null-terminated)
	copy(data[0:16], []byte("XCOM BASE\x00\x00\x00\x00\x00\x00\x00"))

	// Short range radar: 10
	data[16] = 0x0A
	data[17] = 0x00

	// Grid: access lift at position 0
	data[22] = byte(geoscape.AccessLift)
	// Rest of grid defaults to 0 (also AccessLift) - set to Empty (0xFF)
	for i := 23; i < 22+36; i++ {
		data[i] = 0xFF
	}
	data[22] = byte(geoscape.AccessLift)

	// Engineers at offset 94 (16+2+2+2+36+36)
	data[94] = 10
	// Scientists at offset 95
	data[95] = 5

	// Active flag at offset 288 (16+2+2+2+36+36+2+192)
	// Active=true means the int32 is 0 (invertedbool)
	data[288] = 0x00
	data[289] = 0x00
	data[290] = 0x00
	data[291] = 0x00

	// Set remaining bases as inactive (Active=false means int32=1 with invertedbool)
	for b := 1; b < 8; b++ {
		offset := b * 292
		data[offset+288] = 0x01
	}

	var obj geoscape.BaseFile
	err := restruct.Unpack(data, binary.LittleEndian, &obj)
	require.NoError(t, err)
	assert.Equal(t, 8, len(obj.Bases))
	assert.Equal(t, "XCOM BASE", obj.Bases[0].Name.String())
	assert.Equal(t, 10, obj.Bases[0].Engineers)
	assert.Equal(t, 5, obj.Bases[0].Scientists)
	assert.True(t, obj.Bases[0].Active)
	assert.False(t, obj.Bases[1].Active)

	encoded, err := restruct.Pack(binary.LittleEndian, obj)
	require.NoError(t, err)
	assert.Equal(t, data, encoded, "round-trip should produce identical bytes")
}
