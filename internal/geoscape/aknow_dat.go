package geoscape

// AknowFile is meant to hold the alien knowledge (according to UFOpaedia,
// whether the aliens know the locations of X-COM bases). In this game
// version the data is unused and always 0, see below; base detection is
// stored in XBASES.DAT.
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
// References in GEOSCAPE.EXE (linear addresses): allocation 0x21808, new
// game 0x21dca, save 0x14dde, load 0x15b8d.
//
// https://www.ufopaedia.org/index.php/AKNOW.DAT
type AknowFile struct {
	Data [495]byte
}
