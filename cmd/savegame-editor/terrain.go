package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/resources"
	"github.com/redtoad/xcom-editor/savegame"
)

// battlescapePalette is the palette used for the Battlescape (see
// PALETTES.DAT).
const battlescapePalette = 4

// Size of a terrain sprite in pixels.
const (
	terrainSpriteWidth  = 32
	terrainSpriteHeight = 40
)

// terrainSet holds a terrain file (TERRAIN/*.MCD) and its sprites.
type terrainSet struct {
	records []resources.MCDRecord
	// frames is the number of sprites in the sheet.
	frames int
	// sheet holds all sprites (PNG encoded) side by side.
	sheet []byte
}

// terrain holds all terrain sets of the game indexed by name (e.g.
// "JUNGLE"), nil if the game files could not be loaded.
var terrain map[string]*terrainSet

// loadTerrain loads all terrain sets (TERRAIN/*.MCD, *.PCK and *.TAB) from
// the game root folder.
func loadTerrain(root string) (map[string]*terrainSet, error) {
	palettes, err := resources.LoadPalettes(filepath.Join(root, "GEODATA", "PALETTES.DAT"))
	if err != nil {
		return nil, fmt.Errorf("could not load PALETTES.DAT: %w", err)
	}
	files, err := filepath.Glob(filepath.Join(root, "TERRAIN", "*.MCD"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no terrain files found in %s", filepath.Join(root, "TERRAIN"))
	}
	sets := make(map[string]*terrainSet)
	for _, file := range files {
		base := strings.TrimSuffix(file, filepath.Ext(file))
		set, err := loadTerrainSet(base, palettes[battlescapePalette])
		if err != nil {
			return nil, fmt.Errorf("could not load %s: %w", filepath.Base(base), err)
		}
		sets[strings.ToUpper(filepath.Base(base))] = set
	}
	return sets, nil
}

// loadTerrainSet loads a terrain set from base.MCD, base.PCK and base.TAB.
func loadTerrainSet(base string, palette *resources.Palette) (*terrainSet, error) {
	records, err := resources.LoadMCD(base + ".MCD")
	if err != nil {
		return nil, err
	}
	pck, err := os.ReadFile(base + ".PCK")
	if err != nil {
		return nil, err
	}
	offsets, err := resources.LoadTAB(base+".TAB", 2)
	if err != nil {
		return nil, err
	}
	sheet := image.NewRGBA(image.Rect(0, 0, max(1, len(offsets))*terrainSpriteWidth, terrainSpriteHeight))
	for i, offset := range offsets {
		sprite, err := resources.ParsePCK(pck, terrainSpriteWidth, offset)
		if err != nil {
			return nil, fmt.Errorf("could not load sprite %d: %w", i, err)
		}
		pos := image.Pt(i*terrainSpriteWidth, 0)
		img := sprite.Image(palette)
		draw.Draw(sheet, img.Bounds().Add(pos), img, image.Point{}, draw.Src)
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, sheet); err != nil {
		return nil, err
	}
	return &terrainSet{records: records, frames: len(offsets), sheet: buf.Bytes()}, nil
}

// terrainPartJS is the sprite of an MCD record.
type terrainPartJS struct {
	// Set is the index of the terrain set in battleTerrainJS.Sets.
	Set int `json:"set"`
	// Frame is the sprite in the sheet of the set.
	Frame int `json:"frame"`
	// Y is the vertical offset of the sprite in pixels.
	Y int `json:"y"`
}

// battleTerrainJS describes the sprites of the parts of a Battlescape map.
type battleTerrainJS struct {
	// Sets are the names of the terrain sets the map is built from, their
	// sprites are served at /api/terrain/{name}.png.
	Sets []string `json:"sets"`
	// Parts holds the sprite of each MCD index used in MAP.DAT.
	Parts []terrainPartJS `json:"parts"`
	// Tiles holds the MCD indices of the floor, west wall, north wall and
	// object of each tile (4 values per tile), in the order of
	// battleMapJS.Tiles. 0 means the part does not exist.
	Tiles []int `json:"tiles"`
	// Units holds the sprites of the units which can be drawn with the game
	// graphics, empty if the unit graphics are not available.
	Units []battleUnitSpritesJS `json:"units"`
}

// battleTerrain determines the terrain sets of a mission and the sprites of
// all MCD records.
func battleTerrain(b *savegame.Battle, sets map[string]*terrainSet) (*battleTerrainJS, error) {
	names := make([]string, 0, len(sets))
	for name := range sets {
		names = append(names, name)
	}
	sort.Strings(names)
	candidates := make([]savegame.TerrainSet, len(names))
	for i, name := range names {
		types := make([]uint8, len(sets[name].records))
		for j, r := range sets[name].records {
			types[j] = r.TileType
		}
		candidates[i] = savegame.TerrainSet{Name: name, TileTypes: types}
	}
	matched, err := b.MatchTerrainSets(candidates)
	if err != nil {
		return nil, err
	}
	result := &battleTerrainJS{Sets: matched, Parts: []terrainPartJS{}, Units: []battleUnitSpritesJS{}}
	if unitSheets != nil {
		result.Units = battleUnitSprites(b)
	}
	sx, sy, sz := b.Size()
	result.Tiles = make([]int, 0, sx*sy*sz*4)
	for z := 0; z < sz; z++ {
		for y := 0; y < sy; y++ {
			for x := 0; x < sx; x++ {
				t, _ := b.Tile(x, y, z)
				result.Tiles = append(result.Tiles, int(t.Ground), int(t.WestWall), int(t.NorthWall), int(t.Object))
			}
		}
	}
	for i, name := range matched {
		for _, r := range sets[name].records {
			result.Parts = append(result.Parts, terrainPartJS{Set: i, Frame: int(r.Frames[0]), Y: int(r.PLevel)})
		}
	}
	return result, nil
}

// GET /api/games/{slot}/battle/terrain
func handleGetBattleTerrain(w http.ResponseWriter, r *http.Request) {
	entry := getGameEntry(w, r)
	if entry == nil {
		return
	}
	if terrain == nil {
		writeError(w, http.StatusNotFound, "terrain graphics not available")
		return
	}
	entry.mu.RLock()
	defer entry.mu.RUnlock()
	if entry.sg.Battle == nil {
		writeError(w, http.StatusNotFound, "savegame was not saved during a mission")
		return
	}
	result, err := battleTerrain(entry.sg.Battle, terrain)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, result)
}

// GET /api/terrain/{name}.png
func handleGetTerrainSheet(w http.ResponseWriter, r *http.Request) {
	set, ok := terrain[strings.ToUpper(mux.Vars(r)["name"])]
	if !ok {
		writeError(w, http.StatusNotFound, "terrain set not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=86400")
	_, _ = w.Write(set.sheet)
}
