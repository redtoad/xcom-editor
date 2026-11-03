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
	for i, pos := range b.UnitPositions.Units {
		unit := b.Units.Units[i]
		if pos.Y == 0xFF || unit.Type == 0xFF {
			continue // unused slot
		}
		// large units occupy 4 slots, only count the upper left one
		if pos.Flags&battlescape.UnitLarge != 0 &&
			pos.Flags&(battlescape.UnitRightSection|battlescape.UnitLowerSection) != 0 {
			continue
		}
		alive := pos.Flags&battlescape.UnitActive != 0 && unit.Health > 0
		switch pos.Owner {
		case battlescape.OwnerXcom:
			status.XcomTotal++
			if alive {
				status.XcomAlive++
			}
		case battlescape.OwnerAlien:
			status.AliensTotal++
			if alive {
				status.AliensAlive++
			}
		}
	}
	return status
}
