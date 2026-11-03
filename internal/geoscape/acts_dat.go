package geoscape

// ActsFile holds the priorities of the alien missions ("acts") for each zone
// (indexed by MissionZone). A priority of 0 means the act will not be selected. When
// an act is selected in a month, its priority is set to 0. If it was the last
// non-zero act of the zone, the zone's priority in ZONAL.DAT is set to 0 too.
//
// The file is 84 bytes long (12 zones * 7 bytes).
//
// https://www.ufopaedia.org/index.php/ACTS.DAT
type ActsFile struct {
	Zones [NumZones]ZoneActs
}

// ZoneActs holds the mission priorities for a single zone.
type ZoneActs struct {
	// 0: Alien Research.
	Research uint8
	// 1: Alien Harvest.
	Harvest uint8
	// 2: Alien Abduction (not tried in some zones, e.g. Antarctica).
	Abduction uint8
	// 3: Alien Infiltration.
	Infiltration uint8
	// 4: Alien Base.
	Base uint8
	// 5: Alien Terror. Normally unused, terror missions are generated separately.
	Terror uint8
	// 6: Alien Retaliation. Normally unused, retaliation missions are generated separately.
	Retaliation uint8
}
