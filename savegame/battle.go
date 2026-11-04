package savegame

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

// Battle holds the Battlescape data of a savegame that was saved during a
// mission.
type Battle struct {
	// Globals holds the map size, current turn and mission score (WGLOB.DAT).
	Globals battlescape.WglobFile
	// Lighting holds the ambient light of the map (BGLOB.DAT).
	Lighting battlescape.BglobFile
	// Geodata describes how the map was built from map modules (GEODATA.DAT).
	Geodata battlescape.GeodataFile
	// Mission holds the mission data and score (MISDATA.DAT).
	Mission battlescape.MisdataFile
	// Map holds the tiles of the map (MAP.DAT).
	Map battlescape.MapFile
	// Routes holds the route nodes for spawning and path finding (ROUTES.DAT).
	Routes battlescape.RoutesFile
	// FogOfWar holds which tiles X-COM has seen (SEEMAP.DAT).
	FogOfWar battlescape.SeemapFile
	// Smoke holds which tiles have smoke or fire (SMOKBIT.DAT).
	Smoke battlescape.SmokbitFile
	// SmokeClouds holds the smoke and fire on the map (SMOKREF.DAT).
	SmokeClouds battlescape.SmokrefFile
	// TerrainLight holds the fixed and ambient light of each tile (TERMP.DAT).
	TerrainLight battlescape.LightMapFile
	// SourceLight holds the light of mobile light sources (SOURCEMP.DAT).
	SourceLight battlescape.LightMapFile
	// Objects holds all objects that can be picked up (OBPOS.DAT).
	Objects battlescape.ObposFile
	// UnitPositions holds the positions and owners of all units (UNITPOS.DAT).
	UnitPositions battlescape.UnitposFile
	// Units holds the stats of all units (UNITREF.DAT).
	Units battlescape.UnitrefFile
}

// battleFiles maps the Battlescape data files to the fields they are loaded into.
func (b *Battle) battleFiles() map[string]any {
	return map[string]any{
		"WGLOB.DAT":    &b.Globals,
		"BGLOB.DAT":    &b.Lighting,
		"GEODATA.DAT":  &b.Geodata,
		"MISDATA.DAT":  &b.Mission,
		"MAP.DAT":      &b.Map,
		"ROUTES.DAT":   &b.Routes,
		"SEEMAP.DAT":   &b.FogOfWar,
		"SMOKBIT.DAT":  &b.Smoke,
		"SMOKREF.DAT":  &b.SmokeClouds,
		"TERMP.DAT":    &b.TerrainLight,
		"SOURCEMP.DAT": &b.SourceLight,
		"OBPOS.DAT":    &b.Objects,
		"UNITPOS.DAT":  &b.UnitPositions,
		"UNITREF.DAT":  &b.Units,
	}
}

// loadBattle loads the Battlescape data if the savegame was saved during a
// mission (i.e. it contains WGLOB.DAT). Otherwise game.Battle is nil.
func (game *Savegame) loadBattle() error {
	game.Battle = nil
	if _, err := os.Stat(path.Join(game.Path, "WGLOB.DAT")); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	battle := &Battle{}
	for name, obj := range battle.battleFiles() {
		if err := internal.LoadDATFile(path.Join(game.Path, name), obj); err != nil {
			return fmt.Errorf("could not load %s: %w", name, err)
		}
	}
	game.Battle = battle
	return nil
}

func (game *Savegame) saveBattle() error {
	if game.Battle == nil {
		return nil
	}
	for name, obj := range game.Battle.battleFiles() {
		if err := internal.SaveDATFile(path.Join(game.Path, name), obj); err != nil {
			return fmt.Errorf("could not save %s: %w", name, err)
		}
	}
	return nil
}

// MissionStatus summarizes an ongoing mission.
type MissionStatus struct {
	// Turn is the current turn (starting with 1).
	Turn int
	// Craft is the X-COM craft that brought the soldiers (nil for base
	// defence missions).
	Craft *Craft
	// XcomAlive and XcomTotal count the X-COM units (soldiers and tanks).
	XcomAlive, XcomTotal int
	// AliensAlive and AliensTotal count the alien units.
	AliensAlive, AliensTotal int
	// Score is the current mission score.
	Score int
}

