package geoscape

// FacilFile holds the stats of the 17 base facility types (in the order of
// Facility).
//
// The file is 272 bytes long (17 facilities * 16 bytes).
//
// https://www.ufopaedia.org/index.php/FACIL.DAT
type FacilFile struct {
	Facilities [17]FacilityStats
}

// FacilityStats holds the stats of a single base facility type.
type FacilityStats struct {
	// 0-1: Index of the facility name in ENGLISH.DAT.
	NameIndex uint16
	// 2-3: Build cost in $1000.
	BuildCost int16
	// 4: Build time in days.
	BuildDays uint8
	// 5: Monthly upkeep in $1000.
	Upkeep uint8
	// 6: Capacity (e.g. personnel, storage or containment space).
	Capacity uint8
	// 7: Unknown.
	Unknown1 uint8
	// 8-9: Defense strength.
	Defense int16
	// 10: Defense hit ratio in percent.
	HitRatio uint8
	// 11: Unknown.
	Unknown2 uint8
	// 12: Build flags: 0x01 = special facility, 0x02 = defense facility,
	// 0x04 = can be built now, 0x08 = requires research.
	Flags uint8
	// 13-15: Unknown.
	Unknown3 [3]uint8
}
