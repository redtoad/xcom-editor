package resources

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func readInt16(rd *bufio.Reader) (int, error) {
	buf := make([]byte, 2)
	_, err := rd.Read(buf)
	if err != nil {
		return -1, err
	}
	return int(binary.LittleEndian.Uint16(buf)), nil
}

// LoadSPK loads SPK image from path.
//
// Another 320 pixels wide x 200 pixels high image format but using compression, primarily used by UFO for
// background images (eg inventory screens).
//
// https://www.ufopaedia.org/index.php/Image_Formats
func LoadSPK(path string) (*ImageResource, error) {

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	buf := bufio.NewReader(fp)

	var spriteData []uint

	for {

		// Read a 16-bit unsigned integer, call it "a".
		value, err := readInt16(buf)
		if err != nil {
			return nil, err
		}

		switch value {

		// If "a" is 0xFFFF (65535) then you must read another 16-bit
		// integer, and skip that number*2 pixels.
		case 0xffff:
			pixels, err := readInt16(buf)
			if err != nil {
				return nil, err
			}
			for i := 0; i < pixels*2; i++ {
				spriteData = append(spriteData, 0) // transparent background
			}

		// If "a" is 0xFFFE (65534) then the next 16-bit integer*2
		// specifies the number of pixels you are going to draw.
		// Read that number of bytes from the file and draw their indexed
		// color.
		case 0xfffe:
			pixels, err := readInt16(buf)
			if err != nil {
				return nil, err
			}
			colors := make([]byte, pixels*2)
			_, err = io.ReadFull(buf, colors)

			if err != nil {
				return nil, err
			}
			for i := 0; i < len(colors); i++ {
				spriteData = append(spriteData, uint(colors[i]))
			}

		// If "a" is 0xFFFD (65533) then you are done. This is always the
		// last code in the file.
		case 0xfffd:
			return &ImageResource{spriteData, 320}, nil

		default:
			return nil, fmt.Errorf("unknown header byte %x", value)
		}

	}

}
