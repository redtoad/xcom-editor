package geoscape_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestSavegameInfo_Time(t *testing.T) {
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
			data, err := loadHex(tt.hex)
			assert.NoError(t, err)

			var info geoscape.SavegameInfo
			err = restruct.Unpack(data, binary.LittleEndian, &info)
			assert.NoError(t, err, "could not unpack test data: %v", err)

			assert.Equal(t, tt.want, info.Time())
		})
	}
}
