package savegame_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// copySavegame copies a real savegame from testdata/ into a temporary directory.
func copySavegame(t *testing.T, game string) string {
	t.Helper()
	src := filepath.Join("..", "testdata", game)
	dst := t.TempDir()
	entries, err := os.ReadDir(src)
	require.NoError(t, err)
	for _, e := range entries {
		buf, err := os.ReadFile(filepath.Join(src, e.Name()))
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(dst, e.Name()), buf, 0644))
	}
	return dst
}

// savegames are real savegames in testdata/. GAME_4, GAME_5 and GAME_6 were saved
// during a mission.
var savegames = []string{"GAME_1", "GAME_2", "GAME_3", "GAME_4", "GAME_5", "GAME_6"}

func TestSavegame_AllFilesHandled(t *testing.T) {
	handled := map[string]bool{}
	for _, name := range savegame.HandledFiles(&savegame.Savegame{}) {
		handled[name] = true
	}
	for _, game := range savegames {
		entries, err := os.ReadDir(filepath.Join("..", "testdata", game))
		require.NoError(t, err)
		for _, e := range entries {
			assert.True(t, handled[e.Name()], "%s/%s is not loaded as part of the savegame", game, e.Name())
		}
	}
}

func TestSavegame_RoundTrip(t *testing.T) {
	for _, game := range savegames {
		t.Run(game, func(t *testing.T) {
			dir := copySavegame(t, game)

			sg, err := savegame.Load(dir)
			require.NoError(t, err)
			require.NoError(t, sg.Save())

			entries, err := os.ReadDir(dir)
			require.NoError(t, err)
			for _, e := range entries {
				original, err := os.ReadFile(filepath.Join("..", "testdata", game, e.Name()))
				require.NoError(t, err)
				saved, err := os.ReadFile(filepath.Join(dir, e.Name()))
				require.NoError(t, err)
				assert.Equal(t, original, saved, "%s changed after load and save", e.Name())
			}
		})
	}
}

func TestSavegame_SaveGeoscapeFiles(t *testing.T) {
	dir := copySavegame(t, "GAME_3")

	sg, err := savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, uint16(1999), sg.Counters.Year)

	sg.Counters.MonthlyScores[0] = -100
	sg.ZonePriorities.Priorities[0] = 0
	sg.Countries.Countries[0].Satisfaction = 3
	sg.ResearchTopics.Topics[0].Time = 1
	require.NoError(t, sg.Save())

	reloaded, err := savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, int16(-100), reloaded.Counters.MonthlyScores[0])
	assert.Equal(t, uint8(0), reloaded.ZonePriorities.Priorities[0])
	assert.Equal(t, int16(3), reloaded.Countries.Countries[0].Satisfaction)
	assert.Equal(t, uint16(1), reloaded.ResearchTopics.Topics[0].Time)
}

func TestSavegame_Battle(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)
	assert.Nil(t, sg.Battle, "GAME_1 was not saved during a mission")

	dir := copySavegame(t, "GAME_4")
	sg, err = savegame.Load(dir)
	require.NoError(t, err)
	require.NotNil(t, sg.Battle)
	assert.Equal(t, uint16(2), sg.Battle.Globals.Turn)
	assert.Equal(t, "Mark Mitchell", sg.Battle.Units.Units[3].Name.String())

	sg.Battle.Units.Units[3].TimeUnits = 99
	require.NoError(t, sg.Save())

	reloaded, err := savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, uint8(99), reloaded.Battle.Units.Units[3].TimeUnits)
}

