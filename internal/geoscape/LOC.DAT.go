package geoscape

import "fmt"

// Location Data for bases and crafts

type LOC_DAT struct {
	Objects [50]Object
}

func (o Object) String() string {
	return fmt.Sprintf("{%#v %v}", o.Type, o.TableReference)
}

type ObjectType uint

const (
	Unused ObjectType = iota
	AlienShip
	XCOMShip
	XCOMBase
	AlienBase
	CrashSite
	LandedUFO
	Waypoint
	TerrorSite
)

type Object struct {
	Type           ObjectType `struct:"uint8"`
	TableReference uint       `struct:"uint8"`
	X              int        `struct:"int16"`
	Y              int        `struct:"int16"`
}

func (o Object) SizeOf() uint {
	return 20
}
