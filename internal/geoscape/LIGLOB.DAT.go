package geoscape

// This file is used by GEOSCAPE.EXE and it's structure is very simple. Every record is a 4 byte signed
// long integer. Probably the most useful offset is the first 4 bytes where your current money is stored.
// The rest of the bytes are used for the Finance graphs: Expenditure, Maintenance, and Balance (the
// others are stored elsewhere).
// https://www.ufopaedia.org/index.php/LIGLOB.DAT

type LIGLOB_DAT struct {
	// Current balance
	CurrentBalance int32 `struct:"int32"`
	// Expenditure for the last 12 months
	Expenditure []int32 `struct:"[12]int32"`
	// Maintenance costs for the last 12 months
	Maintenance []int32 `struct:"[12]int32"`
	// Balance for the last 12 months
	Balance []int32 `struct:"[12]int32"`
}
