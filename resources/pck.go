package resources

import (
	"bufio"
	"encoding/binary"
	"errors"
	"image"
	"io"
	"os"
)

var ErrUnsupportedTABOffset = errors.New("unsupported offset size for TAB file")

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
			} else {
				return nil, err
			}
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

type ImageResource struct {
	pixels []uint
	width  int
}

func (i *ImageResource) Image(palette *Palette) image.Image {
	height := len(i.pixels) / i.width
	upLeft := image.Point{}
	lowRight := image.Point{X: i.width, Y: height}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
	for idx, col := range i.pixels {
		x := idx % i.width
		y := (idx / i.width)
		if col == 0 {
			img.Set(x, y, image.Transparent)
		} else {
			img.Set(x, y, palette.Colors[col])
		}
	}
	return img
}

func (i *ImageResource) Width() int {
	return i.width
}

func (i *ImageResource) Height() int {
	return len(i.pixels) / i.width
}
