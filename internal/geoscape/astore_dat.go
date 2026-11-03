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
	// 3-10: Unknown, usually 0.
	Unknown [8]uint8
	// 11: Transfer state: 1 = in transit, 2 = selected for transfer.
	TransferState uint8
}

// AlienRace is the race of a live alien in ASTORE.DAT.
type AlienRace uint8

const (
	AlienNone       AlienRace = 0  // none
	AlienSectoid    AlienRace = 1  // Sectoid
	AlienSnakeman   AlienRace = 2  // Snakeman
	AlienEthereal   AlienRace = 3  // Ethereal
	AlienMuton      AlienRace = 4  // Muton
	AlienFloater    AlienRace = 5  // Floater
	AlienCelatid    AlienRace = 6  // Celatid
	AlienSilacoid   AlienRace = 7  // Silacoid
	AlienChryssalid AlienRace = 8  // Chryssalid
	AlienReaper     AlienRace = 9  // Reaper
	AlienSectopod   AlienRace = 10 // Sectopod
	AlienCyberdisc  AlienRace = 11 // Cyberdisc
)

// AlienRank is the rank of a live alien in ASTORE.DAT.
type AlienRank uint8

const (
	AlienCommander AlienRank = 1 // Commander
	AlienLeader    AlienRank = 2 // Leader
	AlienEngineer  AlienRank = 3 // Engineer
	AlienMedic     AlienRank = 4 // Medic
	AlienNavigator AlienRank = 5 // Navigator
	AlienSoldier   AlienRank = 6 // Soldier
	AlienTerrorist AlienRank = 7 // Terrorist
)