func TestSavegame_MissionStatus(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)
	assert.Nil(t, sg.MissionStatus())

	sg, err = savegame.Load(copySavegame(t, "GAME_4"))
	require.NoError(t, err)
	status := sg.MissionStatus()
	require.NotNil(t, status)
	assert.Equal(t, 2, status.Turn)
	assert.Equal(t, -27, status.Score)
	assert.Equal(t, 14, status.XcomTotal)
	assert.Equal(t, 13, status.XcomAlive)
	assert.Equal(t, 13, status.AliensTotal)
	assert.Equal(t, 13, status.AliensAlive)
	require.NotNil(t, status.Craft)
	assert.Equal(t, "SKYRANGER-2", status.Craft.Name())
	assert.Equal(t, 0, status.Craft.Index(), "the craft that is out on the mission")
}

func TestSavegame_Crafts(t *testing.T) {
	tests := []struct {
		game string
		want []string
	}{
		{"GAME_1", []string{"SKYRANGER-1", "FIRESTORM-1", "AVENGER-2", "INTERCEPTOR-3", "INTERCEPTOR-4", "INTERCEPTOR-5", "INTERCEPTOR-6"}},
		{"GAME_2", []string{"SKYRANGER-2", "INTERCEPTOR-1", "INTERCEPTOR-2", "UFO-0"}},
		// only crafts referenced by LOC.DAT exist (the other entries of CRAFT.DAT are zeroed)
		{"GAME_4", []string{"SKYRANGER-2", "SKYRANGER-1", "SKYRANGER-2"}},
	}
	for _, tt := range tests {
		t.Run(tt.game, func(t *testing.T) {
			sg, err := savegame.Load(copySavegame(t, tt.game))
			require.NoError(t, err)
			var names []string
			for _, c := range sg.Crafts() {
				names = append(names, c.Name())
			}
			assert.Equal(t, tt.want, names)
		})
	}
}

func TestCraft_Damaged(t *testing.T) {
	for _, game := range []string{"GAME_1", "GAME_2", "GAME_4"} {
		sg, err := savegame.Load(copySavegame(t, game))
		require.NoError(t, err)
		for _, c := range sg.Crafts() {
			// the entries 1 and 2 of CRAFT.DAT were zeroed in GAME_4
			want := game == "GAME_4" && (c.Index() == 1 || c.Index() == 2)
			assert.Equal(t, want, c.Damaged(), "%s: %s (%d)", game, c.Name(), c.Index())
		}
	}
}

func TestBattle_MapUnits(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_4"))
	require.NoError(t, err)
	units := sg.Battle.MapUnits()
	require.Len(t, units, 27, "14 X-COM units and 13 aliens (large units only once)")

	first := units[0]
	assert.Equal(t, 0, first.Index)
	assert.Equal(t, "Rudi Richter", first.Unit.Name.String())
	assert.Equal(t, [3]int{41, 21, 0}, [3]int{first.X, first.Y, first.Z})
	assert.True(t, first.Alive())

	dead := units[8]
	assert.Equal(t, "Shigeru Sato", dead.Unit.Name.String())
	assert.False(t, dead.Alive())

	disc := units[14]
	assert.Equal(t, "Cyberdisc", disc.Unit.Name.String())
	assert.True(t, disc.Large)
	assert.Equal(t, uint8(battlescape.OwnerAlien), disc.Owner)
	assert.Equal(t, 18, units[15].Index, "the 3 other slots of the large unit are skipped")
	assert.Equal(t, 15, units[15].Ref, "large units have only one entry in UNITREF.DAT")
	assert.Equal(t, 22, units[16].Index)
	assert.Equal(t, 16, units[16].Ref)
}

