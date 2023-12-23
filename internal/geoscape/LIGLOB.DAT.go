package geoscape

// This file is used by GEOSCAPE.EXE and it's structure is very simple. Every record is a 4 byte signed
// long integer. Probably the most useful offset is the first 4 bytes where your current money is stored.
// The rest of the bytes are used for the Finance graphs: Expenditure, Maintenance, and Balance (the
// others are stored elsewhere).
// https://www.ufopaedia.org/index.php/LIGLOB.DAT

type LIGLOB_DAT struct {
	CurrentBalance int64   `struct:"int64"`
	Expenditure    []int64 `struct:"[12]int64"`
	Maintenance    []int64 `struct:"[12]int64"`
	Balance        []int64 `struct:"[12]int64"`
}
