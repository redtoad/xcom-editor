package resources_test

import (
	"testing"

	"github.com/redtoad/xcom-editor/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// testFont is a 4x3 font (baseline 2) with the glyphs '!' and '"'. The
// values 4 and 5 are the outline of the letters, which is not part of the
// glyphs.
var testFont = resources.FontFormat{Width: 4, Height: 3, Baseline: 2}

func testFontData() []byte {
	return []byte{
		// '!': a 2x2 square with an outline
		1, 1, 4, 0,
		1, 1, 5, 0,
		4, 4, 4, 0,
		// '"': an L shape with a descender
		0, 0, 0, 0,
		0, 1, 0, 0,
		0, 1, 1, 1,
	}
}

func TestParseBitmapFont(t *testing.T) {
	f, err := resources.ParseBitmapFont(testFontData(), testFont)
	require.NoError(t, err)
	require.Len(t, f.Glyphs, 2)

	assert.Equal(t, '!', f.Rune(0))
	assert.Equal(t, '"', f.Rune(1))
	assert.True(t, f.Set(0, 1, 1))
	assert.False(t, f.Set(0, 2, 1), "outline")
	assert.Equal(t, 2, f.GlyphWidth(0))
	assert.Equal(t, 4, f.GlyphWidth(1))

	_, err = resources.ParseBitmapFont([]byte{1, 2, 3}, testFont)
	assert.Error(t, err)
}

func TestBitmapFont_ShadedCorners(t *testing.T) {
	// the A and R of the small font only differ in shaded (2) pixels at the
	// sides (top row and middle bar)
	glyphs := map[rune][]byte{
		'A': {4, 2, 1, 1, 2, 4, 0, 0, 4, 1, 4, 4, 1, 4, 0, 0, 4, 1, 1, 1, 1, 4, 0, 0},
		'R': {4, 1, 1, 1, 2, 4, 0, 0, 4, 1, 4, 4, 1, 4, 0, 0, 4, 1, 1, 1, 2, 4, 0, 0},
	}
	format := resources.FontFormat{Width: 8, Height: 3, Baseline: 3}
	rows := func(r rune) []string {
		f, err := resources.ParseBitmapFont(glyphs[r], format)
		require.NoError(t, err)
		var result []string
		for y := 0; y < 3; y++ {
			row := ""
			for x := 0; x < 6; x++ {
				if f.Set(0, x, y) {
					row += "#"
				} else {
					row += "."
				}
			}
			result = append(result, row)
		}
		return result
	}
	assert.Equal(t, []string{"..##..", ".#..#.", ".####."}, rows('A'))
	assert.Equal(t, []string{".###..", ".#..#.", ".###.."}, rows('R'))
}

func TestBitmapFont_RuneCodePage437(t *testing.T) {
	f := &resources.BitmapFont{}
	assert.Equal(t, 'A', f.Rune('A'-33))
	assert.Equal(t, 'Ç', f.Rune(128-33))
	assert.Equal(t, 'ü', f.Rune(129-33))
	assert.Equal(t, 'ß', f.Rune(225-33))
}

func TestBitmapFont_TTF(t *testing.T) {
	f, err := resources.ParseBitmapFont(testFontData(), testFont)
	require.NoError(t, err)

	parsed, err := sfnt.Parse(f.TTF("Test Font"))
	require.NoError(t, err)

	var buf sfnt.Buffer
	name, err := parsed.Name(&buf, sfnt.NameIDFamily)
	require.NoError(t, err)
	assert.Equal(t, "Test Font", name)

	// one font pixel = 128 units, the em is the height of the font
	assert.Equal(t, fixed.Int26_6(3*128), fixed.Int26_6(parsed.UnitsPerEm()))
	ppem := fixed.I(3 * 128) // 1 unit = 1/64 pixel
	metrics, err := parsed.Metrics(&buf, ppem, font.HintingNone)
	require.NoError(t, err)
	assert.Equal(t, fixed.I(2*128), metrics.Ascent)
	assert.Equal(t, fixed.I(1*128), metrics.Descent)

	for _, tc := range []struct {
		r       rune
		advance int
		bounds  fixed.Rectangle26_6
	}{
		// the square stands on the baseline
		{'!', 2, fixed.R(0, -2*128, 2*128, 0)},
		// the L shape reaches below the baseline
		{'"', 4, fixed.R(1*128, -1*128, 4*128, 1*128)},
	} {
		idx, err := parsed.GlyphIndex(&buf, tc.r)
		require.NoError(t, err)
		require.NotZero(t, idx, "glyph for %q", tc.r)

		advance, err := parsed.GlyphAdvance(&buf, idx, ppem, font.HintingNone)
		require.NoError(t, err)
		assert.Equal(t, fixed.I(tc.advance*128), advance, "advance of %q", tc.r)

		bounds, _, err := parsed.GlyphBounds(&buf, idx, ppem, font.HintingNone)
		require.NoError(t, err)
		assert.Equal(t, tc.bounds, bounds, "bounds of %q", tc.r)
	}

	space, err := parsed.GlyphIndex(&buf, ' ')
	require.NoError(t, err)
	assert.NotZero(t, space)
}
