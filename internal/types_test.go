package internal_test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/stretchr/testify/assert"
)

func TestNullString_Unpack(t *testing.T) {

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
			_ = restruct.Unpack(tt.bytes, binary.LittleEndian, &value)
			assert.Equal(t, tt.want, string(value.Name))
		})
	}
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
