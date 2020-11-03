package resources

import (
	"io/ioutil"
)

// LoadSCR load SCR image from path.
//
// These ones are easy, every byte is an uncompressed index into the game's palette.
// Typically SCR files are used for 320x200 background images (and often stored in greyscale
// so they can be re-colored on the fly). DAT files that contain images use the same format,
// though the line width tends to vary depending on the specific use to which they are to
// be put.
func LoadSCR(path string, width int) (*ImageResource, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	data := make([]uint, len(buf))
	for i, bt := range buf {
		data[i] = uint(bt)
	}
	return &ImageResource{pixels: data, width: width}, nil
}
