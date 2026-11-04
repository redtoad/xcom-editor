package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadTestGames loads savegames from testdata/ (at the repository root) into
// a temporary directory and registers them as games.
func loadTestGames(t *testing.T, slots ...string) {
	t.Helper()
	games = make(map[string]*gameEntry)
	for _, slot := range slots {
		src := filepath.Join("..", "..", "testdata", slot)
		dst := filepath.Join(t.TempDir(), slot)
		require.NoError(t, os.MkdirAll(dst, 0755))
		entries, err := os.ReadDir(src)
		require.NoError(t, err)
		for _, e := range entries {
			buf, err := os.ReadFile(filepath.Join(src, e.Name()))
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(filepath.Join(dst, e.Name()), buf, 0644))
		}
		entry, err := loadGame(dst)
		require.NoError(t, err)
		games[slot] = entry
	}
}

// getJSON performs a GET request against the API and decodes the response.
func getJSON(t *testing.T, path string, handler http.HandlerFunc, route string, v any) {
	t.Helper()
	r := mux.NewRouter()
	r.HandleFunc(route, handler)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), v))
}

func TestHandleListGames_Mission(t *testing.T) {
	loadTestGames(t, "GAME_1", "GAME_4")

	var result []struct {
		Slot    string     `json:"slot"`
		Mission *missionJS `json:"mission"`
	}
	getJSON(t, "/api/games", handleListGames, "/api/games", &result)

	require.Len(t, result, 2)
	assert.Equal(t, "GAME_1", result[0].Slot)
	assert.Nil(t, result[0].Mission)
	assert.Equal(t, "GAME_4", result[1].Slot)
	require.NotNil(t, result[1].Mission)
	assert.Equal(t, missionJS{
		Turn: 2, CraftName: "SKYRANGER-2",
		XcomAlive: 13, XcomTotal: 14, AliensAlive: 7, AliensTotal: 7, Score: -27,
	}, *result[1].Mission)
}

func TestHandleGetGame_Mission(t *testing.T) {
	loadTestGames(t, "GAME_4")

	var result struct {
		Mission *missionJS `json:"mission"`
	}
	getJSON(t, "/api/games/GAME_4", handleGetGame, "/api/games/{slot}", &result)
	require.NotNil(t, result.Mission)
	assert.Equal(t, 2, result.Mission.Turn)
}

func TestHandleListLocations_CraftDetails(t *testing.T) {
	loadTestGames(t, "GAME_1", "GAME_2")

	type location struct {
		TypeCode  int    `json:"typeCode"`
		Name      string `json:"name"`
		DestCoord *struct {
			Lat float32 `json:"lat"`
			Lon float32 `json:"lon"`
		} `json:"destCoord"`
		Craft *craftDetails `json:"craft"`
	}
	find := func(locs []location, name string) location {
		for _, loc := range locs {
			if loc.Name == name {
				return loc
			}
		}
		t.Fatalf("location %s not found", name)
		return location{}
	}

	var game1 []location
	getJSON(t, "/api/games/GAME_1/locations", handleListLocations, "/api/games/{slot}/locations", &game1)

	skyranger := find(game1, "SKYRANGER-1")
	assert.NotNil(t, skyranger.DestCoord)
	assert.Equal(t, &craftDetails{
		Type: "Skyranger", Status: "Out", Altitude: "very low", Speed: 760, Fuel: 1493,
		BaseName: "Alpine HQ", Destination: "Crash Site-21",
	}, skyranger.Craft)
	assert.Nil(t, find(game1, "Alpine HQ").Craft)

	var game2 []location
	getJSON(t, "/api/games/GAME_2/locations", handleListLocations, "/api/games/{slot}/locations", &game2)

	ufo := find(game2, "UFO-0 (supply ship, Alien Supply)")
	assert.NotNil(t, ufo.DestCoord)
	require.NotNil(t, ufo.Craft)
	assert.Equal(t, "supply ship", ufo.Craft.Type)
	assert.Equal(t, "Sectoid", ufo.Craft.Race)
	assert.Equal(t, "Alien Supply", ufo.Craft.Mission)
	assert.Equal(t, "Antarctica", ufo.Craft.Zone)
	assert.Equal(t, 3200, ufo.Craft.Speed)
}

func TestHandleListCraft_Damaged(t *testing.T) {
	loadTestGames(t, "GAME_4")

	var result []struct {
		Index   int  `json:"index"`
		Damaged bool `json:"damaged"`
	}
	getJSON(t, "/api/games/GAME_4/craft", handleListCraft, "/api/games/{slot}/craft", &result)

	require.Len(t, result, 3)
	for _, c := range result {
		assert.Equal(t, c.Index != 0, c.Damaged, "craft %d", c.Index)
	}
}
