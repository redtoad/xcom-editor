package geoscape

// ResearchFile holds all research topics.
//
// The file is 2112 bytes long (96 topics * 22 bytes). A new game copies the
// table from GEOSCAPE.EXE (offset 371496); only Unknown1, Completed and
// PrerequisitesMet change during the game.
//
// https://www.ufopaedia.org/index.php/RESEARCH.DAT
type ResearchFile struct {
	Topics [NumResearchTopics]ResearchTopic
}

// ResearchTopic holds the data of a single research topic.
type ResearchTopic struct {
	// 0-1: Base research time in scientist days. It is randomized by +/-50%
	// when the research starts and then copied to PROJECT.DAT.
	Time uint16
	// 2-4: Topic code. For live aliens, byte 3 is the AlienRace and byte 4 the
	// AlienRank. If byte 2 or 3 is not 0, the item must be in stores before
	// the topic can be researched.
	Topic [3]uint8
	// 5: Unknown, 0 in the initial table. For live alien topics this is the
	// index in ASTORE.DAT of the first matching alien (same race and rank)
	// in the current base. GEOSCAPE.EXE stores it while building the
	// research list of a base (0x19fd3) and does not update it afterwards,
	// so it may be stale (e.g. after the alien was researched). No code
	// reading it has been found.
	Unknown1 uint8
	// 6: UFOpaedia category.
	UfopaediaCategory uint8
	// 7: UFOpaedia entry within the category.
	UfopaediaEntry uint8
	// 8: Score points awarded on completion.
	Points uint8
	// 9: Unknown, always 0. Not accessed by any code.
	Unknown2 uint8
	// 10-13: Non-zero if the research is completed.
	Completed int32
	// 14-17: Prerequisite flag (exact meaning unclear). 0 or 1, never changes
	// during the game.
	Prerequisite int32
	// 18-21: Non-zero if all prerequisites are met.
	PrerequisitesMet int32
}
