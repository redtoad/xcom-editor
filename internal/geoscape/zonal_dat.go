package geoscape

// NumZones is the number of zones (see MissionZone) in ZONAL.DAT, ACTS.DAT,
// XBASES.DAT and MISSIONS.DAT. The zones are stored in the order of MissionZone.
const NumZones = 12

// ZonalFile is a dynamic list of the chances that the aliens will conduct a
// mission in a particular zone (indexed by MissionZone).
//
// A zone with priority 0 will only get terror attacks. Once research, harvest,
// abduction, infiltration and base creation missions have all been performed
// in a zone (see ACTS.DAT), the value for that zone is set to 0.
//
// The file is 12 bytes long.
//
// https://www.ufopaedia.org/index.php/ZONAL.DAT
type ZonalFile struct {
	Priorities [NumZones]uint8
}