func TestBattle_Tile(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_4"))
	require.NoError(t, err)
	sx, sy, sz := sg.Battle.Size()
	assert.Equal(t, [3]int{50, 50, 2}, [3]int{sx, sy, sz})

	// every unit stands on a tile with a floor
	for _, u := range sg.Battle.MapUnits() {
		tile, _ := sg.Battle.Tile(u.X, u.Y, u.Z)
		assert.NotZero(t, tile.Ground, "unit %d", u.Index)
	}
	// levels are stored from the top down: only the upper level has tiles
	// without a floor, the ground level has none
	floorless := func(z int) (n int) {
		for y := 0; y < sy; y++ {
			for x := 0; x < sx; x++ {
				if tile, _ := sg.Battle.Tile(x, y, z); tile.Ground == 0 {
					n++
				}
			}
		}
		return n
	}
	assert.Equal(t, 813, floorless(0), "top level")
	assert.Equal(t, 0, floorless(sz-1), "ground level")

	_, seen := sg.Battle.Tile(41, 21, 0)
	assert.True(t, seen, "X-COM has seen the tile of its own soldier")
	tile, seen := sg.Battle.Tile(50, 0, 0)
	assert.Equal(t, battlescape.MapTile{}, tile)
	assert.False(t, seen)
}

func TestBattle_CraftAreas(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_4"))
	require.NoError(t, err)
	assert.Empty(t, sg.Battle.CraftAreas(), "no crafts on the map in this mission")

	g := &sg.Battle.Geodata
	g.Crafts = 2
	g.XcomCraft, g.XcomCraftX, g.XcomCraftY = 0, 1, 2
	g.AlienCraft, g.AlienCraftX, g.AlienCraftY = 2, 3, 0
	assert.Equal(t, []savegame.CraftArea{
		{Alien: false, Type: 0, X: 10, Y: 20, Width: 10, Height: 20},
		{Alien: true, Type: 2, X: 30, Y: 0, Width: 20, Height: 20},
	}, sg.Battle.CraftAreas())
	assert.Equal(t, "Skyranger", sg.Battle.CraftAreas()[0].Name())
	assert.Equal(t, "large scout", sg.Battle.CraftAreas()[1].Name())
}

func TestBattle_MatchTerrainSets(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_4"))
	require.NoError(t, err)
	b := sg.Battle
	// a small map of two terrain sets: BLANKS (2 records) and a set with
	// a floor (2), a west wall (3), a north wall (4) and an object (5)
	b.Globals.SizeX, b.Globals.SizeY, b.Globals.SizeZ = 2, 1, 1
	b.Globals.MCDRecords = 6
	b.Map.Tiles = []battlescape.MapTile{{Ground: 2, WestWall: 3, NorthWall: 4, Object: 5}, {Ground: 2}}

	blanks := savegame.TerrainSet{Name: "BLANKS", TileTypes: []uint8{0, 0}}
	match := savegame.TerrainSet{Name: "MATCH", TileTypes: []uint8{0, 1, 2, 3}}
	// same size, but the walls are swapped
	swapped := savegame.TerrainSet{Name: "SWAPPED", TileTypes: []uint8{0, 2, 1, 3}}
	other := savegame.TerrainSet{Name: "OTHER", TileTypes: []uint8{0, 0, 0}}

	sets, err := b.MatchTerrainSets([]savegame.TerrainSet{other, swapped, match, blanks})
	require.NoError(t, err)
	assert.Equal(t, []string{"BLANKS", "MATCH"}, sets)

	_, err = b.MatchTerrainSets([]savegame.TerrainSet{blanks, swapped})
	assert.Error(t, err, "no combination matches")
}

func TestSavegame_LeftOverBattleFiles(t *testing.T) {
	// a savegame overwritten on the Geoscape still contains the Battlescape
	// files of the earlier mission, but SAVEINFO.DAT marks it as Geoscape save
	dir := copySavegame(t, "GAME_4")
	info := filepath.Join(dir, "SAVEINFO.DAT")
	buf, err := os.ReadFile(info)
	require.NoError(t, err)
	buf[38], buf[39] = 0, 0 // TacticalSave
	require.NoError(t, os.WriteFile(info, buf, 0644))

	sg, err := savegame.Load(dir)
	require.NoError(t, err)
	assert.Nil(t, sg.Battle, "left-over Battlescape files are ignored")
	assert.Nil(t, sg.MissionStatus())

	// saving leaves the left-over files untouched
	before, err := os.ReadFile(filepath.Join(dir, "UNITREF.DAT"))
	require.NoError(t, err)
	require.NoError(t, sg.Save())
	after, err := os.ReadFile(filepath.Join(dir, "UNITREF.DAT"))
	require.NoError(t, err)
	assert.Equal(t, before, after)
}

