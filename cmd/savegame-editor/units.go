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
var unitSheetFiles = []string{"XCOM_0", "XCOM_1", "XCOM_2", "TANKS", "SECTOID", "SNAKEMAN", "MUTON", "CHRYS", "FLOATER", "X_REAP", "CYBER", "CIVM", "CIVF"}

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
	// shapeTank units are large units with a turret (X-COM tanks).
	shapeTank
	// shapeFloater units are humanoid units without legs (Floaters), their
	// torso reaches down to the ground.
	shapeFloater
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

// Arm sprites of Floaters (UNITS/FLOATER.PCK), one per direction. Unlike the
// other humanoid units, the right arm comes first.
const (
	spriteFloaterRightArm = 0
	spriteFloaterLeftArm  = 8
)

// Sprites of tanks (UNITS/TANKS.PCK). The hull has one sprite per tile and
// direction like other large units.
const (
	// spriteHovertank is the first hull sprite of hovertanks.
	spriteHovertank = 32
	// spriteTurrets is the first turret sprite, 8 per turret (Unit.TurretImage).
	spriteTurrets = 64
	// spritePropulsion is the first sprite of the propulsion drawn below the
	// right, left and lower tile of hovertanks (8 animation frames each).
	spritePropulsion = 104
	// turretOffsetY moves the turret up (in pixels); it is drawn on the lower
	// tile.
	turretOffsetY = -4
)

// hovertankTurretOffset moves the turret of hovertanks for each direction
// (in pixels).
var hovertankTurretOffset = [8][2]int{{-2, -1}, {-7, -3}, {-5, -4}, {0, -5}, {5, -4}, {7, -3}, {2, -1}, {0, -1}}

// largeUnitOffsets are the positions of the sprites of the 4 tiles of large
// units relative to the upper left tile: upper left, upper right, lower left,
// lower right.
var largeUnitOffsets = [4][2]int{{0, 0}, {16, 8}, {-16, 8}, {0, 16}}

// unitBodies are the unit types (UNITREF.DAT Type) which are drawn with
// their sprites, other units are drawn as markers.
var unitBodies = map[uint8]unitBody{
	0:  {Sheet: "XCOM_0", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // no armour
	1:  {Sheet: "XCOM_1", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // personal armour
	2:  {Sheet: "XCOM_2", Shape: shapeHumanoid, Torso: 32, FemaleTorso: 267, Dead: 266}, // power or flying suit
	3:  {Sheet: "TANKS", Shape: shapeTank, Dead: -1},                                    // flying suits use type 2
	4:  {Sheet: "SECTOID", Shape: shapeHumanoid, Torso: 32, Dead: 266},
	5:  {Sheet: "SNAKEMAN", Shape: shapeHumanoid, Torso: 24, Dead: 98},
	7:  {Sheet: "MUTON", Shape: shapeHumanoid, Torso: 32, Dead: 266},
	8:  {Sheet: "FLOATER", Shape: shapeFloater, Torso: 16, Dead: 66},
	11: {Sheet: "CHRYS", Shape: shapeHumanoid, Torso: 24, Dead: 226},
	12: {Sheet: "X_REAP", Shape: shapeLarge, Dead: -1},
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

// movementFly is Unit.MovementType of flying units.
const movementFly = 2

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
		// one sprite per tile
		result := make([]unitSpriteJS, 4)
		for part, off := range largeUnitOffsets {
			result[part] = sprite(part*8+dir, off[0], off[1])
		}
		return result
	case shapeTank:
		return tankSprites(u, sprite)
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
	if body.Shape == shapeFloater {
		// no legs and the arms swapped, Floaters cannot kneel
		parts["leftArm"] = sprite(spriteFloaterLeftArm+dir, 0, 0)
		parts["rightArm"] = sprite(spriteFloaterRightArm+dir, 0, 0)
		parts["torso"] = sprite(torso+dir, 0, 0)
		delete(parts, "legs")
	}
	result := make([]unitSpriteJS, 0, 4)
	for _, name := range humanoidOrder[dir] {
		if part, ok := parts[name]; ok {
			result = append(result, part)
		}
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

// tankSprites returns the sprites of a tank: the hull (one sprite per tile,
// with the propulsion below for hovertanks) and the turret on top.
func tankSprites(u savegame.MapUnit, sprite func(frame, dx, dy int) unitSpriteJS) []unitSpriteJS {
	dir := int(u.Unit.Facing % 8)
	hover := u.Unit.MovementType == movementFly
	hull := 0
	if hover {
		hull = spriteHovertank
	}
	result := make([]unitSpriteJS, 0, 8)
	for part, off := range largeUnitOffsets {
		if hover && part > 0 {
			result = append(result, sprite(spritePropulsion+(part-1)*8, off[0], off[1]))
		}
		result = append(result, sprite(hull+part*8+dir, off[0], off[1]))
	}
	dx, dy := largeUnitOffsets[3][0], largeUnitOffsets[3][1]+turretOffsetY
	if hover {
		dx += hovertankTurretOffset[dir][0]
		dy += hovertankTurretOffset[dir][1]
	}
	return append(result, sprite(spriteTurrets+int(u.Unit.TurretImage)*8+dir, dx, dy))
}
