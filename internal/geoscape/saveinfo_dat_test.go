package geoscape_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestSaveinfo_Unpack(t *testing.T) {
	buf := MustLoadFromBase64(testFile_SAVEINFO_DAT)
	var bf geoscape.SaveinfoFile
	err := restruct.Unpack(buf, binary.LittleEndian, &bf)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSaveinfoFile_Time(t *testing.T) {
	tests := []struct {
		hex  string
		want time.Time
	}{
		{
			"01005465737400000000000000000000000000000000000000000000cf0703000d00040014000000",
			time.Date(1999, 4, 13, 4, 20, 0, 0, time.UTC),
		},
		{
			"010057656c6c206f6e20746865207761790000000000000000000000cf0703000e0002002f000000",
			time.Date(1999, 4, 14, 2, 47, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.hex, func(t *testing.T) {
			data := MustLoadFromHex(tt.hex)

			var info geoscape.SaveinfoFile
			err := restruct.Unpack(data, binary.LittleEndian, &info)
			assert.NoError(t, err, "could not unpack test data: %v", err)

			assert.Equal(t, tt.want, info.Time())
		})
	}
}

// base64 -i GAME_1/SAVEINFO.DAT -b 120 | pbcopy
const testFile_SAVEINFO_DAT = `
AQBUZXN0AACgExwAaAEBAEYCAAC2+hkAAAAAAM8HBwANAA8AJAAAAA==
`
