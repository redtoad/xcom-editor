package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMCD(t *testing.T) {
	buf := make([]byte, 2*MCDRecordSize)
	second := buf[MCDRecordSize:]
	second[0] = 7     // first frame
	second[1] = 8     // second frame
	second[20] = 0x34 // ScanG (little endian)
	second[21] = 0x12
	second[44] = 9    // DieMCD
	second[48] = 0xF8 // TLevel = -8
	second[49] = 4    // PLevel
	second[53] = TileNorthWall

	records, err := ParseMCD(buf)
	require.NoError(t, err)
	require.Len(t, records, 2)
	assert.Equal(t, MCDRecord{}, records[0])
	r := records[1]
	assert.Equal(t, [8]uint8{7, 8}, r.Frames)
	assert.Equal(t, uint16(0x1234), r.ScanG)
	assert.Equal(t, uint8(9), r.DieMCD)
	assert.Equal(t, int8(-8), r.TLevel)
	assert.Equal(t, uint8(4), r.PLevel)
	assert.Equal(t, uint8(TileNorthWall), r.TileType)

	_, err = ParseMCD(make([]byte, MCDRecordSize+1))
	assert.Error(t, err)
}
