package geoscape

// XbasesFile holds information about alien attacks on X-COM bases, one record
// for each zone (indexed by MissionZone).
//
// The format is only partially documented. The file is 48 bytes long
// (12 zones * 4 bytes).
//
// https://www.ufopaedia.org/index.php/XBASES.DAT
type XbasesFile struct {
	Zones [NumZones]XbaseTarget
}

// XbaseTarget holds the alien targeting information for a single zone.
type XbaseTarget struct {
	// 0-1: Non-zero if and only if an X-COM base has been detected.
	Detected uint16
	// 2-3: LOC.DAT reference of the X-COM base to be targeted.
	Target uint16
}
