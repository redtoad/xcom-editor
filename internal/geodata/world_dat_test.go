package geodata_test

import (
	"encoding/base64"
	"encoding/binary"
	"image/color"
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
			world := &geodata.WorldFile{}
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

func TestTextureFile_Image(t *testing.T) {
	var buf [geodata.NumTextureSets * geodata.NumTerrains * geodata.TextureSize * geodata.TextureSize]byte
	// first pixel of terrain 2 in texture set 1
	buf[(1*geodata.NumTerrains+2)*geodata.TextureSize*geodata.TextureSize] = 7

	var f geodata.TextureFile
	if err := restruct.Unpack(buf[:], binary.LittleEndian, &f); err != nil {
		t.Fatal(err)
	}
	palette := make(color.Palette, 256)
	for i := range palette {
		palette[i] = color.Gray{Y: uint8(i)}
	}
	img := f.Image(1, 2, palette)
	if got := img.ColorIndexAt(0, 0); got != 7 {
		t.Errorf("ColorIndexAt(0, 0) = %d, want 7", got)
	}
	if got := img.Bounds().Dx(); got != geodata.TextureSize {
		t.Errorf("width = %d, want %d", got, geodata.TextureSize)
	}
}
