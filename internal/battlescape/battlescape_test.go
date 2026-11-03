package battlescape_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// battleSavegames are real savegames in testdata/ (at the repository root)
// that were saved during the same mission (in turn 2 and 3).
var battleSavegames = []string{"GAME_4", "GAME_5"}

// MustLoadTestFile loads a DAT file from a savegame in testdata/.
func MustLoadTestFile(t *testing.T, game, name string) []byte {
	t.Helper()
	buf, err := os.ReadFile(filepath.Join("..", "..", "testdata", game, name))
	require.NoError(t, err)
	return buf
}

// mustUnpack unpacks a DAT file from savegame GAME_4 in testdata/.
func mustUnpack(t *testing.T, name string, obj any) {
	t.Helper()
	require.NoError(t, restruct.Unpack(MustLoadTestFile(t, "GAME_4", name), binary.LittleEndian, obj))
}

// TestBattlescapeFilesRoundTrip loads every file of the battle savegames and
// checks that packing it again results in identical bytes.
func TestBattlescapeFilesRoundTrip(t *testing.T) {
	files := []struct {
		name string
		obj  func() any
	}{
		{"BGLOB.DAT", func() any { return &battlescape.BglobFile{} }},
		{"GEODATA.DAT", func() any { return &battlescape.GeodataFile{} }},
		{"MAP.DAT", func() any { return &battlescape.MapFile{} }},
		{"MISDATA.DAT", func() any { return &battlescape.MisdataFile{} }},
		{"OBPOS.DAT", func() any { return &battlescape.ObposFile{} }},
		{"ROUTES.DAT", func() any { return &battlescape.RoutesFile{} }},
		{"SEEMAP.DAT", func() any { return &battlescape.SeemapFile{} }},
		{"SMOKBIT.DAT", func() any { return &battlescape.SmokbitFile{} }},
		{"SMOKREF.DAT", func() any { return &battlescape.SmokrefFile{} }},
		{"SOURCEMP.DAT", func() any { return &battlescape.LightMapFile{} }},
		{"TERMP.DAT", func() any { return &battlescape.LightMapFile{} }},
		{"UNITPOS.DAT", func() any { return &battlescape.UnitposFile{} }},
		{"UNITREF.DAT", func() any { return &battlescape.UnitrefFile{} }},
		{"WGLOB.DAT", func() any { return &battlescape.WglobFile{} }},
	}

	for _, game := range battleSavegames {
		for _, f := range files {
			t.Run(game+"/"+f.name, func(t *testing.T) {
				buf := MustLoadTestFile(t, game, f.name)
				obj := f.obj()

				require.NoError(t, restruct.Unpack(buf, binary.LittleEndian, obj))
				packed, err := restruct.Pack(binary.LittleEndian, obj)
				require.NoError(t, err)
				if !bytes.Equal(buf, packed) {
					t.Errorf("re-packed %s does not match original data", f.name)
				}
			})
		}
	}
}

func TestWglobFile_Unpack(t *testing.T) {
	var f battlescape.WglobFile
	mustUnpack(t, "WGLOB.DAT", &f)

	assert.Equal(t, uint16(2), f.Turn)
	assert.Equal(t, uint16(50), f.SizeX)
	assert.Equal(t, uint16(50), f.SizeY)
	assert.Equal(t, uint16(2), f.SizeZ)
	assert.Equal(t, 5000, f.Tiles())
	assert.Equal(t, uint16(118), f.Routes)
	assert.Equal(t, uint16(1), f.SoldiersLost)
	assert.Equal(t, int16(-27), f.XcomLossScore)
}

func TestTileFiles_MatchMapSize(t *testing.T) {
	var wglob battlescape.WglobFile
	mustUnpack(t, "WGLOB.DAT", &wglob)

	var m battlescape.MapFile
	mustUnpack(t, "MAP.DAT", &m)
	assert.Len(t, m.Tiles, wglob.Tiles())

	var routes battlescape.RoutesFile
	mustUnpack(t, "ROUTES.DAT", &routes)
	assert.Len(t, routes.Nodes, int(wglob.Routes))
	assert.Equal(t, uint8(0xFF), routes.Nodes[0].Links[0].Target)

	var seemap battlescape.SeemapFile
	mustUnpack(t, "SEEMAP.DAT", &seemap)
	assert.Len(t, seemap.Tiles, wglob.Tiles())

	for _, name := range []string{"TERMP.DAT", "SOURCEMP.DAT"} {
		var light battlescape.LightMapFile
		mustUnpack(t, name, &light)
		assert.Len(t, light.Tiles, wglob.Tiles(), name)
	}
}

