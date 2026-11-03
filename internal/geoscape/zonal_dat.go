package geoscape

//go:generate stringer -type=Zone -output=zonal_dat_string.go -linecomment

// Zone is one of the 12 geographical zones (regions) the aliens plan their
// missions for. The order is the canonical order used by ZONAL.DAT, ACTS.DAT
// and XBASES.DAT.
type Zone int

const (
	ZoneNorthAmerica   Zone = 0  // North America
	ZoneArctic         Zone = 1  // Arctic
	ZoneAntarctica     Zone = 2  // Antarctica
	ZoneSouthAmerica   Zone = 3  // South America
	ZoneEurope         Zone = 4  // Europe
	ZoneNorthAfrica    Zone = 5  // North Africa
	ZoneSouthernAfrica Zone = 6  // Southern Africa
	ZoneCentralAsia    Zone = 7  // Central Asia
	ZoneSoutheastAsia  Zone = 8  // South East Asia
	ZoneSiberia        Zone = 9  // Siberia
	ZoneAustralasia    Zone = 10 // Australasia
	ZonePacific        Zone = 11 // Pacific
)

// NumZones is the number of zones in ZONAL.DAT, ACTS.DAT and XBASES.DAT.
const NumZones = 12

// ZonalFile is a dynamic list of the chances that the aliens will conduct a
// mission in a particular zone (indexed by Zone).
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
