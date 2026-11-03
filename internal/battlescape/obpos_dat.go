package battlescape

// NoOwner is the owner of an object lying on the ground.
const NoOwner = 0xFF

// ObposFile holds all objects on the map that can be picked up (e.g. weapons,
// ammunition and corpses), whether on the ground or carried by units.
//
// The file is 2720 bytes long (170 objects * 16 bytes).
//
// https://www.ufopaedia.org/index.php/OBPOS.DAT
type ObposFile struct {
	Objects [170]BattleObject
}

// BattleObject is a single object on the map.
type BattleObject struct {
	// 0: Object type (index in OBDATA.DAT), 0xFF if unused or destroyed.
	Type uint8
	// 1-3: Position on the ground (Z is not updated while carried).
	Y uint8
	X uint8
	Z uint8
	// 4: Index of the unit carrying the object (UNITPOS.DAT), NoOwner if it
	// is on the ground.
	Owner uint8
	// 5: Inventory slot: 0 = right hand, 1 = left hand, 2-3 = right leg,
	// 4-5 = left leg, 6-7 = right shoulder, 8-9 = left shoulder,
	// 10-18 = backpack, 19-24 = belt.
	Slot uint8
	// 6-7: Index of the object this object is loaded into (e.g. a clip in a
	// rifle), 0xFFFF if none.
	Container uint16
	// 8: Ammunition, painkillers, primed turn of grenades or UNITPOS.DAT
	// index of an unconscious unit (depending on the object type).
	Value1 uint8
	// 9: Stimulants, last thrower of grenades or unit type of corpses.
	Value2 uint8
	// 10: Heals of medikits or rank of corpses.
	Value3 uint8
	// 11: Copy of the MCD floor offset of the tile.
	Floor int8
	// 12: Flags (e.g. 8 = unconscious unit).
	Flags uint8
	// 13-15: Unknown. Byte 13 is always 0, bytes 14-15 are usually 0 but
	// some entries contain values like 0x0004, 0x0008 or 0x0288, 0x028c
	// (probably stale offsets or pointers).
	Unknown [3]uint8
}
