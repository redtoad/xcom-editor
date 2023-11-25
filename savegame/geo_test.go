package savegame_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/redtoad/xcom-editor/savegame"
)

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
			if got := savegame.Longitude(tt.x); got != tt.want {
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
			if got := savegame.Latitude(tt.y); got != tt.want {
				t.Errorf("Latitude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLocation(t *testing.T) {
	tests := []struct {
		name string
		x    int
		y    int
		want savegame.Location
	}{
		{
			"North pole",
			2160, -0,
			savegame.Location{90.0, 0.0},
		},
		{
			"South pole",
			0, -0,
			savegame.Location{-90.0, 0.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := savegame.NewLocation(tt.x, tt.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
