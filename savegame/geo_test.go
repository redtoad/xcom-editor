package savegame_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/redtoad/xcom-editor/savegame"
)

func TestNewCoord(t *testing.T) {
	tests := []struct {
		x    int
		y    int
		want savegame.Coord
	}{
		{0, -0, savegame.Coord{}},
		{2160, 720, savegame.Coord{Lat: -90.0, Lon: -90.0}},
		{2879, -720, savegame.Coord{Lat: 90.0, Lon: -0.125}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("(x=%d,y=%d)", tt.x, tt.y), func(t *testing.T) {
			if got := savegame.NewCoord(tt.x, tt.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCoord() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCoord_String(t *testing.T) {
	tests := []struct {
		coord savegame.Coord
		want  string
	}{
		{savegame.Coord{}, "0.00000° N 0.00000° E"},
		{savegame.Coord{Lat: 90.0, Lon: 270.0}, "90.00000° N 270.00000° E"},
		{savegame.Coord{Lat: 90.0, Lon: -90.0}, "90.00000° N 90.00000° W"},
		{savegame.Coord{Lat: -90.0, Lon: 359.875}, "90.00000° S 359.87500° E"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("(lat=%f,lon=%f)", tt.coord.Lat, tt.coord.Lon), func(t *testing.T) {
			if got := tt.coord.String(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf(".String() = %v, want %v", got, tt.want)
			}
		})
	}
}
