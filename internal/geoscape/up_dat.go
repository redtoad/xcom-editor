package geoscape

// UpFile holds the UFOpaedia entries.
//
// The file is 1392 bytes long (116 entries * 12 bytes). A new game copies the
// table from GEOSCAPE.EXE (offset 386176); only Research changes during the
// game, all unknown fields are constants.
//
// https://www.ufopaedia.org/index.php/UP.DAT
type UpFile struct {
	Entries [116]UpEntry
}

// UpEntry is a single UFOpaedia entry.
type UpEntry struct {
	// 0: Category (0-8).
	Category uint8
	// 1: Unknown, always 0.
	Unknown1 uint8
	// 2: Index within the category. For categories 2 and 3 this is the index
	// in OBDATA.DAT, armours use 100-102.
	SubIndex uint8
	// 3-4: Unknown, always 0.
	Unknown2 [2]uint8
	// 5: Image number (GEOGRAPH/UP0##.SPK).
	Image uint8
	// 6-7: Index of the title in ENGLISH.DAT (1-based).
	StringIndex uint16
	// 8: Research status: 0 = not researched, 2 = researched, 3 = available
	// from the start.
	Research uint8
	// 9-11: Unknown, always 0.
	Unknown3 [3]uint8
}
