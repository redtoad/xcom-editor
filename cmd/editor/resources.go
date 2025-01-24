package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/redtoad/xcom-editor/resources"
)

func ServeImage(w http.ResponseWriter, r *http.Request) {

	pth := r.URL.Path[10:]
	basename := path.Base(pth)
	var img image.Image
	var err error
	paletteNr := r.URL.Query().Get("palette")

	if paletteNr != "" {
		val, _ := strconv.Atoi(paletteNr)
		img, err = loader.LoadImageWithPalette(pth, val)
	} else {
		img, err = loader.LoadImage(pth)
	}
	if err != nil {
		if err == resources.ErrImageNotFound {
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
	if err = png.Encode(buf, img); err != nil {
		log.Printf("Error: Could not load image from %s: %s\n", pth, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", fmt.Sprintf("filename=\"%s\"", basename))
	if _, err = w.Write(buf.Bytes()); err != nil {
		log.Printf("Error: Could not write response: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