// MissionStatus returns a summary of the ongoing mission, or nil if the
// savegame was not saved during a mission.
func (game *Savegame) MissionStatus() *MissionStatus {
	b := game.Battle
	if b == nil {
		return nil
	}
	status := &MissionStatus{
		Turn:  int(b.Globals.Turn),
		Score: int(b.Globals.AlienKillScore) + int(b.Globals.XcomLossScore),
	}
	// MISDATA.DAT refers to the craft by its location (LOC.DAT)
	if idx := int(b.Mission.Craft); idx < len(game.locationFile.Objects) {
		if loc := game.locationFile.Objects[idx]; loc.Type == geoscape.XCOMShip {
			status.Craft = game.Craft(int(loc.TableReference))
		}
	}
	for _, u := range b.MapUnits() {
		switch u.Owner {
		case battlescape.OwnerXcom:
			status.XcomTotal++
			if u.Alive() {
				status.XcomAlive++
			}
		case battlescape.OwnerAlien:
			status.AliensTotal++
			if u.Alive() {
				status.AliensAlive++
			}
		}
	}
	return status
}

// MapUnit is a unit on the Battlescape map.
type MapUnit struct {
	// Index is the slot in UNITPOS.DAT and UNITREF.DAT.
	Index int
	// X, Y and Z is the position of the unit. Z is the map level in the
	// order of MAP.DAT, i.e. 0 is the top level and SizeZ-1 the ground level.
	// Large units are placed at their upper left tile.
	X, Y, Z int
	// Owner is battlescape.OwnerXcom, OwnerAlien, OwnerCivilian or 3 for
	// X-COM units left behind on a previous stage.
	Owner uint8
	// Large is set for units occupying 2x2 tiles.
	Large bool
	// Visible is set if the unit is visible to X-COM.
	Visible bool
	// Unit holds the stats of the unit (UNITREF.DAT).
	Unit *battlescape.Unit
	// active is set if the unit takes part in the battle (not dead).
	active bool
}

// Alive reports whether the unit is still alive.
func (u MapUnit) Alive() bool {
	return u.active && u.Unit.Health > 0
}

// Unconscious reports whether the unit is alive, but stunned.
func (u MapUnit) Unconscious() bool {
	return u.Alive() && u.Unit.Stun >= u.Unit.Health
}

// MapUnits returns all units on the map. Unused slots and the additional
// slots of large units (which occupy 4 slots) are skipped.
func (b *Battle) MapUnits() []MapUnit {
	var units []MapUnit
	for i, pos := range b.UnitPositions.Units {
		unit := &b.Units.Units[i]
		if pos.Y == 0xFF || unit.Type == 0xFF {
			continue // unused slot
		}
		large := pos.Flags&battlescape.UnitLarge != 0
		// large units occupy 4 slots, only use the upper left one
		if large && pos.Flags&(battlescape.UnitRightSection|battlescape.UnitLowerSection) != 0 {
			continue
		}
		units = append(units, MapUnit{
			Index:   i,
			X:       int(pos.X),
			Y:       int(pos.Y),
			Z:       int(pos.Z),
			Owner:   pos.Owner,
			Large:   large,
			Visible: pos.Flags&battlescape.UnitVisible != 0,
			Unit:    unit,
			active:  pos.Flags&battlescape.UnitActive != 0,
		})
	}
	return units
}

// Size returns the dimensions of the map.
func (b *Battle) Size() (x, y, z int) {
	return int(b.Globals.SizeX), int(b.Globals.SizeY), int(b.Globals.SizeZ)
}

// tileIndex returns the index of a tile in MAP.DAT and SEEMAP.DAT, or -1 if
// the position is outside the map.
func (b *Battle) tileIndex(x, y, z int) int {
	sx, sy, sz := b.Size()
	if x < 0 || y < 0 || z < 0 || x >= sx || y >= sy || z >= sz {
		return -1
	}
	return z*sx*sy + y*sx + x
}

// Tile returns the tile at the given position and whether X-COM has seen
// any part of it. Z is the map level in the order of MAP.DAT (0 = top level,
// the same as MapUnit.Z).
func (b *Battle) Tile(x, y, z int) (tile battlescape.MapTile, seen bool) {
	idx := b.tileIndex(x, y, z)
	if idx < 0 || idx >= len(b.Map.Tiles) {
		return battlescape.MapTile{}, false
	}
	if idx < len(b.FogOfWar.Tiles) {
		seen = b.FogOfWar.Tiles[idx] != 0
	}
	return b.Map.Tiles[idx], seen
}

// CraftArea is the area of the map covered by a craft.
type CraftArea struct {
	// Alien is set for UFOs and unset for X-COM crafts.
	Alien bool
	// Type is the craft type as stored in GEODATA.DAT (0 = Skyranger,
	// Lightning, Avenger for X-COM crafts, 0 = small scout, ... for UFOs).
	Type int
	// X, Y, Width and Height is the covered area in tiles. The size is
	// derived from the craft type and is only approximate.
	X, Y, Width, Height int
}

