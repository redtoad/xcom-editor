package battlescape

// GeodataFile describes how the map was put together from map modules.
// Maps consist of up to 6x6 modules of 10x10 tiles.
//
// The file is 92 bytes long.
//
// https://www.ufopaedia.org/index.php/GEODATA.DAT
type GeodataFile struct {
	// 0-1: Number of modules along the Y axis (map size / 10).
	ModulesY uint16
	// 2-3: Number of modules along the X axis (map size / 10).
	ModulesX uint16
	// 4-5: Number of map levels.
	SizeZ uint16
	// 6-7: X-COM craft type (ignored for base missions).
	XcomCraft uint16
	// 8-9: Alien craft type (ignored if not a UFO mission).
	AlienCraft uint16
	// 10-11: Terrain type.
	Terrain uint16
	// 12-47: Module IDs, west to east and north to south. 253-255 mark
	// modules covered by a 2x2 module.
	Modules [36]uint8
	// 48-83: 1 if the routes (RMP) of the module are loaded, 0 otherwise.
	LoadRoutes [36]uint8
	// 84: Crafts on the map: 0 = none, 1 = X-COM only, 2 = X-COM and alien.
	Crafts uint8
	// 85: Copy of XcomCraft.
	XcomCraftCopy uint8
	// 86: Y position of the X-COM craft.
	XcomCraftY uint8
	// 87: X position of the X-COM craft.
	XcomCraftX uint8
	// 88: Copy of AlienCraft.
	AlienCraftCopy uint8
	// 89: Y position of the alien craft.
	AlienCraftY uint8
	// 90: X position of the alien craft.
	AlienCraftX uint8
	// 91: Unknown, always 0.
	Unknown uint8
}
