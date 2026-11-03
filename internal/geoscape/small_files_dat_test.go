package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZonalFile_Unpack(t *testing.T) {
	var f geoscape.ZonalFile
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_1", "ZONAL.DAT"), binary.LittleEndian, &f))

	assert.Equal(t, uint8(0x12), f.Priorities[geoscape.ZoneNorthAmerica])
	assert.Equal(t, uint8(0x04), f.Priorities[geoscape.ZoneArctic])
	assert.Equal(t, uint8(0x10), f.Priorities[geoscape.ZoneEurope])
	assert.Equal(t, uint8(0x05), f.Priorities[geoscape.ZonePacific])
	assert.Equal(t, "South East Asia", geoscape.ZoneSoutheastAsia.String())
}

func TestActsFile_Unpack(t *testing.T) {
	var f geoscape.ActsFile
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_1", "ACTS.DAT"), binary.LittleEndian, &f))

	assert.Equal(t, geoscape.ZoneActs{
		Research: 0x0e, Harvest: 0x10, Abduction: 0x14, Infiltration: 0x14, Base: 0x14,
	}, f.Zones[geoscape.ZoneNorthAmerica])
	// no abductions in Antarctica
	assert.Equal(t, uint8(0), f.Zones[geoscape.ZoneAntarctica].Abduction)
	assert.Equal(t, uint8(0x1c), f.Zones[geoscape.ZonePacific].Base)
}

func TestXbasesFile_Unpack(t *testing.T) {
	var f geoscape.XbasesFile
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_1", "XBASES.DAT"), binary.LittleEndian, &f))

	for zone, target := range f.Zones {
		assert.Equal(t, geoscape.XbaseTarget{}, target, "zone %v", geoscape.Zone(zone))
	}
}

func TestUiglobFile_Unpack(t *testing.T) {
	var f geoscape.UiglobFile
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_3", "UIGLOB.DAT"), binary.LittleEndian, &f))

	assert.Equal(t, uint16(15), f.NextUFO)
	assert.Equal(t, uint16(6), f.NextTerrorSite)
	assert.Equal(t, uint16(6), f.MonthsElapsed)
	assert.Equal(t, uint16(1999), f.Year)
	assert.Equal(t, [12]int16{450, 570, 485, 430, 510}, f.MonthlyScores)
}

func TestXcomFile_Unpack(t *testing.T) {
	var f geoscape.XcomFile
	buf := MustLoadTestFile(t, "GAME_1", "XCOM.DAT")
	require.NoError(t, restruct.Unpack(buf, binary.LittleEndian, &f))

	assert.Equal(t, int32(binary.LittleEndian.Uint32(buf[0:])), f.Months[0].USA)
	assert.Equal(t, int32(binary.LittleEndian.Uint32(buf[0x40:])), f.Months[0].NorthAmerica)
}
