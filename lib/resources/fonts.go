package resources

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
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
// See https://www.ufopaedia.org/index.php/SMALLSET.DAT for more info.
type Font struct {
	Characters [][]byte
	Masks      []image.Rectangle

	// Advance is the glyph advance, in pixels.
	Advance int
	// Width is the glyph width, in pixels.
	Width int
	// Height is the inter-line height, in pixels.
	Height int
	// Ascent is the glyph ascent, in pixels.
	Ascent int
	// Descent is the glyph descent, in pixels.
	Descent int
	// Left is the left side bearing, in pixels. A positive value means that
	// all of a glyph is to the right of the dot.
	Left int

	// Mask contains all of the glyph masks. Its width is typically the Face's
	// Width, and its height a multiple of the Face's Height.
	Mask image.Image

	// Ranges map runes to sub-images of Mask. The rune ranges must not
	// overlap, and must be in increasing rune order.
	Ranges []basicfont.Range
}

// Close implements the io.Close interface.
func (f Font) Close() error { return nil }

// Glyph returns the draw.DrawMask parameters (dr, mask, maskp) to draw r's
// glyph at the sub-pixel destination location dot, and that glyph's
// advance width.
//
// It returns !ok if the face does not contain a glyph for r. This includes
// returning !ok for a fallback glyph (such as substituting a U+FFFD glyph
// or OpenType's .notdef glyph), in which case the other return values may
// still be non-zero.
//
// The contents of the mask image returned by one Glyph call may change
// after the next Glyph call. Callers that want to cache the mask must make
// a copy.
func (f Font) Glyph(dot fixed.Point26_6, r rune) (
	dr image.Rectangle, mask image.Image, maskp image.Point, advance fixed.Int26_6, ok bool) {

	fmt.Printf("rune=%v -> ", r)
	if found, rng := f.find(r); rng != nil {
		maskp.Y = (int(found-rng.Low) + rng.Offset) * (f.Ascent + f.Descent)
		//x := int(dot.X+32)>>6 + f.Left
		//y := int(dot.Y+32) >> 6
		x := int(dot.X)>>6 + f.Left
		y := int(dot.Y) >> 6
		dr = image.Rectangle{
			Min: image.Point{
				X: x,
				Y: y - f.Ascent,
			},
			Max: image.Point{
				X: x + f.Width,
				Y: y + f.Descent,
			},
		}
		fmt.Printf("x=%v y=%v maskp=%v\n", x, y, maskp)
		return dr, f.Mask, maskp, fixed.I(f.Advance), r == found
	}
	return image.Rectangle{}, nil, image.Point{}, 0, false
}

// GlyphBounds returns the bounding box of r's glyph, drawn at a dot equal
// to the origin, and that glyph's advance width.
//
// It returns !ok if the face does not contain a glyph for r. This includes
// returning !ok for a fallback glyph (such as substituting a U+FFFD glyph
// or OpenType's .notdef glyph), in which case the other return values may
// still be non-zero.
//
// The glyph's ascent and descent are equal to -bounds.Min.Y and
// +bounds.Max.Y. The glyph's left-side and right-side bearings are equal
// to bounds.Min.X and advance-bounds.Max.X. A visual depiction of what
// these metrics are is at
// https://developer.apple.com/library/archive/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyphterms_2x.png
func (f Font) GlyphBounds(r rune) (bounds fixed.Rectangle26_6, advance fixed.Int26_6, ok bool) {
	if found, rng := f.find(r); rng != nil {
		return fixed.R(0, -f.Ascent, f.Width, +f.Descent), fixed.I(f.Advance), r == found
	}
	return fixed.Rectangle26_6{}, 0, false
}

// GlyphAdvance returns the advance width of r's glyph.
//
// It returns !ok if the face does not contain a glyph for r. This includes
// returning !ok for a fallback glyph (such as substituting a U+FFFD glyph
// or OpenType's .notdef glyph), in which case the other return values may
// still be non-zero.
func (f Font) GlyphAdvance(r rune) (advance fixed.Int26_6, ok bool) {
	if found, rng := f.find(r); rng != nil {
		return fixed.I(f.Advance), r == found
	}
	return 0, false
}

// Kern returns the horizontal adjustment for the kerning pair (r0, r1). A
// positive kern means to move the glyphs further apart.
func (f Font) Kern(r0, r1 rune) fixed.Int26_6 { return 1 }

