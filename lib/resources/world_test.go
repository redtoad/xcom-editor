package resources

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-restruct/restruct"
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
		want    []Polygon
		wantErr bool
	}{
		{
			"read one polygon",
			MustDecode("7Ape/vQKZf4MC17+9gpW/gEAAAA="),
			[]Polygon{
				{2796, -418, 2804, -411, 2828, -418, 2806, -426, 1},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := &WorldData{}
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

func TestLongitude(t *testing.T) {
	tests := []struct {
		x    int
		want float32
	}{
		{0, 0.0},
		{2160, 270},
		{2879, 359.875},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d -> %f", tt.x, tt.want), func(t *testing.T) {
			if got := Longitude(tt.x); got != tt.want {
				t.Errorf("Longitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLatitude(t *testing.T) {
	tests := []struct {
		y    int
		want float32
	}{
		{0, 0.0},
		{720, +90.0},
		{-720, -90.0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d -> %f", tt.y, tt.want), func(t *testing.T) {
			if got := Latitude(tt.y); got != tt.want {
				t.Errorf("Latitude() = %v, want %v", got, tt.want)
			}
		})
	}
}
