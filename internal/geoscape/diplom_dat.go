package geoscape

// DiplomFile holds the relations to the 16 funding countries (in the same
// order as in ALIEN.DAT).
//
// The file is 576 bytes long (16 countries * 36 bytes).
//
// https://www.ufopaedia.org/index.php/DIPLOM.DAT
type DiplomFile struct {
	Countries [16]CountryRelation
}

// CountryRelation holds the relation to a single funding country.
type CountryRelation struct {
	// 0-1: Index of the country name in ENGLISH.DAT (starts at 602).
	NameIndex int16
	// 2-3: Satisfaction: 3 = happy, 2 = satisfied, 1 = unhappy, 0 = signed a
	// pact with the aliens (stops funding).
	Satisfaction int16
	// 4-27: Monthly funding in $1000 (January to December).
	Funding [12]int16
	// 28-29: Change in funding (in $1000) between last month and this month.
	FundingChange int16
	// 30-31: Non-zero if the country signs a pact with the aliens at the end
	// of the month.
	Pact int16
	// 32-33: Funding at the start of the game (in $1000). No effect after the
	// first save.
	StartingFunding int16
	// 34-35: Maximum funding (in $1000).
	FundingCap int16
}
