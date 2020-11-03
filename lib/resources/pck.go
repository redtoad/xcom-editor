package resources

import (
	"bufio"
	"encoding/binary"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"log"
	"math"
	"os"
)

var (
	// ErrUnsupportedTABOffset is thrown if the offset is not 2 or 4
	ErrUnsupportedTABOffset = errors.New("unsupported offset size for TAB file")
	// ErrNotEnoughSprites is thrown if no (or not enough) images are found in ImageCollection.
	ErrNotEnoughSprites = errors.New("could not load sprites from file")
	// ErrNotEnoughTabs is thrown if not enough offsets are found to load PCK file.
	ErrNotEnoughTabs = errors.New("could not load offsets from tab file")
)

//  32x40 sprites for use in the BattleScape
//  32x48 sprites in BIGOBS.PCK
// 128x68 sprite for X1.PCK

func readInt(rd *bufio.Reader) (int, error) {
	value, err := rd.ReadByte()
	if err != nil {
		return 0, err
	}
	return int(value), nil
}

// LoadTAB loads TAB file from path.
//
// The TAB file is a list of file offsets saying where each image begins in the related PCK archive.
// Some offsets are encoded in 2 bytes while others are encoded in 4. It depends on how many images
// are in the PCK archive.
func LoadTAB(path string, offsetSize int) ([]int, error) {
	var offsets []int
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(fp)

	for {

		values, err := buf.Peek(offsetSize)
		if err != nil {
			if err == io.EOF {
				return offsets, nil
			}
			return nil, err
		}

		switch offsetSize {
		case 2:
			offset := binary.LittleEndian.Uint16(values)
			offsets = append(offsets, int(offset))
		case 4:
			offset := binary.LittleEndian.Uint32(values)
			offsets = append(offsets, int(offset))
		default:
			return nil, ErrUnsupportedTABOffset
		}

		_, err = buf.Discard(offsetSize)
		if err != nil {
			return nil, err
		}
	}

}

// LoadPCK loads PCK image from path.
func LoadPCK(path string, width int, offset int) (*ImageResource, error) {

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(fp)
	_, err = buf.Discard(offset)
	if err != nil {
		return nil, err
	}

	var spriteData []uint

	// read the first byte and skip down that many rows
	skipRows, err := readInt(buf)
	if err != nil {
		return nil, err
	}
	for i := 0; i < width*skipRows; i++ {
		spriteData = append(spriteData, 0)
	}

	for {

		value, err := readInt(buf)
		if err != nil {
			return nil, err
		}

		switch value {

		// If the byte is 254/xFE, the byte following is the number of transparent pixels.
		// Any other byte is a color index (with the exception of 255/xFF which signals the
		// end of that sprite's data).

		case 0xfe: // n transparent pixels follow
			pixels, err := readInt(buf)
			if err != nil {
				return nil, err
			}
			for i := 0; i < pixels; i++ {
				spriteData = append(spriteData, 0)
			}

		case 0xff: // end of data
			return &ImageResource{spriteData, width}, nil

		default: // normal pixel in pallet's color
			spriteData = append(spriteData, uint(value))
		}
	}

}

// ImageResource is a list of pixels (colors) that will create an image of width.
type ImageResource struct {
	pixels []uint
	width  int
}

