package resources

import (
	"bufio"
	"encoding/binary"
	"log"
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

	pxCount := 0 // count number of pixels read from file for debugging

	for {

		// Read a 16-bit unsigned integer, call it "a".
		value, err := readInt16(buf)
		if err != nil {
			log.Printf("error: %v (%d pixels read)", err, pxCount)
			return nil, err
		}

		pxCount++

		switch value {

		// If "a" is 0xFFFF (65535) then you must read another 16-bit
		// integer, and skip that number*2 pixels.
		case 0xffff:
			pixels, err := readInt16(buf)
			if err != nil {
				log.Printf("error: %v (%d pixels read)", err, pxCount)
				return nil, err
			}
			for i := 0; i < pixels*2; i++ {
				log.Printf("skip %d pixels", pixels*2)
				spriteData = append(spriteData, 0) // transparent background
			}

			pxCount += pixels * 2

		// If "a" is 0xFFFE (65534) then the next 16-bit integer*2
		// specifies the number of pixels you are going to draw.
		// Read that number of bytes from the file and draw their indexed
		// color.
		case 0xfffe:
			pixels, err := readInt16(buf)
			if err != nil {
				log.Printf("error: %v (%d pixels read)", err, pxCount)
				return nil, err
			}
			colors := make([]byte, pixels*2)
			_, err = buf.Read(colors)
			if err != nil {
				log.Printf("error: %v (%d pixels read)", err, pxCount)
				return nil, err
			}
			for i := 0; i < len(colors); i++ {
				log.Printf("draw %d pixels", len(colors))
				spriteData = append(spriteData, uint(colors[i]))
			}

			pxCount += len(colors)

		// If "a" is 0xFFFD (65533) then you are done. This is always the
		// last code in the file.
		case 0xfffd:
			log.Printf("done (%d pixels read)", pxCount)
			return &ImageResource{spriteData, 320}, nil

		}
	}

}
