package savegame_test

import (
	"os"
	"path/filepath"
	"testing"

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
