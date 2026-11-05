package geoscape

// ProductFile holds the manufacturing data of all manufacturable items.
//
// The file is 630 bytes long (35 items * 18 bytes). A new game copies the
// table from GEOSCAPE.EXE (offset 373608); only NotResearched changes during
// the game, all unknown fields are constants.
//
// https://www.ufopaedia.org/index.php/PRODUCT.DAT
type ProductFile struct {
	Products [NumProducts]Product
}

// Product holds the manufacturing data of a single item.
type Product struct {
	// 0-3: Manufacturing cost in $.
	Cost int32
	// 4-5: Engineer hours needed to build one item.
	Hours uint16
	// 6-7: Workshop space needed.
	WorkSpace uint16
	// 8-10: Required items (Inventory index in BASE.DAT), 0xFF if none.
	Requirements [3]uint8
	// 11-13: Quantity needed of each required item.
	RequiredQuantities [3]uint8
	// 14: Manufactured item (Inventory index in BASE.DAT).
	Item uint8
	// 15: Unknown, always 1.
	Unknown1 uint8
	// 16: 0 if the item can be built, 1 if it still needs to be researched.
	NotResearched uint8
	// 17: Unknown, always 0.
	Unknown2 uint8
}

// ProductNames holds the English names of all manufacturable items in the
// order of PRODUCT.DAT.
var ProductNames = [NumProducts]string{
	"Fusion Ball Launcher", "Laser Cannon", "Plasma Beam", "Fusion Ball",
	"Tank/Laser Cannon", "Hovertank/Plasma", "Hovertank/Launcher", "HWP Fusion Bomb",
	"Laser Pistol", "Laser Rifle", "Heavy Laser", "Motion Scanner",
	"Medi-Kit", "Psi-Amp", "Heavy Plasma", "Heavy Plasma Clip",
	"Plasma Rifle", "Plasma Rifle Clip", "Plasma Pistol", "Plasma Pistol Clip",
	"Blaster Launcher", "Blaster Bomb", "Small Launcher", "Stun Bomb",
	"Alien Grenade", "Mind Probe", "Personal Armour", "Power Suit",
	"Flying Suit", "Alien Alloys", "UFO Power Source", "UFO Navigation",
	"Firestorm", "Lightning", "Avenger",
}
