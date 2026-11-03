package savegame

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
