package main

import (
	"bytes"
	"fmt"
	"image/png"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/resources"
	"github.com/redtoad/xcom-editor/savegame"
)

// basePalette is the palette used for the base screen (see PALETTES.DAT).
const basePalette = 1

// Sprites in BASEBITS.PCK (the same layout is used by OpenXcom).
const (
	spriteGround        = 0  // ground of every grid cell
	spriteCorridorH     = 7  // corridor to the facility on the right
	spriteCorridorV     = 8  // corridor to the facility below
	spriteHangarShape   = 9  // hangar, 4 sprites (one per quarter)
	spriteFacilityStart = 17 // first facility interior (in order of geoscape.Facility)
	spriteCraftStart    = 33 // first craft (in order of geoscape.CraftType)
)

// tileSize is the size of a grid cell in pixels.
const tileSize = 32

// facilityShapes is the outline sprite of each facility type (1 = square,
// 2 = octagon, 3 = cross), in the order of geoscape.Facility.
var facilityShapes = []int{2, 1, 1, 1, 2, 1, 2, 1, 1, 3, 3, 3, 2, 2, 1, 2}

// baseBits holds the PNG encoded sprites of GEOGRAPH/BASEBITS.PCK, nil if
// the game files could not be loaded.
var baseBits [][]byte

// loadBaseBits loads the base screen sprites (GEOGRAPH/BASEBITS.PCK) from
// the game root folder.
func loadBaseBits(root string) ([][]byte, error) {
	palettes, err := resources.LoadPalettes(filepath.Join(root, "GEODATA", "PALETTES.DAT"))
	if err != nil {
		return nil, fmt.Errorf("could not load PALETTES.DAT: %w", err)
	}
	offsets, err := resources.LoadTAB(filepath.Join(root, "GEOGRAPH", "BASEBITS.TAB"), 2)
	if err != nil {
		return nil, fmt.Errorf("could not load BASEBITS.TAB: %w", err)
	}
	sprites, err := resources.LoadImageCollectionFromPCK(filepath.Join(root, "GEOGRAPH", "BASEBITS.PCK"), 32, offsets)
	if err != nil {
		return nil, fmt.Errorf("could not load BASEBITS.PCK: %w", err)
	}
	result := make([][]byte, len(sprites.Sprites))
	for i, sprite := range sprites.Sprites {
		var buf bytes.Buffer
		if err := png.Encode(&buf, sprite.Image(palettes[basePalette])); err != nil {
			return nil, fmt.Errorf("could not encode sprite %d: %w", i, err)
		}
		result[i] = buf.Bytes()
	}
	return result, nil
}

// baseSprite is a sprite drawn at a position (in pixels) of the base screen.
type baseSprite struct {
	Sprite int `json:"sprite"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// hangarCraft is a craft stationed at a base.
type hangarCraft struct {
	Type geoscape.CraftType
	// Out is true if the craft has left the base.
	Out bool
}

// baseLayers returns the sprites to draw a base like the base screen of the
// game: the ground, facility outlines, corridors between facilities, the
// facility interiors and the crafts in the hangars (in this order).
//
// Like the game, the crafts of the base are assigned to the built hangars in
// order; crafts that are out still occupy their hangar.
func baseLayers(tiles []savegame.BaseTile, crafts []hangarCraft) [][]baseSprite {
	const size = 6
	built := func(i int) bool {
		return tiles[i].Type != geoscape.Empty && tiles[i].DaysToCompletion == 0
	}
	isHangar := func(f geoscape.Facility) bool {
		return f >= geoscape.HangarTopLeft && f <= geoscape.HangarBottomRight
	}

	// empty layers must not be nil, so they are encoded as [] in JSON
	ground, shapes, corridors, interiors, craftSprites := []baseSprite{}, []baseSprite{}, []baseSprite{}, []baseSprite{}, []baseSprite{}
	for i, tile := range tiles {
		x, y := i%size*tileSize, i/size*tileSize
		ground = append(ground, baseSprite{spriteGround, x, y})
		f := tile.Type
		switch {
		case f == geoscape.Empty:
			continue
		case isHangar(f):
			quarter := int(f - geoscape.HangarTopLeft)
			shape := spriteHangarShape + quarter
			if !built(i) {
				shape += 4 // dotted outline while under construction
			}
			shapes = append(shapes, baseSprite{shape, x, y})
			if f == geoscape.HangarTopLeft && built(i) && len(crafts) > 0 {
				craft := crafts[0]
				crafts = crafts[1:]
				if !craft.Out {
					craftSprites = append(craftSprites, baseSprite{
						spriteCraftStart + int(craft.Type), x + tileSize/2 + 2, y + tileSize/2 - 4,
					})
				}
			}
		case int(f) < len(facilityShapes):
			shape := facilityShapes[f]
			if !built(i) {
				shape += 3 // dotted outline while under construction
			}
			shapes = append(shapes, baseSprite{shape, x, y})
			if built(i) {
				interiors = append(interiors, baseSprite{spriteFacilityStart + int(f), x, y})
			}
		}

		// corridors to built neighbours (not within a hangar)
		if !built(i) {
			continue
		}
		if right := i + 1; i%size < size-1 && built(right) &&
			!(f == geoscape.HangarTopLeft && tiles[right].Type == geoscape.HangarTopRight) &&
			!(f == geoscape.HangarBottomLeft && tiles[right].Type == geoscape.HangarBottomRight) {
			corridors = append(corridors, baseSprite{spriteCorridorH, x + tileSize/2, y})
		}
		if below := i + size; i/size < size-1 && built(below) &&
			!(f == geoscape.HangarTopLeft && tiles[below].Type == geoscape.HangarBottomLeft) &&
			!(f == geoscape.HangarTopRight && tiles[below].Type == geoscape.HangarBottomRight) {
			corridors = append(corridors, baseSprite{spriteCorridorV, x, y + tileSize/2})
		}
	}
	return [][]baseSprite{ground, shapes, corridors, interiors, craftSprites}
}

// miniFacility is a facility drawn on the mini base view of the base screen
// (the base selector at the top), in grid cells.
type miniFacility struct {
	X    int `json:"x"`
	Y    int `json:"y"`
	Size int `json:"size"`
	// Built is false while the facility is under construction.
	Built bool `json:"built"`
}

// miniFacilities returns the facilities of a base to draw it like the mini
// base view of the game: one block per facility, hangars are a single block
// of 2x2 cells.
func miniFacilities(tiles []savegame.BaseTile) []miniFacility {
	const size = 6
	result := []miniFacility{}
	for i, tile := range tiles {
		f := tile.Type
		if f == geoscape.Empty || (f > geoscape.HangarTopLeft && f <= geoscape.HangarBottomRight) {
			continue
		}
		s := 1
		if f == geoscape.HangarTopLeft {
			s = 2
		}
		result = append(result, miniFacility{i % size, i / size, s, tile.DaysToCompletion == 0})
	}
	return result
}

// GET /api/basebits/{sprite}.png
func handleGetBaseBit(w http.ResponseWriter, r *http.Request) {
	idx, err := strconv.Atoi(mux.Vars(r)["sprite"])
	if err != nil || idx < 0 || idx >= len(baseBits) {
		writeError(w, http.StatusNotFound, "sprite not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=86400")
	_, _ = w.Write(baseBits[idx])
}
