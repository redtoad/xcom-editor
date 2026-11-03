package battlescape

import "encoding/binary"

// MapFile holds the tiles of the map. Tiles are ordered along X, then Y,
// then from the top level to the bottom level (see WglobFile for the
// dimensions). The file may contain stale tiles of a larger map at the end.
//
// The file is SizeX * SizeY * SizeZ * 4 bytes long.
//
// https://www.ufopaedia.org/index.php/MAP.DAT
type MapFile struct {
	Tiles []MapTile
}

// MapTile holds the parts of a single tile. Each part is an index into the
// MCD records of the map, 0 if the part does not exist.
type MapTile struct {
	Ground    uint8
	WestWall  uint8
	NorthWall uint8
	Object    uint8
}

const mapTileSize = 4

// SizeOf implements restruct.Sizer.
func (f MapFile) SizeOf() int {
	return len(f.Tiles) * mapTileSize
}

// Unpack implements restruct.Unpacker.
func (f *MapFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	tiles, err := unpackRecords[MapTile](buf, order, mapTileSize)
	if err != nil {
		return nil, err
	}
	f.Tiles = tiles
	return buf[f.SizeOf():], nil
}

// Pack implements restruct.Packer.
func (f MapFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	return buf, packRecords(buf, order, mapTileSize, f.Tiles)
}
