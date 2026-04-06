package geoscape

//go:generate stringer -type=Difficulty -output=iglob_dat_string.go -linecomment

import "encoding/binary"

// IglobFile stores the current date, time, and game state counters.
// The file is 60 bytes in the original DOS EU version, or 64 bytes in
// patched versions which also include a difficulty setting.
//
// Note that changing the days, minutes, hours or seconds do not change
// the values seen in the Load/Save game screen, but do change the values
// when you actually load the game. The values seen in the Load/Save
// screen is read from SAVEINFO.DAT.
//
// https://www.ufopaedia.org/index.php/IGLOB.DAT
type IglobFile struct {
	// 0: Current month (1-12).
	Month int32
	// 4: Current day of week (0-6).
	Weekday int32
	// 8: Current day of month (1-31).
	Day int32
	// 12: Current hour (0-23).
	Hour int32
	// 16: Current minute (0-59).
	Minute int32
	// 20: Current second (0-55, multiples of 5).
	Second int32
	// 24: Number of craft in hidden status.
	HiddenCraft int32
	// 28: Number of airborne X-COM craft.
	AirborneXCraft int32
	// 32: Number of craft refueling.
	RefuelingCraft int32
	// 36: Number of damaged craft.
	DamagedCraft int32
	// 40: Number of craft being repaired.
	RepairingCraft int32
	// 44: Total number of airborne craft.
	AirborneCraft int32
	// 48: Number of UFOs on retaliation missions.
	UFOsOnRetaliation int32
	// 52: Number of alien bases.
	AlienBases int32
	// 56: Current MIDI track number.
	CurrentMIDI int32
	// 60: (optional) Difficulty setting. nil if the file is 60 bytes (DOS EU version).
	// If present, the file is 64 bytes (patched version).
	Difficulty *Difficulty
}

// Difficulty is the game difficulty level.
type Difficulty int32

const (
	DifficultyBeginner    Difficulty = 0 // Beginner
	DifficultyExperienced Difficulty = 1 // Experienced
	DifficultyVeteran     Difficulty = 2 // Veteran
	DifficultyGenius      Difficulty = 3 // Genius
	DifficultySuperhuman  Difficulty = 4 // Superhuman
)

// SizeOf implements restruct.Sizer.
func (f IglobFile) SizeOf() int {
	if f.Difficulty != nil {
		return 64
	}
	return 60
}

// Unpack implements restruct.Unpacker.
func (f *IglobFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	f.Month = int32(order.Uint32(buf[0:]))
	f.Weekday = int32(order.Uint32(buf[4:]))
	f.Day = int32(order.Uint32(buf[8:]))
	f.Hour = int32(order.Uint32(buf[12:]))
	f.Minute = int32(order.Uint32(buf[16:]))
	f.Second = int32(order.Uint32(buf[20:]))
	f.HiddenCraft = int32(order.Uint32(buf[24:]))
	f.AirborneXCraft = int32(order.Uint32(buf[28:]))
	f.RefuelingCraft = int32(order.Uint32(buf[32:]))
	f.DamagedCraft = int32(order.Uint32(buf[36:]))
	f.RepairingCraft = int32(order.Uint32(buf[40:]))
	f.AirborneCraft = int32(order.Uint32(buf[44:]))
	f.UFOsOnRetaliation = int32(order.Uint32(buf[48:]))
	f.AlienBases = int32(order.Uint32(buf[52:]))
	f.CurrentMIDI = int32(order.Uint32(buf[56:]))
	if len(buf) >= 64 {
		d := Difficulty(int32(order.Uint32(buf[60:])))
		f.Difficulty = &d
	}
	return buf[f.SizeOf():], nil
}

// Pack implements restruct.Packer.
func (f IglobFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	order.PutUint32(buf[0:], uint32(f.Month))
	order.PutUint32(buf[4:], uint32(f.Weekday))
	order.PutUint32(buf[8:], uint32(f.Day))
	order.PutUint32(buf[12:], uint32(f.Hour))
	order.PutUint32(buf[16:], uint32(f.Minute))
	order.PutUint32(buf[20:], uint32(f.Second))
	order.PutUint32(buf[24:], uint32(f.HiddenCraft))
	order.PutUint32(buf[28:], uint32(f.AirborneXCraft))
	order.PutUint32(buf[32:], uint32(f.RefuelingCraft))
	order.PutUint32(buf[36:], uint32(f.DamagedCraft))
	order.PutUint32(buf[40:], uint32(f.RepairingCraft))
	order.PutUint32(buf[44:], uint32(f.AirborneCraft))
	order.PutUint32(buf[48:], uint32(f.UFOsOnRetaliation))
	order.PutUint32(buf[52:], uint32(f.AlienBases))
	order.PutUint32(buf[56:], uint32(f.CurrentMIDI))
	if f.Difficulty != nil {
		order.PutUint32(buf[60:], uint32(*f.Difficulty))
	}
	return buf, nil
}
