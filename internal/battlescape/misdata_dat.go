package battlescape

// MisdataFile holds the mission data and the score of the mission.
//
// The file is 170 bytes long.
//
// https://www.ufopaedia.org/index.php/MISDATA.DAT
type MisdataFile struct {
	// 0-19: Recovered UFO components (only set if the mission was won): power
	// source, navigation, construction, alien food, reproduction,
	// entertainment, surgery, examination room, alloys and habitat.
	Recovered [10]int16
	// 20-21: Mission type.
	MissionType int16
	// 22-73: Score of the mission, in this order: aliens killed, alien
	// corpses, live aliens, artefacts, alien base control, civilians killed by
	// aliens, civilians killed by X-COM, civilians saved, X-COM soldiers
	// killed, soldiers retired, soldiers missing in action, tanks lost and
	// crafts lost.
	Scores [13]MissionScore
	// 74-75: Time of day (index into a table in the executable).
	TimeOfDay int16
	// 76: LOC.DAT index of the UFO crash site.
	UFOLocation uint8
	// 77: CRAFT.DAT index of the X-COM craft (LOC.DAT index for base defence).
	Craft uint8
	// 78: Result: 0 = aborted, 1 = won, 2 = lost, 3 = aborted without units
	// in the exit area.
	Result uint8
	// 79: 1 for UFO sites, 0 otherwise.
	UFOSite uint8
	// 80-151: Base defence: number of targets left in each base module.
	BaseModules [36]int16
	// 152-167: Equipment flags (1 = enabled).
	Equipment [16]uint8
	// 168-169: Language: 0 = English, 1 = German, 2 = French.
	Language int16
}

// MissionScore holds the score for one category of a mission.
type MissionScore struct {
	Count int16
	Score int16
}
