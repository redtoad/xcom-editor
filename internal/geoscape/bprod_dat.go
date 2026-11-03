package geoscape

// NumProducts is the number of manufacturable items (see PRODUCT.DAT).
const NumProducts = 35

// BprodFile holds the manufacturing progress of each base.
//
// The file is 2800 bytes long (8 bases * 350 bytes).
//
// https://www.ufopaedia.org/index.php/BPROD.DAT
type BprodFile struct {
	Bases [8]BaseProduction
}

// BaseProduction holds the manufacturing progress of a single base. All
// arrays are indexed in the same order as PRODUCT.DAT. Completed projects
// are reset to 0 in all fields.
type BaseProduction struct {
	// 0-139: Hours of work done towards the next unit. An item is only shown
	// as in production if this value is not 0.
	HoursStored [NumProducts]int32
	// 140-209: Engineers allocated.
	Engineers [NumProducts]int16
	// 210-279: Total number of units to produce.
	Total [NumProducts]int16
	// 280-349: Number of units produced.
	Produced [NumProducts]int16
}
