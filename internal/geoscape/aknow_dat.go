package geoscape

// AknowFile determines whether the aliens know the locations of any of the
// X-COM bases, and if so, which one. The flag is cleared after a base defence
// mission.
//
// The format is undocumented, so the data is kept as raw bytes. The file is
// 495 bytes long.
//
// https://www.ufopaedia.org/index.php/AKNOW.DAT
type AknowFile struct {
	Data [495]byte
}
