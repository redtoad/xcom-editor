package savegame

import (
	"encoding/binary"
	"fmt"

	"github.com/go-restruct/restruct"
)

// https://www.ufopaedia.org/index.php/SOLDIER.DAT

type FileSoldier struct {
	Soldiers []Soldier
}

func (s FileSoldier) SizeOf() int {
	return 250 * 68
}

func (s FileSoldier) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	for i := 0; i < 250; i++ {
		data, err := restruct.Pack(order, &s.Soldiers[i])
		if err != nil {
			return nil, err
		}
		offset := i * 68
		for j := 0; j < len(data); j++ {
			buf[offset+j] = data[j]
		}
	}
	return buf, nil
}

func (s *FileSoldier) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	s.Soldiers = make([]Soldier, 250)
	for i := 0; i < 250; i++ {
		offset := i * 68
		data := buf[offset : offset+68]
		if err := restruct.Unpack(data, order, &s.Soldiers[i]); err != nil {
			return nil, err
		}
	}
	return buf[s.SizeOf():], nil
}

type Soldier struct {

	// Basic Information

	// 0-1 / 00-01 (0-5, FFFF): Integer set to Rank or FFFF if soldier is dead (or slot
	// not yet used). Values are actually a pointer within ENGLISH.DAT. Ranks are:
	//
	// 0 Rookie
	// 1 Squaddie
	// 2 Sergeant (SGT)
	// 3 Captain (CPT)
	// 4 Colonel (COL)
	// 5 Commander (CDR)
	// 6 Select Squad for
	// 7 SPACE AVAILABLE>

	/* byte=8bits word=2bytes int=4bytes long=8bytes */

	Rank Rank `struct:"int16"`

	// 2-3 / 02-03 (various, FF): Base the soldier is at, using reference values
	// from LOC.DAT. First base is always 0 (first entry in LOC.DAT). Value of FF
	// means the soldier is being transferred.
	Base int `struct:"int16"`

	// 4-5 / 04-05 (various, FFFF): Craft soldier is assigned to, using values from
	// LOC.DAT. FFFF means not assigned to a craft.
	Craft int `struct:"int16"`

	// 6-7 / 06-07 (various, FFFF): Craft soldier was assigned to before being
	// wounded and spending time in the infirmary.
	CraftBefore int `struct:"int16"`

	// 8-9 / 08-09 (0+): Missions. Signed integer. (The soldier display shows a
	// negative value if highest/16th bit is set!)
	Missions int `struct:"int16"`

	// 10-11 / 0A-0B (0+): Kills. Signed integer.
	Kills int `struct:"int16"`

	// 12-13 / 0C-0D (0+): Wound recovery days. Signed integer.
	RecoveryDays int `struct:"int16"`

	// 14-15 / 0E-0F (20+): Soldier value a.k.a. victory point loss if dies on a
	// mission. Signed integer. Is equal to 20 + Missions + Rank Bonus, as follows:

	//    0  Recruit
	//   +1  Squaddie
	//   +1  SGT
	//   +3  CPT
	//   +6  COL
	//  +10  CDR

	// Note: Hex-editing Missions or Rank and performing a combat (while not
	// touching this field) does not cause this field to reflect the edited values.
	// This field simply adds 1 for each mission and a delta when promoted (e.g.
	// COL to CDR +4). In other words, although it uses the equation shown above,
	// it's not actually computed from Missions and Rank, per se.
	SoldierValue int `struct:"int16"`

	// 16-40 / 10-28 (Text): Name, 25 characters long. Only approx. 21 characters
	// can be entered when editing the field within the game's Soldier display, but
	// longer names can be hex edited and will be seen in the game. Almost any
	// keyboard character can be entered (including within the game) because they are
	// stored directly instead of going through Windows objects. The end of the current
	// name is a null byte; garbage can be present after that (ends of longer previous
	// names, etc.).
	Name string `struct:"[25]byte"`

	// 41 / 29 (various): Always 0 except for existing soldiers being transferred, in
	// which case it equals the LOC.DAT value for destination base. New recruits on
	// their way to a base will also have this set to 0, so if Byte 3 is set to FF
	// and this byte is 0, you can't tell (from these two bytes alone) if it's an
	// existing soldier going to your first base (LOC.DAT=0), or a new recruit going
	// to any base. But the game knows where they're going, via TRANSFER.DAT.
	DestinationBase int `struct:"int8"`

	// Initial Recruit Stats

	// 42 / 2A (50-60): Initial Time Units
	InitialTimeUnits int `struct:"int8"`

	// 43 / 2B (25-40): Initial Health
	InitialHealth int `struct:"int8"`

	// 44 / 2C (40-70): Initial Energy a.k.a. Stamina
	InitialEnergy int `struct:"int8"`

	// 45 / 2D (30-60): Initial Reactions
	InitialReactions int `struct:"int8"`

	// 46 / 2E (20-40): Initial Strength
	InitialStrength int `struct:"int8"`

	// 47 / 2F (40-70): Initial Firing Accuracy
	InitialFiringAccuracy int `struct:"int8"`

	// 48 / 30 (50-80): Initial Throwing Accuracy
	InitialThrowingAccuracy int `struct:"int8"`

	// 49 / 31 (20-40): Initial Melee Accuracy a.k.a. Close Combat Accuracy
	InitialMeleeAccuracy int `struct:"int8"`

	// 50 / 32 (0-100): Psionic Strength. Never changes.
	InitialPsionicStrength int `struct:"int8"`

	// 51 / 33 (0-105+): Current Psionic Skill. 0 = Not yet trained / Psi stats not visible. See Note 1.
	// Unlike other stats, Psi Skill always starts at 0. So it stores its current skill in what would've otherwise been its initial skill byte, if you will. Also see Bytes 64 and 65.
	InitialPsionicSkill int `struct:"int8"`

	// 52 / 34 (5-10): Initial Bravery. Computed as 110-(10*ThisByte). E.g. if ThisByte=9, Bravery=20. So the lower this byte is, the better your Bravery is (see table at right). Initial Bravery can be 10-60, and it only uses increments of 10.
	// [52] Value	Bravery
	// 0	110
	// 1	100
	// 2	90
	// 3	80
	// 4	70
	// 5	60
	// 6	50
	// 7	40
	// 8	30
	// 9	20
	// 10	10
	// 11	0
	InitialBravery int `struct:"int8"`

	// Stat Improvement
	// Add to initial values to get current total.

	// 53 / 35 (0-31): Time Unit Improvement. (Max total TUs possible (initial+improvement) is 81. See Note 1.)
	TimeUnitImprovement int `struct:"int8"`

	// 54 / 36 (0-36): Health Improvement. (Max total 61)
	HealthImprovement int `struct:"int8"`

	// 55 / 37 (0-61): Stamina Improvement. (Max total 101)
	EnergyImprovement int `struct:"int8"`

	// 56 / 38 (0-75): Reaction Improvement. (Max total 105)
	ReactionsImprovement int `struct:"int8"`

	// 57 / 39 (0-51): Strength Improvement. (Max total 71) Note that very high (hacked)
	// values can cause "unable to throw here", presumably because the arc would intercept
	// the "ceiling" of the Battlescape. See throwing distance for more on Strength versus
	// Throwing.
	StrengthImprovement int `struct:"int8"`

	// 58 / 3A (0-85): Firing Accuracy Improvement. (Max total 125)
	FiringAccuracyImprovement int `struct:"int8"`

	// 59 / 3B (0-75): Throwing Accuracy Improvement Improvement. (Max total 125)
	ThrowingAccuracyImprovement int `struct:"int8"`

	// 60 / 3C (0-105): Melee/Close Combat Accuracy Improvement. (Max total 125)
	MeleeAccuracyImprovement int `struct:"int8"`

	// 61 / 3D (0-9): Bravery Improvement (multiply by 10). (Max total 9; 100 Bravery)
	BraveryImprovement int `struct:"int8"`

	// Odds and Ends

	// 62 / 3E (0-3): Armor. 0=None, 1=Personal Armor, 2=Power Suit, 3=Flying Suit.
	// Can be hacked to: 4=Sectoid, 5=Snakeman, 6=Ethereal, 7=Muton, 8=Floater, 9=Celatid, 10=Silacoid, 11=Chryssalid. (Higher values are civilians and parts of big units.) NOTE: These are not real aliens, just skins... more research must be done e.g. on appearance and gender to make real aliens. However, when you start a new mission, the CE version crashes because it cannot find a skin for that value.
	Armor Armor `struct:"int8"`

	// 63 / 3F (0-24): Most-recent month's Psi Lab training increase. Largest values for newly trained (16-24). For Psi Lab point-award functionality, see Psionic Skill.
	MostRecentPsiLabTraining int `struct:"int8"`

	// 64 / 40 (0,1): In Psi Lab training.
	InPsiLabTraining bool `struct:"int8"`

	// 65 / 41 (0,1): Promotion flag (as of most recent combat). Reset as of next combat.
	Promotion bool `struct:"int8"`

	// 66 / 42 (0,1): 0=Male, 1=Female.
	// Editing values greater than 1 does not change the appearance in the soldier
	// equip screens. However, the soldier image will always take the outward appearance
	// of a female during a mission.
	Sex Gender `struct:"int8"`

	// 67 / 43 (0-3): Appearance (for loadout, no armor): 0=blonde, 1=brown hair, 2=Oriental, 3=African.
	// Refers to the UFOGRAPH/MAN.spk images. Because there are only 4 nationalities, editing values greater
	// than 3 will result in an error code: cannot open file ufograph/man_XYZ.spk where X is 0 or 1 (Light armored or personnel armored), Y is "M" or "F" (Male or Female) and Z is the nationality value listed in this field.
	Appearance Appearance `struct:"int8"`
}

func (s *Soldier) String() string {
	if s.Rank == DeadOrUnused {
		return "{Soldier <dead or unused>}"
	} else {
		return fmt.Sprintf("{Soldier %v %s}", s.Rank, s.Name)
	}
}

type Rank int

func (r Rank) String() string {
	switch r {
	case DeadOrUnused:
		return "Dead"
	case Rookie:
		return "Rookie"
	case Squaddie:
		return "Squaddie"
	case Sergeant:
		return "Sergeant"
	case Captain:
		return "Captain"
	case Colonel:
		return "Colonel"
	case Commander:
		return "Commander"
	}
	return "???"
}

const (
	Rookie Rank = iota
	Squaddie
	Sergeant
	Captain
	Colonel
	Commander
	DeadOrUnused Rank = -1
)

type Armor int

const (
	NoArmor Armor = iota
	PersonalArmor
	PowerSuit
	FlyingSuit

	// Can be hacked to
	Sectoid
	Snakeman
	Ethereal
	Muton
	Floater
	Celatid
	Silacoid
	Chryssalid
)

type Gender int

const (
	Male Gender = iota
	Female
)

type Appearance int

const (
	Blonde Appearance = 0
	BrownHair
	Oriental
	African
)
