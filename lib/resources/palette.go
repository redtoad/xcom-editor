package resources

import (
	"encoding/binary"
	"image"
	"image/color"
	"io/ioutil"

	"gopkg.in/restruct.v1"
)

// https://www.ufopaedia.org/index.php/PALETTES.DAT

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

func (p Palette) SizeOf() int {
	// Each palette is made up of the standard 256 colors stored as three byte RGB records (a value
	// for Red, a value for Green and a value for Blue). This makes a total of 768 bytes per palette.
	// Each palette has a six byte 'buffer' following it, whose purpose is unknown.
	return 768 + 6
}

// Implements the restruct.Unpacker interface
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

func (p *Palette) Palette() *color.Palette {
	palette := make(color.Palette, len(p.Colors))
	palette[0] = image.Transparent
	for i := 1; i < len(p.Colors); i++ {
		palette[i] = p.Colors[i]
	}
	return &palette
}
