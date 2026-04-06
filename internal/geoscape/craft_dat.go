package geoscape

//go:generate stringer -type=CraftType,Altitude,FlightMode,WeaponType,MissionType,MissionZone,CraftStatus -output=craft_dat_string.go -linecomment

import (
	"encoding/binary"
	"fmt"

	"github.com/go-restruct/restruct"
)

const maxCrafts = 50
const craftByteLength = 104

// CraftFile represents the CRAFT.DAT file with its fixed set of geoscape craft entries.
type CraftFile struct {
	Crafts [maxCrafts]CraftData
}

// SizeOf implements restruct.Sizer.
func (cf CraftFile) SizeOf() int {
	return maxCrafts * craftByteLength
}

// Pack implements restruct.Packer.
func (cf CraftFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := range maxCrafts {
		data, err := restruct.Pack(order, &cf.Crafts[i])
		if err != nil {
			return nil, fmt.Errorf("could not pack CraftData: %w", err)
		}

		offset := i * craftByteLength
		for j := range data {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

// Unpack implements restruct.Unpacker.
func (cf *CraftFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := range maxCrafts {
		offset := i * craftByteLength
		data := buf[offset : offset+craftByteLength]
		if err := restruct.Unpack(data, order, &cf.Crafts[i]); err != nil {
			return nil, fmt.Errorf("could not unpack CraftData: %w", err)
		}
	}
	return buf[cf.SizeOf():], nil
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

	// 8-9	0x08-0x09	Unused
	// Note: Normally we would use `_` for unused fields, but in this case we want to preserve the exact byte layout of the
	// struct, so we use a named field and just ignore it.
	UnusedWord int16 `struct:"int16"`

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

	// Bytes 44-103: on-board item quantities (bytes 44-99) and status bitfield (bytes 100-103).
	// The rest of the known values detail the items on board of the craft. 49-98 refer to offsets 0-49 in OBDATA.DAT.
	// 44	0x2C	Tank/Cannon
	// 45	0x2D	Tank/Rocket Launcher
	// 46	0x2E	Tank/Laser Cannon
	// 47	0x2F	Hover Tank/Plasma
	// 48	0x30	Hover Tank/Launcher
	// 49	0x31	PISTOL
	// 50	0x32	PISTOL CLIP
	// 51	0x33	RIFLE
	// 52	0x34	RIFLE CLIP
	// 53	0x35	HEAVY CANNON
	// 54	0x36	CANNON AP-AMMO
	// 55	0x37	CANNON HE-AMMO
	// 56	0x38	CANNON I-AMMO
	// 57	0x39	AUTO-CANNON
	// 58	0x3A	AUTO-CANNON AP-AMMO
	// 59	0x3B	AUTO-CANNON HE-AMMO
	// 60	0x3C	AUTO-CANNON I-AMMO
	// 61	0x3D	ROCKET LAUNCHER
	// 62	0x3E	SMALL ROCKET
	// 63	0x3F	LARGE ROCKET
	// 64	0x40	INCENDIARY ROCKET
	// 65	0x41	LASER PISTOL
	// 66	0x42	LASER GUN
	// 67	0x43	HEAVY LASER
	// 68	0x44	GRENADE
	// 69	0x45	SMOKE GRENADE
	// 70	0x46	PROXIMITY GRENADE
	// 71	0x47	HIGH EXPLOSIVE
	// 72	0x48	MOTION SCANNER
	// 73	0x49	MEDI-KIT
	// 74	0x4A	PSI-AMP
	// 75	0x4B	STUN ROD
	// 76	0x4C	Flare
	// 77	0x4D	empty
	// 78	0x4E	empty
	// 79	0x4F	empty
	// 80	0x50	CORPSE
	// 81	0x51	CORPSE & ARMOUR
	// 82	0x52	CORPSE & POWER SUIT
	// 83	0x53	Heavy Plasma
	// 84	0x54	Heavy Plasma Clip
	// 85	0x55	Plasma Rifle
	// 86	0x56	Plasma Rifle Clip
	// 87	0x57	Plasma Pistol
	// 88	0x58	Plasma Pistol Clip
	// 89	0x59	BLASTER LAUNCHER
	// 90	0x5A	BLASTER BOMB
	// 91	0x5B	SMALL LAUNCHER
	// 92	0x5C	STUN MISSILE
	// 93	0x5D	ALIEN GRENADE
	// 94	0x5E	ELERIUM-115
	// 95	0x5F	MIND PROBE
	// 96	0x60	>>UNDEFINED <<
	// 97	0x61	>> empty <<
	// 98	0x62	>> empty <<
	// 99	0x63	No known use for this offset.

	// 100-103	0x64-0x67	Bitfield (only 1st byte is used):
	// bit 0  (1): Craft is in hangar (0=craft's been flew away)
	// bit 1  (2): Want to go home (set after mission or when out of fuel), blocks craft
	//             from re-targeting
	// bit 2  (4): [runtime] Out of elerium (0=fueled)
	// bit 3  (8): [runtime] Out of left ammo (0=left weapon rearmed)
	// bit 4 (16): [runtime] Out of right ammo (0=right weapon rearmed)
	// bit 5 (32): [runtime] Flag of processed UFO in interception window, to avoid
	//             multiple escape/attack timers decrementing in case when UFO is
	//             pursued in more than 1 window at once.
	// bit 6 (64): Hyperwaved extra info
	Items [60]byte `struct:"[60]uint8"`
}

// SizeOf implements restruct.Sizer
func (c *CraftData) SizeOf() int {
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