func TestCraft_Base(t *testing.T) {
	// the base reference of a craft is an index into LOC.DAT, not BASE.DAT
	sg, err := savegame.Load(copySavegame(t, "GAME_2"))
	require.NoError(t, err)

	crafts := sg.Crafts()
	require.Equal(t, "INTERCEPTOR-1", crafts[1].Name())
	require.Equal(t, 0, crafts[1].Base().Index())

	// station the interceptor at "Nortern Base" (LOC.DAT entry 5, BASE.DAT entry 1)
	sg.CraftsFile().Crafts[1].BaseReference = 5
	require.NotNil(t, crafts[1].Base())
	assert.Equal(t, "Nortern Base", crafts[1].Base().Name())

	// LOC.DAT entries which are no X-COM base or out of range
	sg.CraftsFile().Crafts[1].BaseReference = 1 // SKYRANGER-2
	assert.Nil(t, crafts[1].Base())
	sg.CraftsFile().Crafts[1].BaseReference = 8 // beyond the last base, see GAME_6
	assert.Nil(t, crafts[1].Base())
	sg.CraftsFile().Crafts[1].BaseReference = 1000
	assert.Nil(t, crafts[1].Base())
}

func TestCraft_UFO(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_2"))
	require.NoError(t, err)

	crafts := sg.Crafts()
	ufo := crafts[3]
	require.True(t, ufo.IsUFO())
	assert.False(t, crafts[1].IsUFO())

	// UFOs have no base, although their base reference (0) points to the
	// first base in LOC.DAT
	require.Equal(t, geoscape.XCOMBase, sg.LocationFile().Objects[0].Type)
	require.Equal(t, 0, int(sg.CraftsFile().Crafts[ufo.Index()].BaseReference))
	assert.Nil(t, ufo.Base())

	// detected UFOs are numbered by their alien ship entry in LOC.DAT
	for i, obj := range sg.LocationFile().Objects {
		if obj.Type == geoscape.AlienShip && int(obj.TableReference) == ufo.Index() {
			sg.LocationFile().Objects[i].CountSuffix = 7
		}
	}
	assert.Equal(t, "UFO-7", ufo.Name())
}

func TestCraft_BaseNotInBaseFile(t *testing.T) {
	// the interceptors of GAME_1 are stationed at bases in LOC.DAT whose
	// entries in BASE.DAT are empty
	sg, err := savegame.Load(filepath.Join("..", "testdata", "GAME_1"))
	require.NoError(t, err)

	bases := map[string]string{}
	for _, c := range sg.Crafts() {
		if b := c.Base(); b != nil {
			bases[c.Name()] = b.Name()
		}
	}
	assert.Equal(t, map[string]string{
		"SKYRANGER-1": "Alpine HQ", "FIRESTORM-1": "Alpine HQ", "AVENGER-2": "Alpine HQ",
	}, bases)
}

func TestSavegame_BaseOutOfRange(t *testing.T) {
	sg, err := savegame.Load(filepath.Join("..", "testdata", "GAME_1"))
	require.NoError(t, err)
	assert.Nil(t, sg.Base(-1))
	assert.Nil(t, sg.Base(len(sg.BasesData.Bases)))
}

func TestBase_Radars(t *testing.T) {
	sg, err := savegame.Load(filepath.Join("..", "testdata", "GAME_2"))
	require.NoError(t, err)

	bases := sg.Bases()
	require.Len(t, bases, 2)
	assert.Equal(t, []savegame.Radar{
		{Facility: geoscape.LargeRadarSystem, Range: 2578, Chance: 20},
		{Facility: geoscape.SmallRadarSystem, Range: 1696, Chance: 10},
	}, bases[0].Radars())
	assert.Equal(t, []savegame.Radar{
		{Facility: geoscape.LargeRadarSystem, Range: 2578, Chance: 20},
	}, bases[1].Radars())
}

