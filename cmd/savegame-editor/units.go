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
	"strings"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/redtoad/xcom-editor/resources"
	"github.com/redtoad/xcom-editor/savegame"
)

// Size of a unit sprite in pixels (the same as terrain sprites).
const (
	unitSpriteWidth  = 32
	unitSpriteHeight = 40
)

// unitSheets holds the sprites of the units (UNITS/*.PCK, PNG encoded side
// by side) indexed by name (e.g. "XCOM_1"), nil if the game files could not
// be loaded.
var unitSheets map[string][]byte

// unitSheetFiles are the sprite files of the units which are drawn (see
// unitBodies).
var unitSheetFiles = []string{"XCOM_0", "XCOM_1", "XCOM_2", "SECTOID", "SNAKEMAN", "MUTON", "CHRYS", "CYBER", "CIVM", "CIVF"}

// loadUnitSheets loads the unit sprites (UNITS/*.PCK and *.TAB) from the
// game root folder.
func loadUnitSheets(root string) (map[string][]byte, error) {
	palettes, err := resources.LoadPalettes(filepath.Join(root, "GEODATA", "PALETTES.DAT"))
	if err != nil {
		return nil, fmt.Errorf("could not load PALETTES.DAT: %w", err)
	}
	palette := resources.BattlescapePalette(palettes)
	sheets := make(map[string][]byte)
	for _, name := range unitSheetFiles {
		base := filepath.Join(root, "UNITS", name)
		pck, err := os.ReadFile(base + ".PCK")
		if err != nil {
			return nil, err
		}
		offsets, err := resources.LoadTAB(base+".TAB", 2)
		if err != nil {
			return nil, err
		}
		sheet := image.NewRGBA(image.Rect(0, 0, max(1, len(offsets))*unitSpriteWidth, unitSpriteHeight))
		for i, offset := range offsets {
			sprite, err := resources.ParsePCK(pck, unitSpriteWidth, offset)
			if err != nil {
				return nil, fmt.Errorf("could not load sprite %d of %s: %w", i, name, err)
			}
			img := sprite.Image(palette)
			draw.Draw(sheet, img.Bounds().Add(image.Pt(i*unitSpriteWidth, 0)), img, image.Point{}, draw.Src)
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, sheet); err != nil {
			return nil, err
		}
		sheets[name] = buf.Bytes()
	}
	return sheets, nil
}

// unitShape is how the sprites of a unit are put together.
type unitShape int

const (
	// shapeHumanoid units consist of legs, torso and two arms.
	shapeHumanoid unitShape = iota
	// shapeSingle units have one sprite per direction.
	shapeSingle
	// shapeLarge units occupy 2x2 tiles, one sprite per tile.
	shapeLarge
)

// unitBody describes the sprites of a unit type (UNITREF.DAT Type).
type unitBody struct {
	Sheet string
	Shape unitShape
	// Torso is the first torso sprite (one per direction) of humanoid units.
	Torso int
	// FemaleTorso is the first torso sprite of female soldiers, 0 if there
	// is none.
	FemaleTorso int
	// Dead is the sprite of the dead or unconscious unit, -1 if there is none.
	Dead int
}

// Sprites of humanoid units, one per direction.
const (
	spriteLeftArm   = 0
	spriteRightArm  = 8
	spriteLegs      = 16
	spriteLegsKneel = 24
	// kneelOffset moves the upper body of kneeling units down (in pixels).
	kneelOffset = 4
)