// Metrics returns the metrics for this Face.
func (f Font) Metrics() font.Metrics {
	return font.Metrics{
		Height:     fixed.I(f.Height),
		Ascent:     fixed.I(f.Ascent),
		Descent:    fixed.I(f.Descent),
		XHeight:    fixed.I(f.Ascent),
		CapHeight:  fixed.I(f.Ascent),
		CaretSlope: image.Point{X: 0, Y: 1},
	}
}

func (f Font) find(r rune) (rune, *basicfont.Range) {
	for {
		for i, rng := range f.Ranges {
			if (rng.Low <= r) && (r < rng.High) {
				return r, &f.Ranges[i]
			}
		}
		if r == '\ufffd' {
			return 0, nil
		}
		r = '\ufffd'
	}
}

// min returns lowest value of a and b
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns highest value of a and b
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// FindBoundingBox returns the bounding box around a single character in an
// image as Rectangle. Part of See here for an illustration:
// https://developer.apple.com/library/archive/documentation/TextFonts/Conceptual/CocoaTextArchitecture/Art/glyphterms_2x.png
func FindBoundingBox(pixels []byte, width int) (bb image.Rectangle) {
	height := len(pixels) / width
	x0, y0, x1, y1 := 0, 0, width-1, height-1
	for idx := 0; idx < len(pixels); idx++ {
		if pixels[idx] > 0 {
			x := idx % width
			x0 = max(x0, x)
			x1 = min(x1, x)
			y := idx / width
			y0 = max(y0, y)
			y1 = min(y1, y)
		}
	}
	return image.Rect(x0, y0, x1, y1)
}

// yellowish font colour recreated from screenshots
var YellowFontPalette = color.Palette{
	color.Transparent,
	color.RGBA{219, 219, 134, 255},
	color.RGBA{178, 170, 101, 255},
	color.RGBA{142, 121, 73, 255},
	color.RGBA{101, 81, 48, 255},
	color.RGBA{64, 44, 28, 255},
}

// LoadFont loads a font from the game directory
func LoadFont(path string, width int, height int) (Font, error) {

	font := Font{
		Width:   width,
		Height:  height,
		Advance: 2,
		Ascent:  height - 2,
		Descent: 2,
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return font, fmt.Errorf("could not load font file: %w", err)
	}

	// The space character is not included in the data
	// so we simply add it at the beginning
	space := make([]byte, width*height)
	data = append(space, data...)

	font.Mask = &image.Paletted{
		Pix:     data,
		Stride:  width,
		Rect:    image.Rectangle{Max: image.Pt(width, len(data)/width)},
		Palette: YellowFontPalette,
	}

	/*

		// There is no space character (' ') in the DAT files, so we simply
		// add empty data at the beginning.
		font.Characters = append(font.Characters, make([]byte, width*height))
		font.Masks = append(font.Masks, image.Rectangle{Max: image.Point{4, height}})

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
			bb := FindBoundingBox(data, width)
			font.Masks = append(font.Masks, image.Rectangle{
				Max: image.Point{bb.Max.X, height},
			})
		}
	*/

	font.Ranges = []basicfont.Range{
		{Low: '\u0020', High: '\u007f', Offset: 0},
		// TODO ... add ranges for other characters ...
		{Low: '\ufffd', High: '\ufffd', Offset: 102}, // glyph for unsupported rune
	}

	return font, nil
}

// Draw generates a transparent image displaying a text in the specified Font.
func (f *Font) Draw(text string, palette color.Palette) (*image.Paletted, error) {

	d := font.Drawer{Face: f}
	bounds, advance := d.BoundString(text)
	fmt.Printf("bounds=%v advance=%s\n", bounds, advance)

	img := image.NewPaletted(
		image.Rect(0, 0, int(bounds.Max.X), int(bounds.Max.Y)),
		palette,
	)

	x, y := 0, bounds.Min.Y

	dot := fixed.Point26_6{X: fixed.I(x), Y: y}
	prevC := rune(-1)
	for _, c := range text {
		if prevC >= 0 {
			dot.X += f.Kern(prevC, c)
		}
		dr, mask, maskp, advance, _ := f.Glyph(dot, c)
		if !dr.Empty() {
			draw.DrawMask(img, dr, mask,
				image.Pt(dot.X.Round(), dot.Y.Round()),
				mask, maskp,
				draw.Over,
			)
		}
		dot.X += advance
		prevC = c
	}

	return img, nil
}
