package resources

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/go-restruct/restruct"
)

// MCDRecordSize is the size of a record in an MCD file.
const MCDRecordSize = 62

// Tile types in MCDRecord.TileType.
const (
	TileFloor     = 0
	TileWestWall  = 1
	TileNorthWall = 2
	TileObject    = 3
)

// MCDRecord describes a part of a Battlescape tile (a floor, a wall or an
// object). The terrain files (TERRAIN/*.MCD) contain one record per part,
// their sprites are stored in the PCK file of the same name.
//
// https://www.ufopaedia.org/index.php/MCD
type MCDRecord struct {
	// 0-7: Sprites of the animation (indices into the PCK file).
	Frames [8]uint8
	// 8-19: Line of fire templates (LOFTEMPS.DAT), from bottom to top.
	LOFT [12]uint8
	// 20-21: Image of the tile on the minimap (SCANG.DAT).
	ScanG uint16
	// 22-29: Unknown.
	Unknown22 [8]uint8
	// 30: 1 for UFO doors (sliding doors).
	UFODoor uint8
	// 31: 1 if the part blocks the line of sight.
	StopLOS uint8
	// 32: 1 if there is no floor (units fall through).
	NoFloor uint8
	// 33: Big wall: 1 = blocks the whole tile, 2-8 = diagonal or partial.
	BigWall uint8
	// 34: 1 for grav lifts.
	GravLift uint8
	// 35: 1 for normal doors.
	Door uint8
	// 36-37: 1 if the part blocks fire and smoke.
	BlockFire  uint8
	BlockSmoke uint8
	// 38: Unknown.
	Unknown38 uint8
	// 39-41: Time units to walk, slide and fly over the part (255 = not
	// passable).
	TUWalk  uint8
	TUSlide uint8
	TUFly   uint8
	// 42: Armor.
	Armor uint8
	// 43: Resistance against explosions.
	HEBlock uint8
	// 44: Record that replaces the part when it is destroyed.
	DieMCD uint8
	// 45: Flammability (255 = not flammable).
	Flammable uint8
	// 46: Record that replaces the part when a door is opened.
	AltMCD uint8
	// 47: Unknown.
	Unknown47 uint8
	// 48: Height offset for units standing on the part.
	TLevel int8
	// 49: Vertical offset of the sprite in pixels.
	PLevel uint8
	// 50: Unknown.
	Unknown50 uint8
	// 51: Light blocked by the part.
	LightBlock uint8
	// 52: Sound of footsteps.
	Footstep uint8
	// 53: Tile type: TileFloor, TileWestWall, TileNorthWall or TileObject.
	TileType uint8
	// 54-55: Explosion created when the part is destroyed.
	HEType     uint8
	HEStrength uint8
	// 56: Smoke created when the part burns.
	SmokeBlockage uint8
	// 57: Turns the part burns.
	Fuel uint8
	// 58: Light emitted by the part.
	LightSource uint8
	// 59: Special type (e.g. UFO power source, navigation).
	TargetType uint8
	// 60-61: Unknown.
	Unknown60 [2]uint8
}

// ParseMCD decodes the records of an MCD file.
func ParseMCD(buf []byte) ([]MCDRecord, error) {
	if len(buf)%MCDRecordSize != 0 {
		return nil, fmt.Errorf("MCD size %d is not a multiple of %d", len(buf), MCDRecordSize)
	}
	records := make([]MCDRecord, len(buf)/MCDRecordSize)
	for i := range records {
		data := buf[i*MCDRecordSize : (i+1)*MCDRecordSize]
		if err := restruct.Unpack(data, binary.LittleEndian, &records[i]); err != nil {
			return nil, fmt.Errorf("could not unpack MCD record %d: %w", i, err)
		}
	}
	return records, nil
}

// LoadMCD loads the records of an MCD file (e.g. TERRAIN/JUNGLE.MCD).
func LoadMCD(path string) ([]MCDRecord, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseMCD(buf)
}
