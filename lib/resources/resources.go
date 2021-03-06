package resources

//go:generate go run github.com/redtoad/xcom-editor/lib/resources/generate

import (
	"errors"
	"image"
	"path"
)

// path to main palettes file
const palettePath = "GEODATA/PALETTES.DAT"

var (
	// ErrImageNotFound is thrown if path is not found in meta data map
	ErrImageNotFound = errors.New("image not found in meta data")
	// ErrNotImplemented is thrown if image type is not supported
	ErrNotImplemented = errors.New("not implemented yet")
)

// ResourceLoader will load resources from game directory.
type ResourceLoader struct {
	rootPath string
	palettes []*Palette
}

// LoadImage loads an image from pth with the palette as defined in resource list meta data. If pth
// is not found in meta data, an error is returned.
func (rs *ResourceLoader) LoadImage(pth string) (image.Image, error) {
	meta, ok := images[pth]
	if !ok {
		return nil, ErrImageNotFound
	}
	return rs.LoadImageWithPalette(pth, meta.PaletteNr)
}

// LoadImageWithPalette loads an image with a specific palette (rather than the one defined in
// meta data).
func (rs *ResourceLoader) LoadImageWithPalette(pth string, paletteNr int) (image.Image, error) {

	meta, ok := images[pth]
	if !ok {
		return nil, ErrImageNotFound
	}

	imgPath := path.Join(rs.rootPath, pth)
	tabPath := path.Join(rs.rootPath, meta.TabFile)
	palette := rs.palettes[paletteNr]

	ext := path.Ext(pth)
	switch ext {
	case ".PCK":

		if meta.TabFile != "" {
			tabs, err := LoadTAB(tabPath, meta.TabOffset)
			if err != nil {
				return nil, err
			}

			if len(tabs) == 0 {
				return nil, ErrNotEnoughTabs
			}

			collection, err := LoadImageCollectionFromPCK(imgPath, meta.Width, tabs)
			if err != nil {
				return nil, err
			}

			return collection.Gallery(10, meta.Height, palette)
		}

		img, err := LoadPCK(imgPath, meta.Width, 0)
		if err != nil {
			return nil, err
		}
		return img.Image(palette), nil

	case ".SPK":

		img, err := LoadSPK(imgPath)
		if err != nil {
			return nil, err
		}
		return img.Image(palette), nil

	case ".SCR":

		img, err := LoadSCR(imgPath, meta.Width)
		if err != nil {
			return nil, err
		}
		return img.Image(palette), nil

	}

	// unsupported file extension!
	return nil, ErrNotImplemented
}

// NewResourceLoader will return a new instance of ResourceLoader.
func NewResourceLoader(root string) (*ResourceLoader, error) {
	palettes, err := LoadPalettes(path.Join(root, palettePath))
	if err != nil {
		return nil, err
	}

	return &ResourceLoader{
		rootPath: root,
		palettes: palettes,
	}, nil
}

// ImageEntry is a resource file for images in formats
type ImageEntry struct {
	PaletteNr int
	Width     int
	Height    int
	TabFile   string
	TabOffset int
}