// Image will create an Image object with colors from palette.
func (i *ImageResource) Image(palette *Palette) image.Image {
	height := len(i.pixels) / i.width
	upLeft := image.Point{}
	lowRight := image.Point{X: i.width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for idx, col := range i.pixels {
		x := idx % i.width
		y := idx / i.width
		if col == 0 {
			img.Set(x, y, image.Transparent)
		} else {
			img.Set(x, y, palette.Colors[col])
		}
	}
	return img
}

// Paletted returns an image with a limited palette (used for GIFs).
func (i *ImageResource) Paletted(palette *Palette) *image.Paletted {
	height := len(i.pixels) / i.width
	upLeft := image.Point{}
	lowRight := image.Point{X: i.width, Y: height}
	img := image.NewPaletted(image.Rectangle{Min: upLeft, Max: lowRight}, *palette.Palette())
	for idx, col := range i.pixels {
		x := idx % i.width
		y := idx / i.width
		if col == 0 {
			img.Set(x, y, image.Transparent)
		} else {
			img.Set(x, y, palette.Colors[col])
		}
	}
	return img
}

// Width returns the image's width in pixels.
func (i *ImageResource) Width() int {
	return i.width
}

// Height calculates the image's height in pixels.
func (i *ImageResource) Height() int {
	return len(i.pixels) / i.width
}

// LoadImageCollectionFromPCK loads a collection of images from a PCK file.
func LoadImageCollectionFromPCK(path string, width int, offsets []int) (*ImageCollection, error) {

	if len(offsets) == 0 {
		return nil, ErrNotEnoughTabs
	}

	var sprites []*ImageResource
	for _, tab := range offsets {
		sprite, err := LoadPCK(path, width, tab)
		if err != nil {
			return nil, err
		}
		sprites = append(sprites, sprite)

	}
	return &ImageCollection{Sprites: sprites, SpriteWidth: width}, nil
}

// ImageCollection is a list of ImageResources. These will typically be turned into
// a gallery or an animated image.
type ImageCollection struct {
	Sprites     []*ImageResource
	SpriteWidth int
}

// gridSize calculates width and height of grid for collection of images
func gridSize(sprites, perRow int) (width int, height int) {
	if sprites > perRow {
		width = perRow
	} else {
		width = sprites
	}
	height = int(math.Ceil(float64(sprites) / float64(width)))
	return
}

// Gallery creates a collection of all images on a grid with numberPerRow images in each
// row. The final size will depend on the number of Sprites in the collection (making up
// the grid), the size of the single SpriteWidth and the rowHeight.
func (c *ImageCollection) Gallery(numberPerRow int, rowHeight int, palette *Palette) (image.Image, error) {

	if len(c.Sprites) == 0 {
		return nil, ErrNotEnoughSprites
	}

	// create new image with black background
	gridWidth, gridHeight := gridSize(len(c.Sprites), numberPerRow)
	collection := image.NewRGBA(image.Rect(
		0, 0,
		gridWidth*c.SpriteWidth, gridHeight*rowHeight))

	// draw each image onto collection
	for no, sprite := range c.Sprites {

		if sprite.Height() > rowHeight {
			log.Printf("Warning: sprite %dx%d is bigger than specified in meta data (%dx%d)!\n",
				sprite.Width(), sprite.Height(), c.SpriteWidth, rowHeight)
		}

		dstX := (no % gridWidth) * c.SpriteWidth
		dstY := (no / gridWidth) * rowHeight
		dstR := image.Rectangle{
			Min: image.Point{X: dstX, Y: dstY},
			Max: image.Point{X: dstX + sprite.Width(), Y: dstY + sprite.Height()},
		}

		img := sprite.Image(palette)
		draw.Draw(collection, dstR,
			img, img.Bounds().Min,
			draw.Src)
	}

	return collection, nil
}

// Animated converts ImageCollection into a single GIF.
func (c *ImageCollection) Animated(delay int, height int, palette *Palette) *gif.GIF {

	sprites := make([]*image.Paletted, len(c.Sprites))
	for i, sprite := range c.Sprites {
		sprites[i] = sprite.Paletted(palette)
	}

	delays := make([]int, len(sprites))
	disposal := make([]byte, len(sprites))
	for i := 0; i < len(delays); i++ {
		delays[i] = delay
		disposal[i] = gif.DisposalBackground
	}

	return &gif.GIF{
		Image:           sprites,
		Delay:           delays,
		LoopCount:       0,
		Disposal:        disposal,
		BackgroundIndex: 0,
		Config:          image.Config{Width: c.SpriteWidth, Height: height},
	}
}
