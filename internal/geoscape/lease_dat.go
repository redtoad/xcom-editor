package geoscape

// LeaseFile holds the view of the Geoscape globe. Rotations are stored in
// 1/8 degree (0-2880 = 0-360 degrees).
//
// Only the first bytes are documented, the rest is kept as raw bytes. The
// file is 180 bytes long.
//
// Only RotationY, RotationX and Zoom differ between savegames. All unknown
// bytes are identical in every savegame seen and Unknown3 looks like x86
// machine code, so they are most likely uninitialized memory of the buffer
// GEOSCAPE.EXE allocates for the file.
//
// https://www.ufopaedia.org/index.php/LEASE.DAT
type LeaseFile struct {
	// 0-1: Y rotation of the globe.
	RotationY int16
	// 2-5: Unknown.
	Unknown1 [4]byte
	// 6-7: X rotation of the globe.
	RotationX int16
	// 8-11: Unknown.
	Unknown2 [4]byte
	// 12-13: Zoom (globe radius in pixels: 90, 120, 180, 360, 450 or 720).
	Zoom int16
	// 14-179: Unknown.
	Unknown3 [166]byte
}
