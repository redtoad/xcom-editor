package geoscape

// AlienFile contains alien activity scores for countries and regions, used to
// populate the graphs screen. The same format is used for XCOM.DAT (which
// records X-COM activity instead). The first record is the current month, the
// second is last month, and so on.
//
// Each score increments by 1 every 30 minutes that a craft is over that
// territory. For countries, the matching region is also incremented.
//
// https://www.ufopaedia.org/index.php/ALIEN.DAT
type AlienFile struct {
	Months [12]ActivityMonth
}

// ActivityMonth holds activity scores for all countries and regions for a
// single month. Values are signed and may go negative in XCOM.DAT.
type ActivityMonth struct {
	// Countries (offsets 0x00–0x3C)
	USA         int32 `struct:"int32"`
	Russia      int32 `struct:"int32"`
	UK          int32 `struct:"int32"`
	France      int32 `struct:"int32"`
	Germany     int32 `struct:"int32"`
	Italy       int32 `struct:"int32"`
	Spain       int32 `struct:"int32"`
	China       int32 `struct:"int32"`
	Japan       int32 `struct:"int32"`
	India       int32 `struct:"int32"`
	Brazil      int32 `struct:"int32"`
	Australia   int32 `struct:"int32"`
	Nigeria     int32 `struct:"int32"`
	SouthAfrica int32 `struct:"int32"`
	Egypt       int32 `struct:"int32"`
	Canada      int32 `struct:"int32"`
	// Regions (offsets 0x40–0x78)
	NorthAmerica   int32 `struct:"int32"`
	Arctic         int32 `struct:"int32"`
	Antarctica     int32 `struct:"int32"`
	SouthAmerica   int32 `struct:"int32"`
	Europe         int32 `struct:"int32"`
	NorthAfrica    int32 `struct:"int32"`
	SouthernAfrica int32 `struct:"int32"`
	CentralAsia    int32 `struct:"int32"`
	SoutheastAsia  int32 `struct:"int32"`
	Siberia        int32 `struct:"int32"`
	Australasia    int32 `struct:"int32"`
	Pacific        int32 `struct:"int32"`
	NorthAtlantic  int32 `struct:"int32"`
	SouthAtlantic  int32 `struct:"int32"`
	IndianOcean    int32 `struct:"int32"`
}
