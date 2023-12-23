package geodata_test

import (
	"encoding/base64"
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geodata"
)

func MustDecode(txt string) []byte {
	data, err := base64.StdEncoding.DecodeString(txt)
	if err != nil {
		panic(err)
	}
	return data
}

func TestWORLD_DAT_Unpack(t *testing.T) {
	tests := []struct {
		name    string
		buffer  []byte
		want    []geodata.Polygon
		wantErr bool
	}{
		{
			"read one polygon",
			MustDecode("7Ape/vQKZf4MC17+9gpW/gEAAAA="),
			[]geodata.Polygon{
				{2796, -418, 2804, -411, 2828, -418, 2806, -426, 1},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := &geodata.WorldData{}
			err := restruct.Unpack(tt.buffer, binary.LittleEndian, &world)
			if (err != nil) != tt.wantErr {
				t.Errorf("WORLD_DAT.Unpack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(world.Polygons, tt.want) {
				t.Errorf("WORLD_DAT.Unpack() = %v, want %v", world.Polygons, tt.want)
			}
		})
	}
}
