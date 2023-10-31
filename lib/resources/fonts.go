package resources

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
)

// Font handles the font data loaded from the game directory.
//
// The Original UFO/TFTD engine uses two font types:
//
//   - big (BIGLETS.DAT, 16x16 pixels per character) and
//   - small (SMALLSET.DAT, 8x9 pixels per character).
//
// Each character is represented by an bitmap of up to six colors (including
// transparency), using one byte per pixel. These color indexes are remapped to
// different parts of the palette in use at run time (so text using this font
// may appear in many different shades or colors).
//
// The standard table starts at the first printable ASCII character, 33 ('!'),
// and continues to character 161 - For a total of 128 characters.
//
// See https://www.ufopaedia.org/index.php/SMALLSET.DAT for mor info.
type Font struct {
	Characters [][]byte
	Width      int
	Height     int
}

// LoadFont loads a font from the game directory
func LoadFont(path string, width int, height int) (Font, error) {

	font := Font{Width: width, Height: height}

	file, err := os.Open(path)
	if err != nil {
		return font, fmt.Errorf("could not load font file: %w", err)
	}

	for {
		data := make([]byte, width*height)
		bytesRead, err := file.Read(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			return font, fmt.Errorf("could not read data from file: %w", err)
		}
		if bytesRead != len(data) {
			return font, fmt.Errorf("chunk size does not match")
		}
		font.Characters = append(font.Characters, data)
	}

	return font, nil
}

// Text generates a transparent image displaying a text in the specified Font.
func Text(text string, font Font, palette color.Palette) (*image.Paletted, error) {
	width := len(text) * font.Width
	height := font.Height
	box := image.Rectangle{
		Min: image.Point{},                    // top left
		Max: image.Point{X: width, Y: height}, // bottom right
	}
	img := image.NewPaletted(box, palette)

	// Note: characters start first printable ASCII character 33 ('!')
	// so we need to remove this offset in teh calculations below.

	for pos, ch := range text {
		// ASCII character 32 is space (' ')
		if ch == 32 {
			continue
		}
		for idx, colour := range font.Characters[ch-33] {
			x := pos*font.Width + idx%font.Width
			y := idx / font.Width
			if colour == 0x00 {
				img.Set(x, y, image.Transparent)
			} else {
				img.Set(x, y, palette[colour])
			}
		}
	}
	return img, nil
}
