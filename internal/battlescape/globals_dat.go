package battlescape

// BglobFile holds the ambient lighting of the map. Light values range from
// 0 (bright) to 16 (dark), e.g. 0 for day missions and 16 for night missions.
//
// The file is 2 bytes long.
//
// https://www.ufopaedia.org/index.php/BGLOB.DAT
type BglobFile struct {
	// 0: Overall ambient light.
	Ambient uint8
	// 1: Ambient light in shadowed areas (e.g. inside crafts), usually
	// Ambient + 4.
	Shadow uint8
}

// WglobFile holds the map dimensions, the current turn, the mission score
// and the speed options of the Battlescape.
//
// The file is 44 bytes long.
//
// https://www.ufopaedia.org/index.php/WGLOB.DAT
type WglobFile struct {
	// 0-1: Current turn * 2.
	TurnTimesTwo uint16
	// 2-3: Current turn (the first turn is 1).
	Turn uint16
	// 4-5: Map size along the X axis.
	SizeX uint16
	// 6-7: Map size along the Y axis.
	SizeY uint16
	// 8-9: Number of map levels (Z axis).
	SizeZ uint16
	// 10-11: SizeX * SizeY.
	Area uint16
	// 12-13: Unknown, appears to be SizeZ - 1.
	Unknown1 uint16
	// 14-15: Number of nodes in ROUTES.DAT.
	Routes uint16
	// 16-17: Number of MCD records.
	MCDRecords uint16
	// 18-19: Aliens killed.
	AliensKilled uint16
	// 20-21: X-COM soldiers lost.
	SoldiersLost uint16
	// 22-23: X-COM tanks lost.
	TanksLost uint16
	// 24-25: Score for killing aliens.
	AlienKillScore int16
	// 26-27: Score for losing X-COM units (negative).
	XcomLossScore int16
	// 28-29: Civilians killed by aliens.
	CiviliansKilledByAliens uint16
	// 30-31: Civilians killed by X-COM.
	CiviliansKilledByXcom uint16
	// 32-33: Scroll speed (6 - value).
	ScrollSpeed uint16
	// 34-35: X-COM movement speed (6 - value).
	XcomSpeed uint16
	// 36-37: Unknown, appears to mirror XcomSpeed.
	Unknown2 uint16
	// 38-39: Alien movement speed (6 - value).
	AlienSpeed uint16
	// 40-41: Fire speed.
	FireSpeed uint16
	// 42-43: Scrolling: 0 = trigger scroll, 1 = auto scroll.
	ScrollMode uint16
}

// Tiles returns the number of tiles of the map.
func (w WglobFile) Tiles() int {
	return int(w.SizeX) * int(w.SizeY) * int(w.SizeZ)
}
