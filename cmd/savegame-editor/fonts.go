package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/resources"
)

// fonts holds the fonts of the game converted to TrueType, indexed by name
// ("small" and "big"). Empty if the game files could not be loaded.
var fonts = map[string][]byte{}

// loadFonts converts the fonts of the game (GEODATA/SMALLSET.DAT and
// GEODATA/BIGLETS.DAT) in the game root folder to TrueType.
func loadFonts(root string) (map[string][]byte, error) {
	result := map[string][]byte{}
	for name, format := range map[string]resources.FontFormat{
		"small": resources.SmallFont,
		"big":   resources.BigFont,
	} {
		font, err := resources.LoadBitmapFont(root, format)
		if err != nil {
			return nil, fmt.Errorf("could not load %s: %w", format.Path, err)
		}
		result[name] = font.TTF("X-COM " + name)
	}
	return result, nil
}

// GET /api/fonts/{name}.ttf
func handleGetFont(w http.ResponseWriter, r *http.Request) {
	font, ok := fonts[mux.Vars(r)["name"]]
	if !ok {
		writeError(w, http.StatusNotFound, "font not found")
		return
	}
	w.Header().Set("Content-Type", "font/ttf")
	// the fonts are generated, make sure the browser picks up changes
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(font)
}
