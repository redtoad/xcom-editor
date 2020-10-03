package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/resources"
)

// Root path of game (containing all images and save games)
var root = "."

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
					log.Printf("Warning: sprite %dx%d is bigger than specified in meta data (%dx%d)!\n",
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

	img, err := LoadImage(root, pth)
	if err != nil {
		if err == ErrImageNotFound {
			log.Printf("Error: File %s not found!\n", pth)
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}
		log.Printf("Error: Could not load %s: %s\n", pth, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Serving file %s\n", pth)

	buf := new(bytes.Buffer)
	_ = png.Encode(buf, img)

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", basename))
	_, err = w.Write(buf.Bytes())

}

func main() {
	flag.StringVar(&root, "root", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	log.Printf("Starting server...\n")
	log.Printf("Game root: %s\n", root)
	log.Println("Try opening http://localhost:8080/resource/UNITS/ZOMBIE.PCK")

	r := mux.NewRouter()
	r.PathPrefix("/resource").HandlerFunc(ServeImage)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
