package geoscape

// AknowFile determines whether the aliens know the locations of any of the
// X-COM bases, and if so, which one. The flag is cleared after a base defence
// mission.
//
// The format is undocumented, so the data is kept as raw bytes. The file is
// 495 bytes long.
//
// GEOSCAPE.EXE allocates it as "alienknowledge" and clears it for a new
// game as 5 records of 99 bytes: 12 * 7 bytes (zones * mission types, like
// ACTS.DAT), 8 bytes at offset 84 and 7 bytes at offset 92. Apart from
// loading and saving, no code in GEOSCAPE.EXE or TACTICAL.EXE accesses it,
// so it is always 0.
//
// https://www.ufopaedia.org/index.php/AKNOW.DAT
type AknowFile struct {
	Data [495]byte
}
