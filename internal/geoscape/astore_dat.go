package geoscape

//go:generate stringer -type=AlienRace,AlienRank -output=astore_dat_string.go -linecomment

// AstoreFile holds the live aliens kept in alien containment.
//
// The file is 600 bytes long (50 entries * 12 bytes).
//
// https://www.ufopaedia.org/index.php/ASTORE.DAT
type AstoreFile struct {
	Aliens [50]StoredAlien
}

// StoredAlien is a single live alien in alien containment.
type StoredAlien struct {
	// 0: Alien race. AlienNone if the entry is unused (e.g. after the alien
	// was researched).
	Race AlienRace
	// 1: Alien rank.
	Rank AlienRank
	// 2: Base the alien is kept in (index in LOC.DAT). 0xFF while in transfer.
	Base uint8
	// 3-10: Unknown, usually 0. Neither GEOSCAPE.EXE nor TACTICAL.EXE (which
	// only copies the file) reads or writes these bytes: a new game only
	// clears bytes 0 and 11 of each entry (GEOSCAPE.EXE 0x21db0) and
	// storing an alien only sets bytes 0, 1, 2 and 11. The bytes therefore
	// keep uninitialized memory, e.g. 00 a4 3b 20 00 00 00 00 in entry 0.
	Unknown [8]uint8
	// 11: Transfer state: 1 = in transit, 2 = selected for transfer.
	TransferState uint8
}

// AlienRace is the race of a live alien in ASTORE.DAT.
type AlienRace uint8

const (
	AlienNone       AlienRace = iota // none
	AlienSectoid                     // Sectoid
	AlienSnakeman                    // Snakeman
	AlienEthereal                    // Ethereal
	AlienMuton                       // Muton
	AlienFloater                     // Floater
	AlienCelatid                     // Celatid
	AlienSilacoid                    // Silacoid
	AlienChryssalid                  // Chryssalid
	AlienReaper                      // Reaper
	AlienSectopod                    // Sectopod
	AlienCyberdisc                   // Cyberdisc
)

// AlienRank is the rank of a live alien in ASTORE.DAT.
type AlienRank uint8

const (
	AlienCommander AlienRank = iota + 1 // Commander
	AlienLeader                         // Leader
	AlienEngineer                       // Engineer
	AlienMedic                          // Medic
	AlienNavigator                      // Navigator
	AlienSoldier                        // Soldier
	AlienTerrorist                      // Terrorist
)
