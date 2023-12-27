package geoscape

import "fmt"

// LOC.DAT has a row width of 20 bytes. There are a total of 50 records (not all of them necessarily used) for a fixed
// file size of 1,000 bytes.
// https://www.ufopaedia.org/index.php/LOC.DAT

// Location Data for bases and crafts
type LOC_DAT struct {
	Objects [50]LocationObject
}

type LocationObject struct {
	// Object type
	Type LocationType `struct:"uint8"`

	// Object table reference - Possible values - 00 to FF - Just a reference. This just shows how many there are of
	// this type on the geoscape. If the object is either a UFO (Alien Ship for TFTD) or X-COM craft, then this is the
	// index into CRAFT.DAT. If the object is an X-Com base, then this is the index into BASE.DAT and if it is an Alien
	// Base this byte contains the race:
	TableReference uint `struct:"uint8"`

	// Horizontal coordinates or longitude (low bit then high bit respectively). Value range: 0 - 2880
	X int `struct:"int16"`

	// Vertical coordinates or latitude (low bit then high bit respectively). Value range: -720 - 720
	Y int `struct:"int16"`

	// For crash site or terror site - countdown timer (in hours). For moving objects - how many game ticks (5s) have to
	// pass until craft moved to next globe coordinate (cell_size div speed). Note that ground UFOs are treated as moving
	// objects, except for speed = 0.
	Timer int `struct:"int16"`

	// Fractional part of how much is left to the next globe coordinate (cell_size mod speed), used only for moving objects.
	Fraction int `struct:"int16"`

	// Count suffix of the item, eg: Skyranger-1 or Crash Site-47. It appears to have no meaning for XCOM Bases, but for
	// other types where it is set, 0B is the high byte for when you go over 255 UFO's or crafts, etc.
	CountSuffix int `struct:"int16"`

	// unused
	_ int `struct:"int16"`

	// Craft transfer mode
	TransferMode int `struct:"int8"`

	// unused
	_ int `struct:"int8"`

	// Globe object visiblity/mobility bitfield
	Visibility int `struct:"int32"`
}

// SizeOf implemtents the restruct.Sizer interface.
func (o LocationObject) SizeOf() uint {
	return 20
}

func (o LocationObject) String() string {
	return fmt.Sprintf("{%v %v}", o.Type, o.TableReference)
}

// LocationType is the the type of object for which a location is stored.
type LocationType uint

const (
	Unused LocationType = iota
	AlienShip
	XCOMShip
	XCOMBase
	AlienBase
	CrashSite
	LandedUFO
	Waypoint
	TerrorSite
)

func (lt LocationType) String() string {
	values := []string{
		"Unused", "AlienShip", "XCOMShip", "XCOMBase", "AlienBase",
		"CrashSite", "LandedUFO", "Waypoint", "TerrorSite",
	}
	return values[lt]
}
