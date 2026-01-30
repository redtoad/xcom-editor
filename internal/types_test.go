package internal_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/stretchr/testify/assert"
)

func TestNullString_String(t *testing.T) {

	type nullStringStruct struct {
		Name internal.NullString `struct:"[26]byte"`
	}

	tests := []struct {
		bytes   []byte
		want    string
		wantErr bool
	}{
		{
			[]byte{0x4d, 0x69, 0x63, 0x68, 0x61, 0x65, 0x6c, 0x20, 0x53, 0x74, 0x65, 0x77, 0x61, 0x72, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			"Michael Stewart",
			false,
		},
		{
			[]byte{0x00, 0x69, 0x63, 0x68, 0x61, 0x65, 0x6c, 0x20, 0x53, 0x74, 0x65, 0x77, 0x61, 0x72, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%x", tt.bytes), func(t *testing.T) {
			var value nullStringStruct
			err := restruct.Unpack(tt.bytes, binary.LittleEndian, &value)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, value.Name.String())
		})
	}
}

func TestNullString_StringNoNull(t *testing.T) {
	// When no null byte is present, the full string is returned.
	s := internal.NullString("ABCDEF")
	assert.Equal(t, "ABCDEF", s.String())
}

func TestNullString_MultiField(t *testing.T) {

	// Verify NullString works correctly in a struct with multiple fields.
	type multiFieldStruct struct {
		ID   int16               `struct:"int16"`
		Name internal.NullString `struct:"[10]byte"`
		Age  int16               `struct:"int16"`
	}

	buf := []byte{
		0x01, 0x00, // ID = 1
		0x41, 0x6c, 0x69, 0x63, 0x65, 0x00, 0x00, 0x00, 0x00, 0x00, // Name = "Alice\0\0\0\0\0"
		0x1e, 0x00, // Age = 30
	}

	var value multiFieldStruct
	err := restruct.Unpack(buf, binary.LittleEndian, &value)
	assert.NoError(t, err)
	assert.Equal(t, int16(1), value.ID)
	assert.Equal(t, "Alice", value.Name.String())
	assert.Equal(t, int16(30), value.Age)

	// Round-trip: packing the unpacked struct must produce identical bytes.
	encoded, err := restruct.Pack(binary.LittleEndian, &value)
	assert.NoError(t, err)
	assert.Equal(t, buf, encoded)
}

func TestNullString_Pack(t *testing.T) {

	type nullStringStruct struct {
		Name string `struct:"[26]byte"`
	}

	tests := []struct {
		bytes   []byte
		want    string
		wantErr bool
	}{
		{
			[]byte{
				0x4d, 0x69, 0x63, 0x68, 0x61, 0x65, 0x6c, 0x20, 0x53, 0x74, 0x65, 0x77, 0x61, 0x72, 0x74, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			"Michael Stewart",
			false,
		},
		{
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%x", tt.bytes), func(t *testing.T) {
			value := nullStringStruct{Name: tt.want}
			data, _ := restruct.Pack(binary.LittleEndian, &value)
			assert.Equal(t, tt.bytes, data)

		})
	}
}
