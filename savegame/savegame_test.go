package savegame_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/redtoad/xcom-editor/internal/battlescape"
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

// savegames are real savegames in testdata/. GAME_4 and GAME_5 were saved
// during a mission.
var savegames = []string{"GAME_1", "GAME_2", "GAME_3", "GAME_4", "GAME_5"}

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
	assert.Equal(t, 7, status.AliensTotal)
	assert.Equal(t, 7, status.AliensAlive)
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
	require.Len(t, units, 21, "14 X-COM units and 7 aliens (large units only once)")

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
