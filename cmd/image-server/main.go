package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"path"

	"github.com/redtoad/xcom-editor/resources"
)

type ImageEntry struct {
	PaletteNr int
	Width     int
	Height    int
	TabFile   string
	TabOffset int
}

var images = map[string]ImageEntry{

	// ...
	"GEOGRAPH/BASEBITS.PCK": {1, 32, 40, "GEOGRAPH/BASEBITS.TAB", 2},
	"GEOGRAPH/INTICON.PCK":  {1, 32, 40, "GEOGRAPH/INTICON.TAB", 2},
	// ...
	"UNITS/BIGOBS.PCK": {4, 32, 48, "UNITS/BIGOBS.TAB", 2},
	"UNITS/XCOM_0.PCK": {1, 32, 40, "UNITS/XCOM_0.TAB", 2},
	"UNITS/XCOM_1.PCK": {1, 32, 40, "UNITS/XCOM_1.TAB", 2},
	"UNITS/XCOM_2.PCK": {1, 32, 40, "UNITS/XCOM_2.TAB", 2},
	"UNITS/X_REAP.PCK": {1, 32, 40, "UNITS/X_REAP.TAB", 2},
	"UNITS/X_ROB.PCK":  {1, 32, 40, "UNITS/X_ROB.TAB", 2},
	// ...
	"UNITS/ZOMBIE.PCK": {1, 32, 40, "UNITS/ZOMBIE.TAB", 2},

	// INTERWIN.DAT - 10 images 160px wide no compression
	// LANG1.DAT - GeoScape control panel overlay in German. Direct palette indexes, no compression, 64x154.
	// LANG2.DAT - GeoScape control panel overlay in French. Direct palette indexes, no compression, 64x154.

}

// root path of X-Com game
var root = os.Args[1]

func init() {

	// add all 56 BIGOB
	for i := 0; i <= 56; i++ {
		images[fmt.Sprintf("UFOGRAPH/BIGOB_%02d.PCK", i)] = ImageEntry{1, 32, 40, "", 4}
	}

	// add soldier images
	for armour := 0; armour <= 1; armour++ {
		for _, sex := range []string{"F", "M"} {
			for race := 0; race <= 3; race++ {
				images[fmt.Sprintf("UFOGRAPH/MAN_%d%s%d.SPK", armour, sex, race)] = ImageEntry{4, 32, 40, "", 4}
			}
		}
	}
	// soldiers in armoured suits
	images["UFOGRAPH/MAN_2.SPK"] = ImageEntry{4, 32, 40, "", 4}
	images["UFOGRAPH/MAN_3.SPK"] = ImageEntry{4, 32, 40, "", 4}
}

var PalettePath = "GEODATA/PALETTES.DAT"

var (
	ErrImageNotFound    = errors.New("image not found in meta data")
	ErrNotImplemented   = errors.New("not implemented yet")
	ErrNotEnoughSprites = errors.New("could not load sprites from file")
	ErrNotEnoughTabs    = errors.New("could not load offsets from tab file")
)

// min returns the larger of a and b
func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func LoadImage(root string, pth string) (image.Image, error) {

	meta, ok := images[pth]
	if !ok {
		return nil, ErrImageNotFound
	}

	palettes, err := resources.LoadPalettes(path.Join(root, PalettePath))
	if err != nil {
		return nil, err
	}

	ext := path.Ext(pth)
	switch ext {
	case ".PCK":

		if meta.TabFile != "" {
			tabs, err := resources.LoadTAB(path.Join(root, meta.TabFile), meta.TabOffset)
			if err != nil {
				return nil, err
			}

			if len(tabs) == 0 {
				return nil, ErrNotEnoughTabs
			}

			var sprites []*resources.ImageResource
			for _, tab := range tabs {
				sprite, err := resources.LoadPCK(path.Join(root, pth), meta.Width, tab)
				if err != nil {
					return nil, err
				}
				sprites = append(sprites, sprite)

				if sprite.Height() > meta.Height {
					fmt.Printf("Warning: sprite %dx%d is bigger than specified in meta data (%dx%d)!\n",
						sprite.Width(), sprite.Height(), meta.Width, meta.Height)
				}

			}

			if len(sprites) == 0 {
				return nil, ErrNotEnoughSprites
			}

			// calculate grid size for collection
			width := min(len(sprites), 10)
			height := (len(sprites) / width) + 1

			// create new image with black background
			collection := image.NewRGBA(image.Rect(
				0, 0,
				width*meta.Width, height*meta.Height))

			// draw each image onto collection
			for no, sprite := range sprites {

				dstX := (no % width) * meta.Width
				dstY := (no / width) * meta.Height
				dstR := image.Rectangle{
					Min: image.Point{dstX, dstY},
					Max: image.Point{dstX + sprite.Width(), dstY + sprite.Height()},
				}

				img := sprite.Image(palettes[meta.PaletteNr])
				draw.Draw(collection, dstR,
					img, img.Bounds().Min,
					draw.Src)

			}

			return collection, nil
		}

		img, err := resources.LoadPCK(path.Join(root, pth), meta.Width, 0)
		if err != nil {
			return nil, err
		}
		return img.Image(palettes[meta.PaletteNr]), nil

	case ".SPK":

		img, err := resources.LoadSPK(path.Join(root, pth))
		if err != nil {
			return nil, err
		}
		return img.Image(palettes[meta.PaletteNr]), nil

	}

	// currently only PCKs and SPK are supported
	return nil, ErrNotImplemented

}

func ServeImage(w http.ResponseWriter, r *http.Request) {

	pth := r.RequestURI[10:]
	basename := path.Base(pth)

	println(pth)

	img, err := LoadImage(root, pth)
	if err != nil {
		if err == ErrImageNotFound {
			http.Error(w, "image not found", http.StatusNotFound)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", basename))
	_, err = w.Write(buf.Bytes())

}

func main() {
	fmt.Printf("Game root: %s\n", root)
	fmt.Printf("Starting server...\n")
	fmt.Println("Try opening http://localhost:8080/resource/UNITS/ZOMBIE.PCK")
	http.HandleFunc("/", ServeImage)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
