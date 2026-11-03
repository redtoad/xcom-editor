package battlescape

// SmokbitFile is a bit field with one bit per tile (in the same order as in
// MAP.DAT, least significant bit first). A set bit means that the tile has
// smoke or fire (see SMOKREF.DAT). The game recalculates it every turn.
//
// The file is 1250 bytes (10000 bits) long, regardless of the map size.
//
// https://www.ufopaedia.org/index.php/SMOKBIT.DAT
type SmokbitFile struct {
	Bits [1250]byte
}

// HasSmoke returns whether the tile with the given index has smoke or fire.
func (f SmokbitFile) HasSmoke(tile int) bool {
	return f.Bits[tile/8]&(1<<(tile%8)) != 0
}

// SmokrefFile holds the smoke and fire on the map.
//
// The file is 3600 bytes long (400 entries * 9 bytes).
//
// https://www.ufopaedia.org/index.php/SMOKREF.DAT
type SmokrefFile struct {
	Clouds [400]SmokeCloud
}

// SmokeCloud is smoke or fire on a single tile.
type SmokeCloud struct {
	// 0-2: Position.
	Y uint8
	X uint8
	Z uint8
	// 3: Smoke intensity (and remaining time).
	Smoke uint8
	// 4: Remaining time of the fire.
	Fire uint8
	// 5: Intensity of new spread (only used at runtime).
	Spread uint8
	// 6: Type: 0 = empty, 1 = fire, 2 = smoke.
	Type uint8
	// 7: Visibility flag.
	Visible uint8
	// 8: Animation frame offset (0-3).
	Frame uint8
}
