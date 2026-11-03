package geoscape

// NumActs is the number of alien mission types (see MissionType) in
// ACTS.DAT and MISSIONS.DAT. Alien supply missions are not included.
const NumActs = 7

// MissionInactive is the UFO counter of an inactive mission.
const MissionInactive = 0xFFFF

// MissionsFile holds the ongoing alien missions, indexed by MissionZone and
// MissionType.
//
// The file is 672 bytes long (12 zones * 7 mission types * 8 bytes).
//
// https://www.ufopaedia.org/index.php/MISSIONS.DAT
type MissionsFile struct {
	Zones [NumZones][NumActs]AlienMission
}

// AlienMission is a single ongoing alien mission.
type AlienMission struct {
	// 0-1: Number of UFOs sent so far, MissionInactive if the mission is
	// not active.
	UFOCounter uint16
	// 2-3: Secondary UFO counter, only used by some missions.
	SubCounter uint16
	// 4-5: Time until the next UFO is sent (counts down every 30 minutes).
	Timer uint16
	// 6-7: Alien race (same values as CraftData.AlienRace). Undefined for
	// inactive missions.
	Race uint16
}
