package savegame

import "github.com/redtoad/xcom-editor/internal/geoscape"

// HandledFiles returns the names of all data files that are loaded as part of
// the savegame.
func HandledFiles(game *Savegame) []string {
	names := []string{
		"SAVEINFO.DAT", "SOLDIER.DAT", "BASE.DAT", "TRANSFER.DAT",
		"CRAFT.DAT", "LOC.DAT", "LIGLOB.DAT",
	}
	for name := range game.geoscapeFiles() {
		names = append(names, name)
	}
	for name := range (&Battle{}).battleFiles() {
		names = append(names, name)
	}
	return names
}

// CraftsFile gives tests access to the raw data of CRAFT.DAT.
func (game *Savegame) CraftsFile() *geoscape.CraftFile {
	return &game.craftsFile
}

// LocationFile gives tests access to the raw data of LOC.DAT.
func (game *Savegame) LocationFile() *geoscape.LocFile {
	return &game.locationFile
}

// SoldierFile gives tests access to the raw data of SOLDIER.DAT.
func (game *Savegame) SoldierFile() *geoscape.SoldierFile {
	return &game.soldierFile
}
