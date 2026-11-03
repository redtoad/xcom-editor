package geodata

import (
	"image"
	"image/color"
)

const (
	// NumTerrains is the number of terrain types in WORLD.DAT (Polygon.Terrain).
	NumTerrains = 13
	// NumTextureSets is the number of texture sets (one per zoom level).
	NumTextureSets = 3
	// TextureSize is the width and height of a texture in pixels.
	TextureSize = 32
)

// TextureFile holds the textures used to draw the terrain of the globe
// (GEOGRAPH/TEXTURE.DAT). There are 3 sets of 13 textures (one per terrain
// type, see Polygon.Terrain), one set for each zoom level. Each texture is
// 32x32 pixels, one byte (palette index) per pixel.
//
// The file is 39936 bytes long.
//
// https://www.ufopaedia.org/index.php/TEXTURE.DAT
type TextureFile struct {
	Textures [NumTextureSets * NumTerrains][TextureSize * TextureSize]uint8
}

// Image returns the texture for the given terrain type and texture set
// (zoom level) using the given palette.
func (f *TextureFile) Image(set, terrain int, palette color.Palette) *image.Paletted {
	img := image.NewPaletted(image.Rect(0, 0, TextureSize, TextureSize), palette)
	copy(img.Pix, f.Textures[set*NumTerrains+terrain][:])
	return img
}
