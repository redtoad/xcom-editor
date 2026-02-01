package geoscape

//go:generate stringer -type=CraftType,Altitude,FlightMode,WeaponType,MissionType,MissionZone,CraftStatus -output=craft_dat_string.go -linecomment

import (
	"encoding/binary"
	"fmt"

	"github.com/go-restruct/restruct"
)

type CraftFile struct {
	Crafts []CraftData
}

const maxCrafts = 50
const craftByteLength = 104

func (cf CraftFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < maxCrafts; i++ {
		data, err := restruct.Pack(order, &cf.Crafts[i])
		if err != nil {
			return nil, err
		}
		offset := i * craftByteLength
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

func (cf *CraftFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	cf.Crafts = make([]CraftData, maxCrafts)
	for i := 0; i < maxCrafts; i++ {
		offset := i * craftByteLength
		data := buf[offset : offset+craftByteLength]
		if err := restruct.Unpack(data, order, &cf.Crafts[i]); err != nil {
			return nil, fmt.Errorf("could not unpack []Crafts: %w", err)
		}
	}
	return buf[cf.SizeOf():], nil
}

func (cf *CraftFile) SizeOf() int {
	return maxCrafts * craftByteLength
}

type CraftData struct {

	// 0	0x00	Craft type, Possible values are the same as GEODATA.DAT:
	// *HUMAN*          *ALIEN*
	// 0 - Skyranger     5 - Small Scout      255 - Entry Not Used
	// 1 - Lightning     6 - Medium Scout
	// 2 - Avenger       7 - Large Scout
	// 3 - Interceptor   8 - Harvester
	// 4 - Firestorm     9 - Abductor
	//                  10 - Terror Ship
	//                  11 - Battleship
	//                  12 - Supply Ship
	Type CraftType `struct:"int8"`

	// Offsets 1 and 5 refer to the weapon placed in the left and right slots respectively. (The Lightning does not have a center weapon type, only a left). Their possible values are listed here:
	// 0 - Stingray
	// 1 - Avalanche
	// 2 - Cannon
	// 3 - Fusion Ball
	// 4 - Laser Cannon
	// 5 - Plasma Beam
	// 255 - No Weapon

	// 1	0x01	Left weapon type
	LeftWeapon WeaponType `struct:"int8"`

	// 2-3	0x02-0x03	Left ammo
	LeftAmmo int `struct:"int16"`

	// 4	0x04	Indicates flight mode
	FlightMode FlightMode `struct:"int8"`

	// 5	0x05	Right weapon type
	RightWeapon WeaponType `struct:"int8"`

	// 6-7	0x06-0x07	Right ammo
	RightAmmo int `struct:"int16"`

	// 8-9	0x08-0x09	Unused.
	Unused int `struct:"int16"`

	// Damage, that is the amount it currently has taken. This value divided by the
	// crafts damage capacity gives the percentage shown in-game.
	Damage int `struct:"int16"`

	// 12-13	0x0C-0x0D	Altitude of craft. Is a index within ENGLISH.DAT for string.
	// 0 = GROUND *
	// 1 = VERY LOW
	// 2 = LOW
	// 3 = HIGH
	// 4 = VERY HIGH
	// (*NOTE: If craft is airborne and you change it to this value the altitude will remain the same. Speed must be edited to 0 for the change to hold.)
	Altitude Altitude `struct:"int16"`

	// 14-15	0x0E-0x0F	Speed of craft.
	Speed int `struct:"int16"`

	// 16-17	0x10-0x11	Index into LOC.DAT referencing the destination - for example, waypoints for X-COM craft, or X-COM bases for alien craft.
	Destination int `struct:"int16"`

	// 18-19	0x12-0x13	Index into INTER.DAT when the ship is in interception mode.
	InterceptionRef int `struct:"int16"`

	// 20-21	0x14-0x15	Next UFO waypoint coordinate X (longitude).
	NextUFOWaypointLon int `struct:"int16"`

	// 22-23	0x16-0x17	Next UFO waypoint coordinate Y (latitude).
	NextUFOWaypointLat int `struct:"int16"`

	// 24-25	0x18-0x19	Fuel, amount remaining. This value divided by the crafts total fuel capacity gives the percentage shown in-game.
	FuelAmountRemaining int `struct:"int16"`

	// 26-27	0x1A-0x1B	BaseData reference as an index to LOC.DAT.
	BaseReference Altitude `struct:"int16"`

	// 28-29	0x1C-0x1D	Mission type craft is on. Is an index within ENGLISH.DAT for the string (558 + this value).
	MissionType MissionType `struct:"int16"`

	// 30-31	0x1E-0x1F	Zone where mission is being carried out. Is an index within ENGLISH.DAT for string (543 + this value).
	MissionZone MissionZone `struct:"int16"`

	// 32-33	0x20-0x21	UFO trajectory segment (ranges from 0-7).
	UFOTrajectorySegment int `struct:"int16"`

	// 34-35	0x22-0x23	UFO trajectory type (ranges from 0-9).
	UFOTrajectoryType int `struct:"int16"`

	// 36-37	0x24-0x25	Alien Race found on craft. Is index within ENGLISH.DAT for the string (466 + this value).
	//     UFO                TFTD
	// 0 = Sectoid            Aquatoid
	// 1 = Snakeman           Gillman
	// 2 = Ethereal           Lobsterman
	// 3 = Muton              Tasoth
	// 4 = Floater            Mixed Crew (Type I)
	// 5 = Final mission mix  Mixed Crew (Type II)
	AlienRace int `struct:"int16"`

	// 38-39	0x26-0x27	UFO attack timer.
	UFOAttackTimer int `struct:"int16"`

	// 40-41	0x28-0x29	UFO escape manuever timer.
	UFOEscapeManueverTimer int `struct:"int16"`

	// 42-43	0x2A-0x2B	Craft status. Is an index within ENGLISH.DAT for the string (268 + this value).
	Status CraftStatus `struct:"int16"`

	// 44-103: Cargo items and flags. Offsets 44-98 contain item counts on board the craft
	// (mapping to OBDATA.DAT entries). Offsets 100-103 contain a bitfield for craft state flags.
	Cargo [60]byte `struct:"[60]byte"`
}

// SizeOf implements restruct.Sizer
func (c CraftData) SizeOf() int {
	return 104
}

type CraftType int

const (
	Skyranger CraftType = iota
	Lightning
	Avenger
	Interceptor
	Firestorm

	SmallScout  // small scout
	MediumScout // medium scout
	LargeScout  // large scout
	Harvester
	Abductor

	TerrorShip // terror ship
	Battleship
	SupplyShip // supply ship

	EntryNotUsed = -1 // entry not used
)

type Altitude int

const (
	Ground  Altitude = iota
	VeryLow          // very low
	Low
	High
	VeryHigh // very high
)

type FlightMode int

const (
	NoDestination        FlightMode = iota // No destination set (at base)
	SingleDestination                      // Single destination
	MultipleDestinations                   // Multiple destinations (UFO only)

)

type WeaponType int

const (
	Stingray WeaponType = iota
	Avalanche
	Cannon
	FusionBall  // fusion ball
	LaserCannon // laser cannon
	PlasmaBeam  // plasma beam
	NoWeapon    = 255
)

type MissionType int

const (
	MissionAlienResearch     MissionType = iota // Alien Research
	MissionAlienHarvest                         // Alien Harvest
	MissionAlienAbduction                       // Alien Abduction
	MissionAlienInfiltration                    // Alien Infiltration
	MissionAlienBase                            // Alien Base
	MissionAlienTerror                          // Alien Terror
	MissionAlienRetaliation                     // Alien Retaliation
	MissionAlienSupply                          // Alien Supply
)

type MissionZone int

const (
	NorthAmerica   MissionZone = iota // North America
	Arctic                            // Arctic
	Antarctica                        // Antarctica
	SouthAmerica                      // South America
	Europe                            // Europe
	NorthAfrica                       // North Africa
	SouthernAfrica                    // Southern Africa
	CentralAsia                       // Central Asia
	SouthEastAsia                     // South East Asia
	Siberia                           // Siberia
	Australasia                       // Australasia
	Pacific                           // Pacific
	NorthAtlantic                     // North Atlantic (unused)
	SouthAtlantic                     // South Atlantic (unused)
	IndianOcean                       // Indian Ocean (unused)
)

type CraftStatus int

const (
	Ready CraftStatus = iota
	Out
	Repairs
	Refueling
	Rearming
)
