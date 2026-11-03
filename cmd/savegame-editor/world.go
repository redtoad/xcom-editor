package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geodata"
	"github.com/redtoad/xcom-editor/resources"
	"github.com/redtoad/xcom-editor/savegame"
)

// geoscapePalette is the palette used for the Geoscape (see PALETTES.DAT).
const geoscapePalette = 0

// polarSeas is the terrain type of polar seas with icebergs, its most
// common color is used for the ocean.
const polarSeas = 12

// worldPolygon is a polygon of the globe in GPS coordinates.
type worldPolygon struct {
	Terrain int          `json:"terrain"`
	Points  [][2]float32 `json:"points"` // [lon, lat]
}

// worldData holds the globe terrain and textures from the game files.
type worldData struct {
	Polygons []worldPolygon
	Ocean    string
	// Textures holds the PNG encoded textures, indexed by texture set and
	// terrain type.
	Textures [geodata.NumTextureSets][geodata.NumTerrains][]byte
}

// world is nil if the game files could not be loaded.
var world *worldData

// loadWorld loads the globe terrain (GEODATA/WORLD.DAT) and textures
// (GEOGRAPH/TEXTURE.DAT) from the game root folder.
func loadWorld(root string) (*worldData, error) {
	var worldFile geodata.WorldFile
	if err := internal.LoadDATFile(filepath.Join(root, "GEODATA", "WORLD.DAT"), &worldFile); err != nil {
		return nil, err
	}
	var textures geodata.TextureFile
	if err := internal.LoadDATFile(filepath.Join(root, "GEOGRAPH", "TEXTURE.DAT"), &textures); err != nil {
		return nil, err
	}
	palettes, err := resources.LoadPalettes(filepath.Join(root, "GEODATA", "PALETTES.DAT"))
	if err != nil {
		return nil, fmt.Errorf("could not load PALETTES.DAT: %w", err)
	}
	palette := palettes[geoscapePalette].Palette()

	data := &worldData{}
	for _, p := range worldFile.Polygons {
		xs := []int{p.X0, p.X1, p.X2, p.X3}
		ys := []int{p.Y0, p.Y1, p.Y2, p.Y3}
		if p.Type() == geodata.Triangle {
			xs, ys = xs[:3], ys[:3]
		}
		poly := worldPolygon{Terrain: p.Terrain}
		for i := range xs {
			c := savegame.NewCoord(xs[i], ys[i])
			poly.Points = append(poly.Points, [2]float32{c.Lon, c.Lat})
		}
		data.Polygons = append(data.Polygons, poly)
	}

	for set := range data.Textures {
		for terrain := range data.Textures[set] {
			var buf bytes.Buffer
			if err := png.Encode(&buf, textures.Image(set, terrain, *palette)); err != nil {
				return nil, fmt.Errorf("could not encode texture: %w", err)
			}
			data.Textures[set][terrain] = buf.Bytes()
		}
	}

	ocean := textures.Image(0, polarSeas, *palette)
	data.Ocean = hexColor(ocean.Palette[mostCommon(ocean.Pix)])
	return data, nil
}

// mostCommon returns the most common value.
func mostCommon(values []uint8) uint8 {
	var counts [256]int
	var best uint8
	for _, v := range values {
		counts[v]++
		if counts[v] > counts[best] {
			best = v
		}
	}
	return best
}

func hexColor(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

// GET /api/world
func handleGetWorld(w http.ResponseWriter, r *http.Request) {
	if world == nil {
		writeJSON(w, map[string]any{"available": false})
		return
	}
	writeJSON(w, map[string]any{
		"available":   true,
		"ocean":       world.Ocean,
		"polygons":    world.Polygons,
		"textureSets": geodata.NumTextureSets,
		"terrains":    geodata.NumTerrains,
	})
}

// GET /api/world/textures/{set}/{terrain}.png
func handleGetTexture(w http.ResponseWriter, r *http.Request) {
	if world == nil {
		writeError(w, http.StatusNotFound, "game files not found")
		return
	}
	set, err1 := strconv.Atoi(mux.Vars(r)["set"])
	terrain, err2 := strconv.Atoi(mux.Vars(r)["terrain"])
	if err1 != nil || err2 != nil ||
		set < 0 || set >= geodata.NumTextureSets || terrain < 0 || terrain >= geodata.NumTerrains {
		writeError(w, http.StatusNotFound, "texture not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=86400")
	_, _ = w.Write(world.Textures[set][terrain])
}
