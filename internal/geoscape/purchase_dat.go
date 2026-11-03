package geoscape

// PurchaseFile holds prices and transfer times of everything that can be
// bought or sold.
//
// The file is 1820 bytes long (91 entries * 20 bytes).
//
// https://www.ufopaedia.org/index.php/PURCHASE.DAT
type PurchaseFile struct {
	Items [91]PurchaseItem
}

// PurchaseItem holds the purchase data of a single item.
type PurchaseItem struct {
	// 0: Category: 0 = item, 1 = craft, 2 = soldier, 3 = scientist, 4 = engineer.
	Category uint8
	// 1: Index of the item within its category (e.g. Inventory index or craft type).
	Index uint8
	// 2-5: Unknown, usually 1000.
	Unknown1 int32
	// 6-9: Buy price in $.
	BuyPrice int32
	// 10-13: Sell price in $.
	SellPrice int32
	// 14: Transfer time in hours.
	TransferHours uint8
	// 15: Unknown, always 0.
	Unknown2 uint8
	// 16: Quantity last purchased.
	LastPurchased uint8
	// 17: Unknown, always 0.
	Unknown3 uint8
	// 18: 1 if the item can be bought, 0 otherwise.
	Purchasable uint8
	// 19: Unknown, always 0.
	Unknown4 uint8
}
