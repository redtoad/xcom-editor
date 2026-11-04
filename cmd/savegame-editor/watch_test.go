package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// copyGame copies a savegame from testdata/ into root and sets the
// modification time of its files.
func copyGame(t *testing.T, root, slot string, mtime time.Time) string {
	t.Helper()
	copySlot(t, root, slot)
	dir := filepath.Join(root, slot)
	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	for _, e := range entries {
		touch(t, dir, e.Name(), mtime)
	}
	return dir
}

// touch sets the modification time of a file of a savegame.
func touch(t *testing.T, dir, name string, mtime time.Time) {
	t.Helper()
	require.NoError(t, os.Chtimes(filepath.Join(dir, name), mtime, mtime))
}

// post performs a POST request against the API.
func post(t *testing.T, path string, handler http.HandlerFunc, route string) {
	t.Helper()
	r := mux.NewRouter()
	r.HandleFunc(route, handler)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
}

// listStale returns the stale flag of all savegames by slot.
func listStale(t *testing.T) map[string]bool {
	t.Helper()
	var result []struct {
		Slot  string `json:"slot"`
		Stale bool   `json:"stale"`
	}
	getJSON(t, "/api/games", handleListGames, "/api/games", &result)
	stale := make(map[string]bool)
	for _, g := range result {
		stale[g.Slot] = g.Stale
	}
	return stale
}

// watchRoot returns a new, empty savegame folder that is watched.
func watchRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	games = make(map[string]*gameEntry)
	currentRoot = root
	loadFailures = make(map[string]time.Time)
	t.Cleanup(func() { currentRoot = "" })
	return root
}

func TestScanGames_LoadsNewSavegames(t *testing.T) {
	root := watchRoot(t)
	old := time.Now().Add(-time.Hour)

	scanGames()
	assert.Empty(t, listStale(t))

	copyGame(t, root, "GAME_1", old)
	scanGames()
	assert.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))
}

func TestScanGames_WaitsUntilSavegameIsWritten(t *testing.T) {
	root := watchRoot(t)
	dir := copyGame(t, root, "GAME_1", time.Now().Add(-time.Hour))
	touch(t, dir, "SOLDIER.DAT", time.Now()) // still being written

	scanGames()
	assert.Empty(t, listStale(t))

	touch(t, dir, "SOLDIER.DAT", time.Now().Add(-2*settleTime))
	scanGames()
	assert.Contains(t, listStale(t), "GAME_1")
}

func TestScanGames_SkipsBrokenSavegames(t *testing.T) {
	root := watchRoot(t)
	old := time.Now().Add(-time.Hour)
	dir := copyGame(t, root, "GAME_1", old)
	require.NoError(t, os.Remove(filepath.Join(dir, "SOLDIER.DAT")))

	scanGames()
	scanGames()
	assert.Empty(t, listStale(t))
	assert.Contains(t, loadFailures, dir)
}

func TestScanGames_MarksChangedSavegamesAsStale(t *testing.T) {
	root := watchRoot(t)
	dir := copyGame(t, root, "GAME_1", time.Now().Add(-time.Hour))
	scanGames()
	require.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))

	// unchanged files do not make the savegame stale
	scanGames()
	assert.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))

	// the game saved the savegame again
	touch(t, dir, "SAVEINFO.DAT", time.Now())
	scanGames()
	assert.Equal(t, map[string]bool{"GAME_1": true}, listStale(t))

	var detail struct {
		Stale bool `json:"stale"`
	}
	getJSON(t, "/api/games/GAME_1", handleGetGame, "/api/games/{slot}", &detail)
	assert.True(t, detail.Stale)

	// reloading picks up the changes
	post(t, "/api/games/GAME_1/reload", handleReload, "/api/games/{slot}/reload")
	assert.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))
	scanGames()
	assert.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))
}

func TestScanGames_IgnoresSavegamesOfPreviousFolder(t *testing.T) {
	root := watchRoot(t)
	copyGame(t, root, "GAME_1", time.Now().Add(-time.Hour))

	// another folder was loaded while the savegame was loaded
	currentRoot = t.TempDir()
	addGame(root, "GAME_1", &gameEntry{})
	assert.Empty(t, games)
}

func TestScanGames_SavingInEditorIsNotAChange(t *testing.T) {
	root := watchRoot(t)
	copyGame(t, root, "GAME_1", time.Now().Add(-time.Hour))
	scanGames()

	post(t, "/api/games/GAME_1/save", handleSave, "/api/games/{slot}/save")
	scanGames()
	assert.Equal(t, map[string]bool{"GAME_1": false}, listStale(t))
}
