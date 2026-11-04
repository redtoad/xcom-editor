package resources

import (
	"encoding/binary"
	"image"
	"image/color"
	"io/ioutil"

	"github.com/go-restruct/restruct"
)

// LoadPalettes loads Palettes from path.
// (see https://www.ufopaedia.org/index.php/PALETTES.DAT)
func LoadPalettes(path string) ([]*Palette, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	palettes := make([]*Palette, 5)
	for i := 0; i < 5; i++ {
		palette := &Palette{}
		data := buf[i*palette.SizeOf() : (i+1)*palette.SizeOf()]
		if err = restruct.Unpack(data, binary.LittleEndian, &palette); err != nil {
			return nil, err
		}
		palettes[i] = palette
	}
	return palettes, nil
}

// Palette stores colours used for game images.
//
// Back in the days of low video memory, storing each pixel as a multi-byte color value (such as in
// today's "High Color" and "True Color" modes common under MS Windows) was not practical. Instead,
// a common method for was to create a palette of 256 different colors (usually of 3 bytes each), and
// then use single byte values to index into that. The ability of a video card to use such a palette
// was known as "VGA compatibility".
//
// This dramatically decreased the memory requirements of a given image, though it lowered the amount
// of colors that could be used in any single moment.
type Palette struct {
	Colors [256]color.Color `struct:"[768]byte"`
	Buffer []byte           `struct:"[6]byte"`
}

// SizeOf implements the restruct.Sizer interface
func (p Palette) SizeOf() int {
	// Each palette is made up of the standard 256 colors stored as three byte RGB records (a value
	// for Red, a value for Green and a value for Blue). This makes a total of 768 bytes per palette.
	// Each palette has a six byte 'buffer' following it, whose purpose is unknown.
	return 768 + 6
}

// Unpack implements the restruct.Unpacker interface
func (p *Palette) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < 256; i++ {
		rgb := struct{ R, G, B uint8 }{}
		data := buf[i*3 : (i+1)*3]
		if err := restruct.Unpack(data, order, &rgb); err != nil {
			return nil, err
		}

		// Each RGB value has a maximum intensity of 63/x3F (making a total of 262,144 shades
		// available, which seems to be standard for VGA), as opposed to 255/xFF (which would
		// have provided 16,777,216). If you want to use the 64 intensity values in a 256
		// intensity world, just multiply each value by 4 and you'll achieve the correct
		// strength (or near enough - at least, this is how the CE versions of the game
		// deal with the issue).
		p.Colors[i] = color.NRGBA{R: rgb.R * 4, G: rgb.G * 4, B: rgb.B * 4, A: 0xff}
	}
	p.Buffer = buf[768 : 768+6]
	return buf[p.SizeOf():], nil
}

// Palette returns the colors as color.Palette object.
func (p *Palette) Palette() *color.Palette {
	palette := make(color.Palette, len(p.Colors))
	palette[0] = image.Transparent
	for i := 1; i < len(p.Colors); i++ {
		palette[i] = p.Colors[i]
	}
	return &palette
}

// BattlescapePaletteIndex is the palette of the Battlescape in PALETTES.DAT.
const BattlescapePaletteIndex = 4

// battlescapeGradient replaces the last 16 colors of the Battlescape palette.
// The game sets them at runtime to a greyish gradient (used e.g. by the
// X-COM crafts, roads and armours), PALETTES.DAT holds other colors there.
// The values are the same as in OpenXcom.
var battlescapeGradient = [16]color.NRGBA{
	{140, 152, 148, 255}, {132, 136, 140, 255}, {116, 124, 132, 255}, {108, 116, 124, 255},
	{92, 104, 108, 255}, {84, 92, 100, 255}, {76, 80, 92, 255}, {56, 68, 84, 255},
	{48, 56, 68, 255}, {40, 48, 56, 255}, {32, 36, 48, 255}, {24, 28, 32, 255},
	{16, 20, 24, 255}, {8, 12, 16, 255}, {3, 4, 8, 255}, {3, 3, 6, 255},
}

// BattlescapePalette returns the palette used by the Battlescape: the
// Battlescape palette of PALETTES.DAT with the last 16 colors replaced by
// the gradient the game uses.
func BattlescapePalette(palettes []*Palette) *Palette {
	result := *palettes[BattlescapePaletteIndex]
	for i, c := range battlescapeGradient {
		result.Colors[256-len(battlescapeGradient)+i] = c
	}
	return &result
}
