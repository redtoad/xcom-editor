package geoscape

// ProductFile holds the manufacturing data of all manufacturable items.
//
// The file is 630 bytes long (35 items * 18 bytes).
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
	// 15: Unknown, usually 1.
	Unknown1 uint8
	// 16: 0 if the item can be built, 1 if it still needs to be researched.
	NotResearched uint8
	// 17: Unknown, usually 0.
	Unknown2 uint8
}
