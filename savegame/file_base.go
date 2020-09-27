package savegame

import (
	"encoding/binary"
	"fmt"

	"gopkg.in/restruct.v1"
)

type FileBase struct {
	Bases []Base
}

func (s FileBase) SizeOf() int {
	return 292 * 8
}

func (s FileBase) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < 8; i++ {
		data, err := restruct.Pack(order, &s.Bases[i])
		if err != nil {
			return nil, err
		}
		offset := i * 292
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

func (s *FileBase) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	s.Bases = make([]Base, 8)
	for i := 0; i < 8; i++ {
		offset := i * 292
		data := buf[offset : offset+292]
		if err := restruct.Unpack(data, order, &s.Bases[i]); err != nil {
			return nil, err
		}
	}
	return buf[s.SizeOf():], nil
}

// https://www.ufopaedia.org/index.php/BASE.DAT

type Base struct {

	// Each entry is 292 Bytes long

	// 00-0E: Base Name, pretty obvious
	// 0F: Presumably the Null character if the Base Name uses all 15 characters
	Name string `struct:"[16]byte"`

	// Logical values for the detection capabilities:
	//
	//10 short, 0 long: This base has small radar(s) only.
	//
	//20 short, 20 long: This base has large radar(s) only.
	//
	//30 short, 20 long: This base has small and large radar(s).
	//
	//100 hyperwave: This base has a hyperwave decoder(s).
	//
	//The radar values can be set to 100 for perfect short range detection (presumably -- it definitely makes UFOs appear more often), but these reset to the correct values any time you complete a build in that base.

	// 10-11: Base's short range detection capability.
	ShortRange int `struct:"int16"`

	// 12-13: Base's long range detection capability.
	LongRange int `struct:"int16"`

	// 14-15: Base's hyperwave detection capability.
	Hyperwave int `struct:"int16"`

	// 16-39: The next offsets are arranged so they're easier to understand. They are for facilities in the base:
	Grid [36]Facility `struct:"[36]uint8"`

	// 3A-5D: The next offsets represent the days until a facility is completed. They're set up the same way:
	DaysToCompletion [36]uint `struct:"[36]uint8"`

	Engineers  int `struct:"int8"`
	Scientists int `struct:"int8"`

	// 60-11E inventory
	Inventory [96]int `struct:"[96]int16"`

	// 0120: Active/Inactive Base. Inactive entries have a value of 1. Active entries have a value of 0. Creating a new base will overwrite the first inactive entry. If a base is dismantled, the only change to the record is this value so it is possible to restore a dismantled base (Access lift removed) by restoring this value to 0. --SeulDragon 12:24, 11 July 2008 (PDT)
	Active bool `struct:"int8,invertedbool"`

	// 0121~0123: 0120 is stored as an integer. These fields are the unused portion of that integer.
}

type Facility uint

const (
	AccessLift Facility = iota
	LivingQuarters
	Laboratory
	Workshop
	SmallRadarSystem
	LargeRadarSystem
	MissileDefense
	GeneralStores
	AlienContainment
	LaserDefense
	PlasmaDefense
	FusionBallDefense
	GravShield
	MindShield
	PsionicLaboratory
	HyperwaveDecoder
	HangarTopLeft
	HangarTopRight
	HangarBottomLeft
	HangarBottomRight
	Empty Facility = 0xff
)

func (f Facility) Tile() string {
	switch f {
	case Empty:
		return "  "
	case AccessLift:
		return "↑↑"
	case HangarTopLeft:
		return " ⌜"
	case HangarTopRight:
		return "⌝ "
	case HangarBottomLeft:
		return " ⌞"
	case HangarBottomRight:
		return "⌟ "
	default:
		return fmt.Sprintf("%#v", f)
	}
}
