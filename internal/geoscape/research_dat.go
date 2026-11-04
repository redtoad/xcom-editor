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

// ResearchTopicNames holds the English names of all research topics in the
// order of RESEARCH.DAT.
//
// https://www.ufopaedia.org/index.php/RESEARCH.DAT
var ResearchTopicNames = [NumResearchTopics]string{
	"Laser Weapons", "Motion Scanner", "Medi-Kit", "Psi-Amp",
	"Heavy Plasma", "Heavy Plasma Clip", "Plasma Rifle", "Plasma Rifle Clip",
	"Plasma Pistol", "Plasma Pistol Clip", "Blaster Launcher", "Blaster Bomb",
	"Small Launcher", "Stun Bomb", "Alien Grenade", "Elerium-115",
	"Mind Probe", "UFO Power Source", "UFO Navigation", "UFO Construction",
	"Alien Food", "Alien Reproduction", "Alien Entertainment", "Alien Surgery",
	"Examination Room", "Alien Alloys", "New Fighter Craft", "New Fighter-Transporter",
	"Ultimate Craft", "Laser Pistol", "Laser Rifle", "Heavy Laser",
	"Laser Cannon", "Plasma Cannon", "Fusion Missile", "Laser Defense",
	"Plasma Defense", "Fusion Defense", "Grav Shield", "Mind Shield",
	"Psionic Laboratory", "Hyper-Wave Decoder",
	// 42-44 are not used by the game
	"Tank/Laser Cannon", "Hovertank/Plasma", "Hovertank/Launcher",
	"Sectoid Corpse", "Snakeman Corpse", "Ethereal Corpse", "Muton Corpse",
	"Floater Corpse", "Celatid Corpse", "Silacoid Corpse", "Chryssalid Corpse",
	"Reaper Corpse", "Sectopod Corpse", "Cyberdisc Corpse",
	"Alien Origins", "The Martian Solution", "Cydonia or Bust",
	"Personal Armour", "Power Suit", "Flying Suit",
	"Sectoid Commander", "Sectoid Leader", "Sectoid Engineer", "Sectoid Medic", "Sectoid Navigator", "Sectoid Soldier",
	"Snakeman Commander", "Snakeman Leader", "Snakeman Engineer", "Snakeman Medic", "Snakeman Navigator", "Snakeman Soldier",
	"Ethereal Commander", "Ethereal Leader", "Ethereal Engineer", "Ethereal Medic", "Ethereal Navigator", "Ethereal Soldier",
	"Muton Commander", "Muton Leader", "Muton Engineer", "Muton Medic", "Muton Navigator", "Muton Soldier",
	"Floater Commander", "Floater Leader", "Floater Engineer", "Floater Medic", "Floater Navigator", "Floater Soldier",
	"Celatid Terrorist", "Silacoid Terrorist", "Chryssalid Terrorist", "Reaper Terrorist",
}
