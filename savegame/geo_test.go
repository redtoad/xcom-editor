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
		{0, -0, savegame.Coord{0.0, 0.0}},
		{2160, 720, savegame.Coord{270.0, 90.0}},
		{2879, -720, savegame.Coord{359.875, -90.0}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("(x=%d,y=%d)", tt.x, tt.y), func(t *testing.T) {
			if got := savegame.NewCoord(tt.x, tt.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}
