package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type Savegame struct {
	Path         string
	meta         geoscape.SaveinfoFile
	Financials   geoscape.LiglobFile
	BasesData    geoscape.BaseFile
	locationFile geoscape.LocFile
	soldierFile  geoscape.SoldierFile
	transferFile geoscape.TransferFile
	craftsFile   geoscape.CraftFile

	// AlienActivity holds the monthly alien activity per country and region (ALIEN.DAT).
	AlienActivity geoscape.AlienFile
	// XcomActivity holds the monthly X-COM activity per country and region (XCOM.DAT).
	XcomActivity geoscape.XcomFile
	// Counters holds the naming counters, current year and monthly scores (UIGLOB.DAT).
	Counters geoscape.UiglobFile
	// ZonePriorities holds the chances of alien missions per zone (ZONAL.DAT).
	ZonePriorities geoscape.ZonalFile
	// AlienActs holds the alien mission priorities per zone (ACTS.DAT).
	AlienActs geoscape.ActsFile
	// AlienBaseTargets holds the alien targeting of X-COM bases per zone (XBASES.DAT).
	AlienBaseTargets geoscape.XbasesFile
}

type loadOrSaveFunc func() error

// Load loads a savegame from disk. This includes loading all required data files
// one by one.
func Load(root string) (*Savegame, error) {
	game := &Savegame{Path: root}
	for _, fnc := range []loadOrSaveFunc{
		game.loadMetadata, // read-only
		game.loadSoldiers,
		game.loadBases,
		game.loadTransfers,
		game.loadCrafts,
		game.loadLocations,
		game.loadFinancials,
		game.loadGeoscapeFiles,
	} {
		if err := fnc(); err != nil {
			return nil, err
		}
	}
	return game, nil
}

// Reload reloads the savegame from disk, discarding any unsaved changes.
func (game *Savegame) Reload() error {
	for _, fnc := range []loadOrSaveFunc{
		game.loadMetadata,
		game.loadSoldiers,
		game.loadBases,
		game.loadTransfers,
		game.loadCrafts,
		game.loadLocations,
		game.loadFinancials,
		game.loadGeoscapeFiles,
	} {
		if err := fnc(); err != nil {
			return err
		}
	}
	return nil
}

func (game *Savegame) loadFinancials() error {
	filePath := path.Join(game.Path, "LIGLOB.DAT")
	if err := internal.LoadDATFile(filePath, &game.Financials); err != nil {
		return fmt.Errorf("could not load LIGLOB.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveFinancials() error {
	filePath := path.Join(game.Path, "LIGLOB.DAT")
	return internal.SaveDATFile(filePath, game.Financials)
}

// Save saves the entire savegame on disk at its original location.
func (game *Savegame) Save() error {
	for _, fnc := range []loadOrSaveFunc{
		game.saveSoldiers,
		game.saveBases,
		game.saveTransfers,
		game.saveCrafts,
		game.saveFinancials,
		game.saveGeoscapeFiles,
	} {
		if err := fnc(); err != nil {
			return err
		}
	}
	return nil
}