// unitBodies are the unit types (UNITREF.DAT Type) which are drawn with
// their sprites, other units are drawn as markers.
var unitBodies = map[uint8]unitBody{
	0:  {Sheet: "XCOM_0", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // no armour
	1:  {Sheet: "XCOM_1", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // personal armour
	2:  {Sheet: "XCOM_2", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // power suit
	3:  {Sheet: "XCOM_2", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // flying suit
	4:  {Sheet: "SECTOID", Shape: shapeHumanoid, Torso: 32, Dead: 266},
	5:  {Sheet: "SNAKEMAN", Shape: shapeHumanoid, Torso: 24, Dead: 98},
	7:  {Sheet: "MUTON", Shape: shapeHumanoid, Torso: 32, Dead: 266},
	11: {Sheet: "CHRYS", Shape: shapeHumanoid, Torso: 24, Dead: 226},
	14: {Sheet: "CYBER", Shape: shapeLarge, Dead: -1},
	16: {Sheet: "CIVM", Shape: shapeSingle, Dead: 74},
	17: {Sheet: "CIVF", Shape: shapeSingle, Dead: 74},
}

// humanoidOrder is the order in which the parts of a humanoid unit are
// drawn for each direction, so that the parts in front are drawn last.
var humanoidOrder = [8][4]string{
	{"leftArm", "legs", "torso", "rightArm"},
	{"leftArm", "legs", "torso", "rightArm"},
	{"leftArm", "legs", "torso", "rightArm"},
	{"legs", "torso", "leftArm", "rightArm"},
	{"legs", "rightArm", "torso", "leftArm"},
	{"rightArm", "legs", "torso", "leftArm"},
	{"rightArm", "legs", "torso", "leftArm"},
	{"leftArm", "rightArm", "legs", "torso"},
}

// Flags in Unit.Flags.
const unitKneeling = 1 << 6

// unitSpriteJS is a sprite of a unit.
type unitSpriteJS struct {
	// Sheet is the sprite file, served at /api/unitsprites/{sheet}.png.
	Sheet string `json:"sheet"`
	// Frame is the sprite in the sheet.
	Frame int `json:"frame"`
	// DX and DY is the position of the sprite relative to the sprite of the
	// upper left tile of the unit, in pixels.
	DX int `json:"dx"`
	DY int `json:"dy"`
}

// unitSprites returns the sprites of a unit in the order they are drawn, nil
// if the unit type has no sprites (it is drawn as a marker).
func unitSprites(u savegame.MapUnit) []unitSpriteJS {
	body, ok := unitBodies[u.Unit.Type]
	if !ok {
		return nil
	}
	dir := int(u.Unit.Facing % 8)
	sprite := func(frame, dx, dy int) unitSpriteJS {
		return unitSpriteJS{Sheet: body.Sheet, Frame: frame, DX: dx, DY: dy}
	}
	if !u.Alive() || u.Unconscious() {
		if body.Dead < 0 {
			return nil
		}
		return []unitSpriteJS{sprite(body.Dead, 0, 0)}
	}
	switch body.Shape {
	case shapeSingle:
		return []unitSpriteJS{sprite(dir, 0, 0)}
	case shapeLarge:
		// one sprite per tile: upper left, upper right, lower left, lower right
		offsets := [4][2]int{{0, 0}, {16, 8}, {-16, 8}, {0, 16}}
		result := make([]unitSpriteJS, 4)
		for part, off := range offsets {
			result[part] = sprite(part*8+dir, off[0], off[1])
		}
		return result
	}
	torso := body.Torso
	if body.FemaleTorso > 0 && u.Unit.Gender == 1 {
		torso = body.FemaleTorso
	}
	legs, upper := spriteLegs, 0
	if u.Unit.Flags&unitKneeling != 0 {
		legs, upper = spriteLegsKneel, kneelOffset
	}
	parts := map[string]unitSpriteJS{
		"leftArm":  sprite(spriteLeftArm+dir, 0, upper),
		"rightArm": sprite(spriteRightArm+dir, 0, upper),
		"legs":     sprite(legs+dir, 0, 0),
		"torso":    sprite(torso+dir, 0, upper),
	}
	result := make([]unitSpriteJS, 0, 4)
	for _, name := range humanoidOrder[dir] {
		result = append(result, parts[name])
	}
	return result
}

// battleUnitSpritesJS holds the sprites of a unit on the map.
type battleUnitSpritesJS struct {
	// Index is the slot of the unit (battleUnitJS.Index).
	Index   int            `json:"index"`
	Sprites []unitSpriteJS `json:"sprites"`
}

// battleUnitSprites returns the sprites of all units on the map which can be
// drawn with sprites.
func battleUnitSprites(b *savegame.Battle) []battleUnitSpritesJS {
	result := []battleUnitSpritesJS{}
	for _, u := range b.MapUnits() {
		if u.Owner > battlescape.OwnerCivilian {
			continue // left behind on a previous stage
		}
		if sprites := unitSprites(u); sprites != nil {
			result = append(result, battleUnitSpritesJS{Index: u.Index, Sprites: sprites})
		}
	}
	return result
}

// GET /api/unitsprites/{name}.png
func handleGetUnitSheet(w http.ResponseWriter, r *http.Request) {
	sheet, ok := unitSheets[strings.ToUpper(mux.Vars(r)["name"])]
	if !ok {
		writeError(w, http.StatusNotFound, "unit sprites not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	// the sprites depend on the game folder, which can be switched
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(sheet)
}