func TestBase_RadarsIgnoresUnfinishedFacilities(t *testing.T) {
	sg, err := savegame.Load(filepath.Join("..", "testdata", "GAME_2"))
	require.NoError(t, err)

	data := &sg.BasesData.Bases[1]
	built := false
	for i := range data.Grid {
		switch {
		case data.Grid[i] == geoscape.LargeRadarSystem:
			data.DaysToCompletion[i] = 5 // still under construction
		case data.Grid[i] == geoscape.Empty && !built:
			data.Grid[i] = geoscape.HyperwaveDecoder
			data.DaysToCompletion[i] = 0
			built = true
		}
	}
	require.True(t, built, "no empty tile for the hyperwave decoder")
	assert.Equal(t, []savegame.Radar{
		{Facility: geoscape.HyperwaveDecoder, Range: 2759, Chance: 100},
	}, sg.Bases()[1].Radars())
}

func TestBattle_MapUnits_LargeUnitFirst(t *testing.T) {
	// In GAME_6 the tank occupies UNITPOS.DAT slots 0-3, but only entry 0 of
	// UNITREF.DAT: all following units are shifted by 3 slots.
	sg, err := savegame.Load(copySavegame(t, "GAME_6"))
	require.NoError(t, err)
	units := sg.Battle.MapUnits()
	byName := map[string]savegame.MapUnit{}
	for _, u := range units {
		byName[u.Unit.Name.String()] = u
	}
	tank := units[0]
	assert.Equal(t, "Tank/Cannon", tank.Unit.Name.String())
	assert.True(t, tank.Large)

	// soldiers that are still aboard the Skyranger
	for name, pos := range map[string][3]int{
		"Gudrun Brehme": {24, 34, 2},
		"Sarah Jonlan":  {24, 35, 2},
		"Frank Sharpe":  {25, 35, 2},
	} {
		u, ok := byName[name]
		require.True(t, ok, name)
		assert.Equal(t, pos, [3]int{u.X, u.Y, u.Z}, name)
		assert.Equal(t, uint8(battlescape.OwnerXcom), u.Owner, name)
	}
	// the soldiers in UNITREF.DAT 1-3 are no longer hidden behind the tank
	assert.Contains(t, byName, "Donald McNeil")
	for _, u := range units {
		if strings.HasPrefix(u.Unit.Name.String(), "Sectoid") {
			assert.Equal(t, uint8(battlescape.OwnerAlien), u.Owner, "%s (slot %d)", u.Unit.Name, u.Index)
		}
	}
}

func TestSoldier_Base(t *testing.T) {
	// the base reference of a soldier is an index into LOC.DAT, not BASE.DAT
	// (GAME_6 has soldiers at LOC.DAT entries 0, 4, 5, 8 and 9)
	sg, err := savegame.Load(filepath.Join("..", "testdata", "GAME_6"))
	require.NoError(t, err)

	bases := map[int]int{} // LOC.DAT entry -> BASE.DAT entry
	for _, s := range sg.Soldiers() {
		b := s.Base()
		require.NotNil(t, b, "soldier %s has no base", s.Name())
		bases[sg.SoldierFile().Soldiers[s.Index()].Base] = b.Index()
	}
	assert.Equal(t, map[int]int{0: 0, 4: 2, 5: 1, 8: 3, 9: 4}, bases)

	// LOC.DAT entries which are no X-COM base or out of range, and soldiers
	// being transferred
	s := sg.Soldiers()[0]
	for _, ref := range []int{-1, 0xFF, 1000} {
		sg.SoldierFile().Soldiers[s.Index()].Base = ref
		assert.Nil(t, s.Base(), "base reference %d", ref)
	}
}
