package savegame

import "time"

// SAVEINFO_DAT contains name and game time of the saved game.
// This file is 40 bytes long, no separate entries.
type SAVEINFO_DAT struct {

	// 0-1	0x00-0x01	Ignore this if the file is not in the missdat folder. If 0, then this is a savegame made on the beginning of a new battlescape game. If 1, then check DIRECT.DAT to see where which save slot to load from.
	_ bool `struct:"int16"`

	// 2-27	0x02-0x1D	This a 26 byte null terminated string, which details the name of the save file. The name may be 25 characters long; the final byte is always of value 0. A 0 also marks the end of the save name, should it not use all 25 characters.
	Name string `struct:"[26]byte"`

	// 28-29	0x1C-0x1D	The current year.
	Year int `struct:"int16"`

	// 30-31	0x1E-0x1F	The current month. Note that 0 is January.
	Month int `struct:"int16"`

	// 32-33	0x20-0x21	The current day of the month.
	DayOfMonth int `struct:"int16"`

	// 34-35	0x22-0x23	The current hour (24 hour time).
	Hour int `struct:"int16"`

	// 36-37	0x24-0x25	The current minute.
	Minute int `struct:"int16"`

	// 38-39	0x26-0x27	0 for geoscape save, 1 for tactical save.
	TacticalSave bool `struct:"int16"`
}

// SizeOf imlements restruct.Sizer
func (s SAVEINFO_DAT) SizeOf() int {
	return 40
}

func (s SAVEINFO_DAT) Time() time.Time {
	return time.Date(s.Year, time.Month(s.Month+1), s.DayOfMonth, s.Hour, s.Minute, 0, 0, time.UTC)
}