// moduleSize is the size of a map module in tiles.
const moduleSize = 10

// xcomCraftModules and ufoModules are the (approximate) sizes of the crafts
// in map modules (width, height).
var (
	xcomCraftModules = [][2]int{{1, 2}, {2, 1}, {1, 2}}
	ufoModules       = [][2]int{{1, 1}, {1, 1}, {2, 2}, {2, 2}, {2, 2}, {2, 2}, {3, 3}, {3, 3}}
)

// Name returns the name of the craft type.
func (c CraftArea) Name() string {
	t := geoscape.CraftType(c.Type)
	if c.Alien {
		t += geoscape.SmallScout
	}
	return t.String()
}

// CraftAreas returns the areas of the map covered by the X-COM craft and the
// UFO (see GEODATA.DAT). Missions without crafts on the map (e.g. base
// missions) return none.
func (b *Battle) CraftAreas() []CraftArea {
	g := b.Geodata
	area := func(alien bool, typ, x, y int, sizes [][2]int) CraftArea {
		size := [2]int{1, 1}
		if typ >= 0 && typ < len(sizes) {
			size = sizes[typ]
		}
		return CraftArea{
			Alien: alien, Type: typ,
			X: x * moduleSize, Y: y * moduleSize,
			Width: size[0] * moduleSize, Height: size[1] * moduleSize,
		}
	}
	var areas []CraftArea
	if g.Crafts >= 1 {
		areas = append(areas, area(false, int(g.XcomCraft), int(g.XcomCraftX), int(g.XcomCraftY), xcomCraftModules))
	}
	if g.Crafts >= 2 {
		areas = append(areas, area(true, int(g.AlienCraft), int(g.AlienCraftX), int(g.AlienCraftY), ufoModules))
	}
	return areas
}

// TerrainSet is a terrain file (TERRAIN/*.MCD) with the tile type of each
// record (0 = floor, 1 = west wall, 2 = north wall, 3 = object).
type TerrainSet struct {
	Name      string
	TileTypes []uint8
}

// maxTerrainSearch limits the number of combinations MatchTerrainSets tries.
const maxTerrainSearch = 1_000_000

// MatchTerrainSets returns the terrain sets the map is built from, in the
// order of the MCD indices in MAP.DAT.
//
// The savegame does not name the terrain files, it only stores the total
// number of MCD records (WGLOB.DAT). The sets are determined by searching
// for a sequence of sets with this number of records in which the type of
// each record matches the part of the tile it is used for (floor, walls or
// object). A few mismatches are tolerated (e.g. for destroyed parts).
func (b *Battle) MatchTerrainSets(sets []TerrainSet) ([]string, error) {
	total := int(b.Globals.MCDRecords)
	// how often each MCD index is used as floor, west wall, north wall or object
	used := make([][4]int, 256)
	for _, tile := range b.Map.Tiles[:min(b.Globals.Tiles(), len(b.Map.Tiles))] {
		for part, idx := range []uint8{tile.Ground, tile.WestWall, tile.NorthWall, tile.Object} {
			if idx != 0 {
				used[idx][part]++
			}
		}
	}

	var (
		best      []string
		bestScore = -1
		nodes     int
		seq       []int
		types     []uint8
		taken     = make([]bool, len(sets))
	)
	// score returns the number of matching and mismatching uses of the
	// records from index start on.
	score := func(start int) (ok, bad int) {
		for idx := start; idx < len(types) && idx < len(used); idx++ {
			for part, n := range used[idx] {
				if int(types[idx]) == part {
					ok += n
				} else {
					bad += n
				}
			}
		}
		return
	}
	var search func(ok, bad int)
	search = func(ok, bad int) {
		nodes++
		if len(types) == total && ok-bad > bestScore {
			bestScore = ok - bad
			best = make([]string, len(seq))
			for i, s := range seq {
				best[i] = sets[s].Name
			}
		}
		if len(types) >= total || nodes > maxTerrainSearch {
			return
		}
		for i, set := range sets {
			if taken[i] || len(types)+len(set.TileTypes) > total {
				continue
			}
			start := len(types)
			types = append(types, set.TileTypes...)
			o, b := score(start)
			// tolerate about 1% mismatches
			if (bad+b)*100 <= ok+o+100 {
				taken[i] = true
				seq = append(seq, i)
				search(ok+o, bad+b)
				seq = seq[:len(seq)-1]
				taken[i] = false
			}
			types = types[:start]
		}
	}
	search(0, 0)
	if best == nil {
		return nil, fmt.Errorf("no combination of terrain sets matches the map (%d MCD records)", total)
	}
	return best, nil
}
