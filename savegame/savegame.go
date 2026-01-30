package savegame

import (
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type Savegame struct {
	Path          string
	meta          geoscape.SaveinfoFile
	Financials    geoscape.LiglobFile
	BasesData     geoscape.BaseFile
	locationFile  geoscape.LocFile
	soldierFile   geoscape.SoldierFile
	transferFile  geoscape.TransferFile
	craftsFile    geoscape.CraftFile
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
	} {
		if err := fnc(); err != nil {
			return nil, err
		}
	}
	return game, nil
}

// Save saves the entire savegame on disk at its original location.
func (game *Savegame) Save() error {
	for _, fnc := range []loadOrSaveFunc{
		game.saveSoldiers,
		//game.saveBases,
		//game.saveTransfers,
		//game.saveCrafts,
		//game.savelocations,
	} {
		if err := fnc(); err != nil {
			return err
		}
	}
	return nil
}
