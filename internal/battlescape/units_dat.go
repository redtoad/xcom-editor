package battlescape

import "github.com/redtoad/xcom-editor/internal"

// NumUnits is the number of unit slots in UNITPOS.DAT and UNITREF.DAT.
const NumUnits = 80

// Owners in UnitPosition.Owner.
const (
	OwnerXcom     = 0
	OwnerAlien    = 1
	OwnerCivilian = 2
)

// Flags in UnitPosition.Flags.
const (
	UnitVisible        = 1 << 0
	UnitActive         = 1 << 1
	UnitLarge          = 1 << 2
	UnitRightSection   = 1 << 3
	UnitLowerSection   = 1 << 4
	UnitMindControlled = 1 << 6
	UnitDamaged        = 1 << 7
)

// UnitposFile holds the positions and owners of all units on the map. Each
// entry belongs to the entry with the same index in UNITREF.DAT.
//
// The file is 1120 bytes long (80 units * 14 bytes).
//
// https://www.ufopaedia.org/index.php/UNITPOS.DAT
type UnitposFile struct {
	Units [NumUnits]UnitPosition
}

// UnitPosition holds the position of a single unit.
type UnitPosition struct {
	// 0-2: Position, Y is 0xFF if the slot is unused.
	Y uint8
	X uint8
	Z uint8
	// 3: Unknown.
	Unknown1 uint8
	// 4-7: Memory address of the UNITREF.DAT entry (not used on disk).
	UnitRef uint32
	// 8: Turns since the unit was last spotted.
	TurnsSinceSpotted uint8
	// 9: Owner: OwnerXcom, OwnerAlien, OwnerCivilian or 3 for X-COM units
	// left behind on a previous stage.
	Owner uint8
	// 10: Flags (see UnitVisible etc.).
	Flags uint8
	// 11-13: Unknown.
	Unknown2 [3]uint8
}

// UnitrefFile holds the stats of all units on the map. Each entry belongs to
// the entry with the same index in UNITPOS.DAT.
//
// The file is 9920 bytes long (80 units * 124 bytes).
//
// https://www.ufopaedia.org/index.php/UNITREF.DAT
type UnitrefFile struct {
	Units [NumUnits]Unit
}

// Unit holds the stats of a single unit.
type Unit struct {
	// 0: Unit type, 0xFF if the slot is unused.
	Type uint8
	// 1: Inventory paper doll image.
	PaperDoll uint8
	// 2-9: Memory addresses of the unit's sprites (reset on load).
	PCKPointer uint32
	TABPointer uint32
	// 10: Facing (0-7).
	Facing uint8
	// 11: Movement type: 0 = walk, 1 = slide, 2 = fly.
	MovementType uint8
	// 12-17: Current stats.
	TimeUnits uint8
	Health    uint8
	Stun      uint8
	Energy    uint8
	Reactions uint8
	Strength  uint8
	// 18-22: Current armor (front, left, right, rear, under).
	Armor [5]uint8
	// 23-24: Firing and throwing accuracy.
	FiringAccuracy   uint8
	ThrowingAccuracy uint8
	// 25-28: Maximum stats.
	MaxTimeUnits uint8
	MaxHealth    uint8
	MaxEnergy    uint8
	MaxReactions uint8
	// 29-33: Maximum armor (front, left, right, rear, under).
	MaxArmor [5]uint8
	// 34: Health recharge (unused).
	HealthRecharge uint8
	// 35: Energy recharge per turn.
	EnergyRecharge int8
	// 36: Armor recharge (unused).
	ArmorRecharge uint8
	// 37: Psionic skill.
	PsiSkill uint8
	// 38: Item created when the unit dies or falls unconscious (OBDATA.DAT).
	Corpse uint8
	// 39: Score value of the unit.
	Value uint8
	// 40: Index in SOLDIER.DAT, 0xFF if the unit is not a soldier.
	Soldier uint8
	// 41: Animation frame counter.
	AnimationFrame uint8
	// 42: Rank, 0xFF if dead or unused.
	Rank uint8
	// 43: Deployment: 0 = normal, 1 = civilian spawn.
	Deployment uint8
	// 44: Alien aggression.
	Aggression uint8
	// 45: Energy loss divisor.
	EnergyLossDivisor uint8
	// 46: Tank turret image or inventory layout.
	TurretImage uint8
	// 47: Copy of the MCD floor offset of the tile.
	Floor int8
	// 48: Line of fire template.
	LOFTemplate uint8
	// 49-51: Height when standing, kneeling and floating.
	StandingHeight uint8
	KneelingHeight uint8
	FloatingHeight uint8
	// 52: Line of fire template radius.
	LOFRadius uint8
	// 53-54: Image of the item in the left and right hand, 0xFF if none.
	LeftHandImage  uint8
	RightHandImage uint8
	// 55: Damage modifier category.
	DamageModifiers uint8
	// 56: Melee accuracy.
	MeleeAccuracy uint8
	// 57: Psionic strength.
	PsiStrength uint8
	// 58: Current morale.
	Morale uint8
	// 59: Bravery (110 - 10 * value).
	Bravery uint8
	// 60: Panic: 0 = none, 1 = freeze, 2 = run, 3 = berserk.
	Panic uint8
	// 61-62: Motion scanner.
	MotionScannerStep       uint8
	MotionScannerVisibility uint8
	// 63-68: Fatal wounds (head, torso, right arm, left arm, right leg,
	// left leg).
	FatalWounds [6]uint8
	// 69-70: Patrol start and end node (ROUTES.DAT).
	RouteStart uint8
	RouteEnd   uint8
	// 71: Light emitted by the unit.
	Lighting uint8
	// 72: Morale restored by painkillers.
	PainkillerMorale uint8
	// 73: Intelligence (2-8).
	Intelligence uint8
	// 74: AI mode, 0xFF to recalculate.
	AIMode uint8
	// 75: Light direction (unused).
	LightDirection uint8
	// 76-77: Missions.
	Missions int16
	// 78-79: Kills.
	Kills int16
	// 80-85: Experience (reactions, firing, melee, throwing, psionics,
	// bravery).
	Experience [6]uint8
	// 86-110: Name (null terminated).
	Name internal.NullString `struct:"[25]byte"`
	// 111: 1 while firing.
	BigFire uint8
	// 112: 1 if the AI is done with the unit this turn.
	Finished uint8
	// 113: 1 for tanks.
	Tank uint8
	// 114: Turns the unit is burning.
	BurningTurns uint8
	// 115: Gender: 0 = male, 1 = female.
	Gender uint8
	// 116: Appearance (0-3).
	Appearance uint8
	// 117-118: Tank turret weapon (0xFF if unused) and ammunition.
	TurretWeapon uint8
	TurretAmmo   uint8
	// 119: Unknown.
	Unknown uint8
	// 120-123: Flags: 1 = moving, 2 = can fly, 4 = flying, 8 = selected this
	// turn, 16 = selection disabled, 32 = right hand item selected,
	// 64 = kneeling, 128 = wears power or flying suit.
	Flags uint32
}
