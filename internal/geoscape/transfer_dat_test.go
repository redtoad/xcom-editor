package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_TRANSFER_DAT_RoundTrip(t *testing.T) {
	// Create a buffer of 800 bytes (100 entries * 8 bytes each), all zeros
	data := make([]byte, 800)

	// Set up a couple of active transfers
	// Transfer 0: Origin=0, Destination=1, HoursLeft=10, Type=2, ReferenceNumber=5, Quantity=3
	data[0] = 0x00 // Origin
	data[1] = 0x01 // Destination
	data[2] = 0x0A // HoursLeft
	data[3] = 0x02 // Type
	data[4] = 0x05 // ReferenceNumber low byte
	data[5] = 0x00 // ReferenceNumber high byte
	data[6] = 0x03 // Quantity
	data[7] = 0x00 // padding

	var obj geoscape.TransferFile
	err := restruct.Unpack(data, binary.LittleEndian, &obj)
	require.NoError(t, err)
	assert.Equal(t, 100, len(obj.Transfers))
	assert.Equal(t, 0, obj.Transfers[0].Origin)
	assert.Equal(t, 1, obj.Transfers[0].Destination)
	assert.Equal(t, 10, obj.Transfers[0].HoursLeft)
	assert.Equal(t, 2, obj.Transfers[0].Type)
	assert.Equal(t, 5, obj.Transfers[0].ReferenceNumber)
	assert.Equal(t, 3, obj.Transfers[0].Quantity)

	encoded, err := restruct.Pack(binary.LittleEndian, obj)
	require.NoError(t, err)
	assert.Equal(t, data, encoded, "round-trip should produce identical bytes")
}
