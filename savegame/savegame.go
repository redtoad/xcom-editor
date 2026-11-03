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
	// AlienMissions holds the ongoing alien missions per zone (MISSIONS.DAT).
	AlienMissions geoscape.MissionsFile
	// AlienKnowledge holds whether the aliens know about X-COM bases (AKNOW.DAT).
	AlienKnowledge geoscape.AknowFile
	// AlienContainment holds the live aliens in containment (ASTORE.DAT).
	AlienContainment geoscape.AstoreFile
	// GameState holds the current date and time and some counters (IGLOB.DAT).
	GameState geoscape.IglobFile
	// Interceptions holds the open interception windows (INTER.DAT).
	Interceptions geoscape.InterFile
	// GlobeView holds the rotation and zoom of the Geoscape globe (LEASE.DAT).
	GlobeView geoscape.LeaseFile
	// Countries holds the relations to the funding countries (DIPLOM.DAT).
	Countries geoscape.DiplomFile
	// Facilities holds the stats of the base facility types (FACIL.DAT).
	Facilities geoscape.FacilFile
	// Products holds the manufacturing data of all items (PRODUCT.DAT).
	Products geoscape.ProductFile
	// Production holds the manufacturing progress per base (BPROD.DAT).
	Production geoscape.BprodFile
	// ResearchTopics holds all research topics (RESEARCH.DAT).
	ResearchTopics geoscape.ResearchFile
	// Research holds the research progress per base (PROJECT.DAT).
	Research geoscape.ProjectFile
	// Purchases holds prices and transfer times of all items (PURCHASE.DAT).
	Purchases geoscape.PurchaseFile
	// Ufopaedia holds the UFOpaedia entries (UP.DAT).
	Ufopaedia geoscape.UpFile
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
