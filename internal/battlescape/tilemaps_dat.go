package battlescape

import (
	"encoding/binary"
	"fmt"
)

// Flags in SeemapFile.Tiles.
const (
	SeenWestWall  = 1 << 0
	SeenNorthWall = 1 << 1
	SeenContent   = 1 << 2
)

// SeemapFile holds the fog of war: which parts of each tile X-COM has seen
// (see SeenWestWall, SeenNorthWall and SeenContent). Tiles are stored in the
// same order as in MAP.DAT.
//
// The file is one byte per tile long.
//
// https://www.ufopaedia.org/index.php/SEEMAP.DAT
type SeemapFile struct {
	Tiles []uint8
}

// SizeOf implements restruct.Sizer.
func (f SeemapFile) SizeOf() int {
	return len(f.Tiles)
}

// Unpack implements restruct.Unpacker.
func (f *SeemapFile) Unpack(buf []byte, _ binary.ByteOrder) ([]byte, error) {
	f.Tiles = append([]uint8(nil), buf...)
	return buf[len(buf):], nil
}

// Pack implements restruct.Packer.
func (f SeemapFile) Pack(buf []byte, _ binary.ByteOrder) ([]byte, error) {
	copy(buf, f.Tiles)
	return buf, nil
}

// LightMapFile holds the light level of each tile, from 0 (bright) to 15
// (dark). Soldiers can only see tiles with a light level of 9 or less. The
// effective light level is the lower value of TERMP.DAT (fixed and ambient
// light) and SOURCEMP.DAT (mobile light sources such as units). Tiles are
// stored in the same order as in MAP.DAT.
//
// The file is one byte per tile plus 4 unknown bytes long.
//
// https://www.ufopaedia.org/index.php/TERMP.DAT
// https://www.ufopaedia.org/index.php/SOURCEMP.DAT
type LightMapFile struct {
	Tiles []uint8
	// Unknown trailing bytes.
	Trailer [4]byte
}

// SizeOf implements restruct.Sizer.
func (f LightMapFile) SizeOf() int {
	return len(f.Tiles) + len(f.Trailer)
}

// Unpack implements restruct.Unpacker.
func (f *LightMapFile) Unpack(buf []byte, _ binary.ByteOrder) ([]byte, error) {
	if len(buf) < len(f.Trailer) {
		return nil, fmt.Errorf("light map too short: %d bytes", len(buf))
	}
	n := len(buf) - len(f.Trailer)
	f.Tiles = append([]uint8(nil), buf[:n]...)
	copy(f.Trailer[:], buf[n:])
	return buf[len(buf):], nil
}

// Pack implements restruct.Packer.
func (f LightMapFile) Pack(buf []byte, _ binary.ByteOrder) ([]byte, error) {
	copy(buf, f.Tiles)
	copy(buf[len(f.Tiles):], f.Trailer[:])
	return buf, nil
}
