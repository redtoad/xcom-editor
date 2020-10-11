package resources

//go:generate go run github.com/redtoad/xcom-editor/resources/generate

import (
	"errors"
	"image"
	"path"
)

// path to main palettes file
const palettePath = "GEODATA/PALETTES.DAT"

var (
	ErrImageNotFound  = errors.New("image not found in meta data")
	ErrNotImplemented = errors.New("not implemented yet")
)

type ResourceLoader struct {
	rootPath string
	palettes []*Palette
}

func (rs *ResourceLoader) LoadImage(pth string) (image.Image, error) {

	meta, ok := images[pth]
	if !ok {
		return nil, ErrImageNotFound
	}

	imgPath := path.Join(rs.rootPath, pth)
	tabPath := path.Join(rs.rootPath, meta.TabFile)
	palette := rs.palettes[meta.PaletteNr]

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
