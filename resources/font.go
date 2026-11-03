package resources

import (
	"fmt"
	"os"
)

// Fonts of the game (see https://www.ufopaedia.org/index.php/BIGLETS.DAT).
var (
	// SmallFont is GEODATA/SMALLSET.DAT, used for most texts.
	SmallFont = FontFormat{Path: "GEODATA/SMALLSET.DAT", Width: 8, Height: 9, Baseline: 7}
	// BigFont is GEODATA/BIGLETS.DAT, used for titles.
	BigFont = FontFormat{Path: "GEODATA/BIGLETS.DAT", Width: 16, Height: 16, Baseline: 12}
)

// firstChar is the character of the first glyph in a font file.
const firstChar = 33

// FontFormat describes a bitmap font file of the game.
type FontFormat struct {
	Path string
	// Width and Height of each glyph in pixels.
	Width, Height int
	// Baseline is the row (from the top) the characters stand on.
	Baseline int
}

// BitmapFont is a font of the game. Each glyph is a bitmap with one byte
// (a palette offset, 0 = transparent) per pixel. Glyph i is the character
// 33+i in code page 437.
//
// The letters are drawn with the offsets 1-3 and surrounded by an outline
// (offsets 4 and 5), which the game draws in a dark color.
type BitmapFont struct {
	FontFormat
	Glyphs [][]uint8
}

// LoadBitmapFont loads a font from the game root folder.
func LoadBitmapFont(root string, format FontFormat) (*BitmapFont, error) {
	buf, err := os.ReadFile(root + "/" + format.Path)
	if err != nil {
		return nil, err
	}
	return ParseBitmapFont(buf, format)
}

// ParseBitmapFont parses the glyphs of a font file.
func ParseBitmapFont(buf []byte, format FontFormat) (*BitmapFont, error) {
	size := format.Width * format.Height
	if len(buf) == 0 || len(buf)%size != 0 {
		return nil, fmt.Errorf("invalid font size %d (glyphs have %d bytes)", len(buf), size)
	}
	font := &BitmapFont{FontFormat: format}
	for i := 0; i < len(buf); i += size {
		font.Glyphs = append(font.Glyphs, buf[i:i+size])
	}
	return font, nil
}

// Rune returns the character of glyph i.
func (f *BitmapFont) Rune(i int) rune {
	return cp437(firstChar + i)
}

// Palette offsets of the pixels of a glyph.
const (
	// shadedOffset is a darker pixel of a letter. At the sides of a letter
	// it rounds corners and marks notches (e.g. it is the only difference
	// between A and R).
	shadedOffset = 2
	// outlineOffset is the first offset of the outline around the letters.
	outlineOffset = 4
)

// Set returns whether pixel (x, y) of glyph i belongs to the letter (not to
// its outline). Shaded pixels at the left or right side of a letter are left
// out, so the letters keep their shape when drawn in a single color.
func (f *BitmapFont) Set(i, x, y int) bool {
	if !f.inLetter(i, x, y) {
		return false
	}
	if f.Glyphs[i][y*f.Width+x] == shadedOffset &&
		(!f.inLetter(i, x-1, y) || !f.inLetter(i, x+1, y)) {
		return false
	}
	return true
}

// inLetter returns whether pixel (x, y) of glyph i is part of the letter
// (including shaded pixels).
func (f *BitmapFont) inLetter(i, x, y int) bool {
	if x < 0 || y < 0 || x >= f.Width || y >= f.Height {
		return false
	}
	v := f.Glyphs[i][y*f.Width+x]
	return v != 0 && v < outlineOffset
}

// GlyphWidth returns the width of glyph i in pixels (up to its rightmost set
// pixel).
func (f *BitmapFont) GlyphWidth(i int) int {
	width := 0
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			if f.Set(i, x, y) && x+1 > width {
				width = x + 1
			}
		}
	}
	return width
}

// cp437High maps the characters 128-255 of code page 437 to Unicode.
var cp437High = []rune("ÇüéâäàåçêëèïîìÄÅÉæÆôöòûùÿÖÜ¢£¥₧ƒáíóúñÑªº¿⌐¬½¼¡«»░▒▓│┤╡╢╖╕╣║╗╝╜╛┐└┴┬├─┼╞╟╚╔╩╦╠═╬╧╨╤╥╙╘╒╓╫╪┘┌█▄▌▐▀αßΓπΣσµτΦΘΩδ∞φε∩≡±≥≤⌠⌡÷≈°∙·√ⁿ²■ ")

// cp437 converts a character of code page 437 to Unicode.
func cp437(c int) rune {
	if c < 128 {
		return rune(c)
	}
	return cp437High[c-128]
}