func TestGeodataFile_Unpack(t *testing.T) {
	var f battlescape.GeodataFile
	mustUnpack(t, "GEODATA.DAT", &f)

	assert.Equal(t, uint16(5), f.ModulesX)
	assert.Equal(t, uint16(5), f.ModulesY)
	assert.Equal(t, uint16(2), f.SizeZ)
	assert.Equal(t, uint16(4), f.Terrain)
}

func TestBglobFile_Unpack(t *testing.T) {
	var f battlescape.BglobFile
	mustUnpack(t, "BGLOB.DAT", &f)

	// alien bases are always dark
	assert.Equal(t, battlescape.BglobFile{Ambient: 16, Shadow: 16}, f)
}

func TestMisdataFile_Unpack(t *testing.T) {
	var f battlescape.MisdataFile
	mustUnpack(t, "MISDATA.DAT", &f)

	assert.Equal(t, int16(3), f.MissionType)
	assert.Equal(t, uint8(1), f.Craft)
	assert.Equal(t, int16(0), f.Language)
}

func TestUnits_Unpack(t *testing.T) {
	var pos battlescape.UnitposFile
	mustUnpack(t, "UNITPOS.DAT", &pos)
	var ref battlescape.UnitrefFile
	mustUnpack(t, "UNITREF.DAT", &ref)

	soldier := ref.Units[3]
	assert.Equal(t, "Mark Mitchell", soldier.Name.String())
	assert.Equal(t, uint8(3), soldier.Soldier)
	assert.Equal(t, uint8(45), soldier.Health)
	assert.Equal(t, uint8(45), soldier.MaxHealth)
	assert.Equal(t, int16(8), soldier.Kills)
	assert.Equal(t, uint8(battlescape.OwnerXcom), pos.Units[3].Owner)
	assert.Equal(t, uint8(11), pos.Units[3].X)
	assert.NotZero(t, pos.Units[3].Flags&battlescape.UnitActive)

	// dead soldier
	assert.Equal(t, "Shigeru Sato", ref.Units[8].Name.String())
	assert.Equal(t, uint8(0), ref.Units[8].Health)
	assert.Zero(t, pos.Units[8].Flags&battlescape.UnitActive)

	// large alien unit
	assert.Equal(t, "Cyberdisc", ref.Units[14].Name.String())
	assert.Equal(t, uint8(battlescape.OwnerAlien), pos.Units[14].Owner)
	assert.NotZero(t, pos.Units[14].Flags&battlescape.UnitLarge)
	assert.NotZero(t, pos.Units[15].Flags&battlescape.UnitRightSection)

	assert.Equal(t, uint8(0xFF), ref.Units[33].Type)
}

func TestObposFile_Unpack(t *testing.T) {
	var f battlescape.ObposFile
	mustUnpack(t, "OBPOS.DAT", &f)

	var used int
	for _, o := range f.Objects {
		if o.Type != 0xFF {
			used++
		}
	}
	assert.Equal(t, 48, used)

	// a clip with 35 rounds loaded into the weapon (object 0) of unit 22
	assert.Equal(t, uint8(22), f.Objects[1].Owner)
	assert.Equal(t, uint16(0), f.Objects[1].Container)
	assert.Equal(t, uint8(35), f.Objects[1].Value1)
	assert.Equal(t, uint16(0xFFFF), f.Objects[0].Container)
}

func TestSmokeFiles_Unpack(t *testing.T) {
	var bits battlescape.SmokbitFile
	mustUnpack(t, "SMOKBIT.DAT", &bits)
	var ref battlescape.SmokrefFile
	mustUnpack(t, "SMOKREF.DAT", &ref)

	for i, c := range ref.Clouds {
		assert.Zero(t, c.Type, "cloud %d", i)
	}
	assert.False(t, bits.HasSmoke(0))
}
