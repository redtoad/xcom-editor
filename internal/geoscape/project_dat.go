package geoscape

// NumResearchTopics is the number of research topics (see RESEARCH.DAT).
const NumResearchTopics = 96

// ProjectFile holds the research progress of each base.
//
// The file is 2304 bytes long (8 bases * 288 bytes).
//
// https://www.ufopaedia.org/index.php/PROJECT.DAT
type ProjectFile struct {
	Bases [8]BaseResearch
}

// BaseResearch holds the research progress of a single base. Both arrays are
// indexed in the same order as RESEARCH.DAT.
type BaseResearch struct {
	// 0-191: Scientist days remaining. 0 if the topic is not being researched.
	DaysRemaining [NumResearchTopics]uint16
	// 192-287: Scientists assigned.
	Scientists [NumResearchTopics]uint8
}
