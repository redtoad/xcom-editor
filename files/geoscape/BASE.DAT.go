package geoscape

import (
	"encoding/binary"
	"fmt"

	"github.com/go-restruct/restruct"
)

const maxBases = 8

// Each base entry is 292 bytes long.
const baseByteLength = 292

// BASE_DAT has all of the base layout and contents information, as well as
// base name info.
//
// https://www.ufopaedia.org/index.php/BASE.DAT
type BASE_DAT struct {
	Bases []Base
}

func (s BASE_DAT) SizeOf() int {
	return baseByteLength * maxBases
}

func (s BASE_DAT) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
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

func (s *BASE_DAT) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	s.Bases = make([]Base, maxBases)
	for i := 0; i < maxBases; i++ {
		offset := i * baseByteLength
		data := buf[offset : offset+baseByteLength]
		if err := restruct.Unpack(data, order, &s.Bases[i]); err != nil {
			return nil, err
		}
	}
	return buf[s.SizeOf():], nil
}

type Base struct {

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
	case LivingQuarters:
		return "LQ"
	case SmallRadarSystem:
		return "SR"
	case LargeRadarSystem:
		return "LR"
	case Workshop:
		return "WS"
	case Laboratory:
		return "LB"
	case GeneralStores:
		return "GS"
	default:
		return fmt.Sprintf("%#v", f)
	}
}

const (
	InventoryStingrayLauncher     = 0
	InventoryAvalancheLauncher    = 1
	InventoryCannon               = 2
	InventoryFusionBallLauncher   = 3
	InventoryLaserCannon          = 4
	InventoryPlasmaBeam           = 5
	InventoryStingrayMissile      = 6
	InventoryAvalancheMissile     = 7
	InventoryCannonRounds         = 8
	InventoryFusionBalls          = 9
	InventoryTankCannon           = 10
	InventoryTankRocketLauncher   = 11
	InventoryTankLaserCannon      = 12
	InventoryHovertankPlasma      = 13
	InventoryHovertankLauncher    = 14
	InventoryPistol               = 15
	InventoryPistolClip           = 16
	InventoryRifle                = 17
	InventoryRifleClip            = 18
	InventoryHeavyCannon          = 19
	InventoryHCAPAmmo             = 20
	InventoryHCHEAmmo             = 21
	InventoryHCINAmmo             = 22
	InventoryAutoCannon           = 23
	InventoryACAPAmmo             = 24
	InventoryACHEAmmo             = 25
	InventoryACINAmmo             = 26
	InventoryRocketLauncher       = 27
	InventorySmallRocket          = 28
	InventoryLargeRocket          = 29
	InventoryIncendiaryRocket     = 30
	InventoryLaserPistol          = 31
	InventoryLaserRifle           = 32
	InventoryHeavyLaser           = 33
	InventoryGrenade              = 34
	InventorySmokeGrenade         = 35
	InventoryProximityGrenade     = 36
	InventoryHighExplosive        = 37
	InventoryMotionScanner        = 38
	InventoryMediKit              = 39
	InventoryPsiAmp               = 40
	InventoryStunRod              = 41
	InventoryElectroFlare         = 42
	InventoryCORPSE               = 46
	InventoryCORPSE_ARMOUR        = 47
	InventoryCORPSE_POWERSUIT     = 48
	InventoryHeavyPlasma          = 49
	InventoryHeavyPlasmaClip      = 50
	InventoryPlasmaRifle          = 51
	InventoryPlasmaRifleClip      = 52
	InventoryPlasmaPistol         = 53
	InventoryPlasmaPistolClip     = 54
	InventoryBlasterLauncher      = 55
	InventoryBlasterBomb          = 56
	InventorySmallLauncher        = 57
	InventoryStunBomb             = 58
	InventoryAlienGrenade         = 59
	InventoryElerium115           = 60
	InventoryMindProbe            = 61
	InventorySectoidCorpse        = 65
	InventorySnakemanCorpse       = 66
	InventoryEtherealCorpse       = 67
	InventoryMutonCorpse          = 68
	InventoryFloaterCorpse        = 69
	InventoryCelatidCorpse        = 70
	InventorySilacoidCorpse       = 71
	InventoryChryssalidCorpse     = 72
	InventoryReaperCorpse         = 73
	InventorySectopodCorpse       = 74
	InventoryCyberdiscCorpse      = 75
	InventoryHovertankCorpse      = 76
	InventoryTankCorpse           = 77
	InventoryMaleCivilianCorpse   = 78
	InventoryFemaleCivilianCorpse = 79
	InventoryUFOPowerSource       = 80
	InventoryUFONavigation        = 81
	InventoryUFOConstruction      = 82
	InventoryAlienFood            = 83
	InventoryAlienReproduction    = 84
	InventoryAlienEntertainment   = 85
	InventoryAlienSurgery         = 86
	InventoryExaminationRoom      = 87
	InventoryAlienAlloys          = 88
	InventoryAlienHabitat         = 89
	InventoryPersonalArmour       = 90
	InventoryPowerSuit            = 91
	InventoryFlyingSuit           = 92
	InventoryHWPCannonShell       = 93
	InventoryHWPRockets           = 94
	InventoryHWPFusionBomb        = 95
)
