package resources

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBattlescapePalette(t *testing.T) {
	palettes := make([]*Palette, 5)
	for i := range palettes {
		palettes[i] = &Palette{}
		for c := range palettes[i].Colors {
			palettes[i].Colors[c] = color.NRGBA{B: 255, A: 255}
		}
	}
	p := BattlescapePalette(palettes)
	assert.Equal(t, color.NRGBA{B: 255, A: 255}, p.Colors[239], "colors before the gradient are kept")
	assert.Equal(t, color.NRGBA{140, 152, 148, 255}, p.Colors[240])
	assert.Equal(t, color.NRGBA{3, 3, 6, 255}, p.Colors[255])
	assert.Equal(t, color.NRGBA{B: 255, A: 255}, palettes[4].Colors[240], "PALETTES.DAT is not changed")
}
