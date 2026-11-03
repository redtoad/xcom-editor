package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mustUnpack unpacks a DAT file from savegame GAME_3 in testdata/.
func mustUnpack(t *testing.T, name string, obj any) {
	t.Helper()
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_3", name), binary.LittleEndian, obj))
}

func TestAknowFile_Unpack(t *testing.T) {
	var f geoscape.AknowFile
	mustUnpack(t, "AKNOW.DAT", &f)
	assert.Equal(t, [495]byte{}, f.Data)
}

func TestAstoreFile_Unpack(t *testing.T) {
	var f geoscape.AstoreFile
	mustUnpack(t, "ASTORE.DAT", &f)

	assert.Equal(t, geoscape.AlienFloater, f.Aliens[0].Race)
	assert.Equal(t, geoscape.AlienNavigator, f.Aliens[0].Rank)
	assert.Equal(t, uint8(0), f.Aliens[0].Base)
	assert.Equal(t, geoscape.AlienSoldier, f.Aliens[1].Rank)
	assert.Equal(t, geoscape.AlienReaper, f.Aliens[2].Race)
	assert.Equal(t, "Terrorist", f.Aliens[2].Rank.String())
}

func TestBprodFile_Unpack(t *testing.T) {
	var f geoscape.BprodFile
	mustUnpack(t, "BPROD.DAT", &f)
	assert.Equal(t, [geoscape.NumProducts]int32{}, f.Bases[0].HoursStored)
}

func TestDiplomFile_Unpack(t *testing.T) {
	var f geoscape.DiplomFile
	mustUnpack(t, "DIPLOM.DAT", &f)

	usa := f.Countries[0]
	assert.Equal(t, int16(602), usa.NameIndex)
	assert.Equal(t, int16(0), usa.Satisfaction)
	assert.Equal(t, [12]int16{786, 653, 653, 653, 536, 0, 0, 0, 0, 0, 0, 786}, usa.Funding)
	assert.Equal(t, int16(-536), usa.FundingChange)
	assert.Equal(t, int16(0), usa.Pact)
	assert.Equal(t, int16(600), usa.StartingFunding)
	assert.Equal(t, int16(10000), usa.FundingCap)
	for i, c := range f.Countries {
		assert.Equal(t, int16(602+i), c.NameIndex)
	}
}

func TestFacilFile_Unpack(t *testing.T) {
	var f geoscape.FacilFile
	mustUnpack(t, "FACIL.DAT", &f)

	lift := f.Facilities[geoscape.AccessLift]
	assert.Equal(t, uint16(0x249), lift.NameIndex)
	assert.Equal(t, int16(300), lift.BuildCost)
	assert.Equal(t, uint8(1), lift.BuildDays)
	assert.Equal(t, uint8(4), lift.Upkeep)
	assert.Equal(t, int16(400), f.Facilities[1].BuildCost)
	assert.Equal(t, uint8(30), f.Facilities[1].Capacity)
}

func TestLeaseFile_Unpack(t *testing.T) {
	var f geoscape.LeaseFile
	mustUnpack(t, "LEASE.DAT", &f)

	assert.Equal(t, int16(0x180), f.RotationY)
	assert.Equal(t, int16(0x9f), f.RotationX)
	assert.Equal(t, int16(90), f.Zoom)
}

func TestMissionsFile_Unpack(t *testing.T) {
	var f geoscape.MissionsFile
	mustUnpack(t, "MISSIONS.DAT", &f)

	infiltration := f.Zones[geoscape.NorthAmerica][geoscape.MissionAlienInfiltration]
	assert.Equal(t, uint16(2), infiltration.UFOCounter)
	assert.Equal(t, uint16(344), infiltration.Timer)
	assert.Equal(t, uint16(geoscape.MissionInactive), f.Zones[geoscape.NorthAmerica][geoscape.MissionAlienResearch].UFOCounter)
	assert.Equal(t, uint16(geoscape.MissionInactive), f.Zones[geoscape.Arctic][geoscape.MissionAlienInfiltration].UFOCounter)
}

func TestProductFile_Unpack(t *testing.T) {
	var f geoscape.ProductFile
	mustUnpack(t, "PRODUCT.DAT", &f)

	assert.Equal(t, geoscape.Product{
		Cost: 226000, Hours: 500, WorkSpace: 8,
		Requirements: [3]uint8{60, 255, 255}, RequiredQuantities: [3]uint8{15, 0, 0},
		Item: 5, Unknown1: 1, NotResearched: 1,
	}, f.Products[2])
}

func TestProjectFile_Unpack(t *testing.T) {
	var f geoscape.ProjectFile
	mustUnpack(t, "PROJECT.DAT", &f)

	var days, scientists int
	for i := range f.Bases[0].DaysRemaining {
		days += int(f.Bases[0].DaysRemaining[i])
		scientists += int(f.Bases[0].Scientists[i])
	}
	assert.Equal(t, 69, days)
	assert.Equal(t, 50, scientists)
}

func TestPurchaseFile_Unpack(t *testing.T) {
	var f geoscape.PurchaseFile
	mustUnpack(t, "PURCHASE.DAT", &f)

	soldier := f.Items[0]
	assert.Equal(t, uint8(2), soldier.Category)
	assert.Equal(t, int32(40000), soldier.BuyPrice)
	assert.Equal(t, uint8(72), soldier.TransferHours)
	assert.Equal(t, uint8(1), soldier.Purchasable)
	assert.Equal(t, int32(115000), f.Items[90].SellPrice)
}

func TestResearchFile_Unpack(t *testing.T) {
	var f geoscape.ResearchFile
	mustUnpack(t, "RESEARCH.DAT", &f)

	assert.Equal(t, uint16(50), f.Topics[0].Time)
	assert.Equal(t, uint8(10), f.Topics[0].Points)
	assert.Equal(t, int32(1), f.Topics[0].PrerequisitesMet)
	assert.Equal(t, uint16(180), f.Topics[1].Time)
	assert.Equal(t, uint8(0x17), f.Topics[1].UfopaediaEntry)
}

func TestUpFile_Unpack(t *testing.T) {
	var f geoscape.UpFile
	mustUnpack(t, "UP.DAT", &f)

	assert.Equal(t, uint16(4), f.Entries[0].StringIndex)
	assert.Equal(t, uint8(3), f.Entries[0].Research)
	assert.Equal(t, uint8(1), f.Entries[1].SubIndex)
	assert.Equal(t, uint8(2), f.Entries[1].Research)
}
