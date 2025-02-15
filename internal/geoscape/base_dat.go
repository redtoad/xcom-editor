package geoscape

//go:generate stringer -type=Facility,Inventory -output=base_dat_string.go -linecomment -trimprefix Inventory

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
)

const maxBases = 8

// Each base entry is 292 bytes long.
const baseByteLength = 292

// BaseFile (BASE.DAT) has all of the base layout and contents information, as well as
// base name info.
//
// https://www.ufopaedia.org/index.php/BASE.DAT
type BaseFile struct {
	Bases []BaseData
}

// SizeOf implemtents the restruct.Sizer interface.
func (s BaseFile) SizeOf() int {
	return baseByteLength * maxBases
}

// Pack implements the restruct.Packer interface.
func (s BaseFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < maxBases; i++ {
		data, err := restruct.Pack(order, &s.Bases[i])
		if err != nil {
			return nil, err
		}
		offset := i * baseByteLength
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

// Unpack implements the restruct.Unpacker interface.
func (s *BaseFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	s.Bases = make([]BaseData, maxBases)
	for i := 0; i < maxBases; i++ {
		offset := i * baseByteLength
		data := buf[offset : offset+baseByteLength]
		if err := restruct.Unpack(data, order, &s.Bases[i]); err != nil {
			return nil, err
		}
	}
	return buf[s.SizeOf():], nil
}

type BaseData struct {

	// 00-0E: BaseData Name, pretty obvious
	// 0F: Presumably the Null character if the BaseData Name uses all 15 characters
	Name string `struct:"[16]byte"`

	// Logical values for the detection capabilities:
	//
	// 10 short, 0 long: This base has small radar(s) only.
	// 20 short, 20 long: This base has large radar(s) only.
	// 30 short, 20 long: This base has small and large radar(s).
	// 100 hyperwave: This base has a hyperwave decoder(s).
	//
	// The radar values can be set to 100 for perfect short range detection (presumably -- it definitely makes
	// UFOs appear more often), but these reset to the correct values any time you complete a build in that base.

	// 10-11: BaseData's short range detection capability.
	ShortRange int `struct:"int16"`

	// 12-13: BaseData's long range detection capability.
	LongRange int `struct:"int16"`

	// 14-15: BaseData's hyperwave detection capability.
	Hyperwave int `struct:"int16"`

	// 16-39: The next offsets are arranged so they're easier to understand. They are for facilities in the base:
	Grid [36]Facility `struct:"[36]uint8"`

	// 3A-5D: The next offsets represent the days until a facility is completed. They're set up the same way:
	DaysToCompletion [36]uint `struct:"[36]uint8"`

	// Number of engineers in the base.
	Engineers int `struct:"int8"`

	// Number of scientists in the base.
	Scientists int `struct:"int8"`

	// 60-11E inventory
	Inventory [96]int `struct:"[96]int16"`

	// 0120: Active/Inactive BaseData. Inactive entries have a value of 1. Active entries have a value of 0. Creating a new base will overwrite the first inactive entry. If a base is dismantled, the only change to the record is this value so it is possible to restore a dismantled base (Access lift removed) by restoring this value to 0. --SeulDragon 12:24, 11 July 2008 (PDT)
	Active bool `struct:"int8,invertedbool"`

	// 0121~0123: 0120 is stored as an integer. These fields are the unused portion of that integer.
}

type Facility uint

const (
	AccessLift     Facility = iota // Access Lift
	LivingQuarters                 // Living Quarters
	Laboratory
	Workshop
	SmallRadarSystem           // Small Radar System
	LargeRadarSystem           // Large Radar System
	MissileDefense             // Missile Defense
	GeneralStores              // General Stores
	AlienContainment           // Alien Containment
	LaserDefense               // Laser Defense
	PlasmaDefense              // Plasma Defense
	FusionBallDefense          // Fusion Ball Defense
	GravShield                 // Grav Shield
	MindShield                 // Mind Shield
	PsionicLaboratory          // Psionic laboratory
	HyperwaveDecoder           // Hyperwave Decoder
	HangarTopLeft              // Hangar (top left)
	HangarTopRight             // Hangar (top right)
	HangarBottomLeft           // Hangar (bottom left)
	HangarBottomRight          // Hangar (bottom right)
	Empty             Facility = 0xff
)

type Inventory int

const (
	InventoryStingrayLauncher  Inventory = iota // stingray launcher
	InventoryAvalancheLauncher                  // avalanche launcher
	InventoryCannon
	InventoryFusionBallLauncher
	InventoryLaserCannon
	InventoryPlasmaBeam
	InventoryStingrayMissile
	InventoryAvalancheMissile
	InventoryCannonRounds
	InventoryFusionBalls
	InventoryTankCannon
	InventoryTankRocketLauncher
	InventoryTankLaserCannon
	InventoryHovertankPlasma
	InventoryHovertankLauncher
	InventoryPistol
	InventoryPistolClip
	InventoryRifle
	InventoryRifleClip
	InventoryHeavyCannon
	InventoryHCAPAmmo
	InventoryHCHEAmmo
	InventoryHCINAmmo
	InventoryAutoCannon
	InventoryACAPAmmo
	InventoryACHEAmmo
	InventoryACINAmmo
	InventoryRocketLauncher
	InventorySmallRocket
	InventoryLargeRocket
	InventoryIncendiaryRocket
	InventoryLaserPistol
	InventoryLaserRifle
	InventoryHeavyLaser
	InventoryGrenade
	InventorySmokeGrenade
	InventoryProximityGrenade
	InventoryHighExplosive
	InventoryMotionScanner
	InventoryMediKit
	InventoryPsiAmp
	InventoryStunRod
	InventoryElectroFlare
	_
	_
	_
	InventoryCorpse
	InventoryCorpseArmour
	InventoryCorpsePowersuit
	InventoryHeavyPlasma
	InventoryHeavyPlasmaClip
	InventoryPlasmaRifle
	InventoryPlasmaRifleClip
	InventoryPlasmaPistol
	InventoryPlasmaPistolClip
	InventoryBlasterLauncher
	InventoryBlasterBomb
	InventorySmallLauncher
	InventoryStunBomb
	InventoryAlienGrenade
	InventoryElerium115
	InventoryMindProbe
	_
	_
	_
	InventorySectoidCorpse
	InventorySnakemanCorpse
	InventoryEtherealCorpse
	InventoryMutonCorpse
	InventoryFloaterCorpse
	InventoryCelatidCorpse
	InventorySilacoidCorpse
	InventoryChryssalidCorpse
	InventoryReaperCorpse
	InventorySectopodCorpse
	InventoryCyberdiscCorpse
	InventoryHovertankCorpse
	InventoryTankCorpse
	InventoryMaleCivilianCorpse
	InventoryFemaleCivilianCorpse
	InventoryUFOPowerSource
	InventoryUFONavigation
	InventoryUFOConstruction
	InventoryAlienFood
	InventoryAlienReproduction
	InventoryAlienEntertainment
	InventoryAlienSurgery
	InventoryExaminationRoom
	InventoryAlienAlloys
	InventoryAlienHabitat
	InventoryPersonalArmour
	InventoryPowerSuit
	InventoryFlyingSuit
	InventoryHWPCannonShell
	InventoryHWPRockets
	InventoryHWPFusionBomb
)
