package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
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

func TestHandleListLocations_Radars(t *testing.T) {
	loadTestGames(t, "GAME_2")

	var locations []struct {
		Name   string    `json:"name"`
		Radars []radarJS `json:"radars"`
	}
	getJSON(t, "/api/games/GAME_2/locations", handleListLocations, "/api/games/{slot}/locations", &locations)

	radars := map[string][]radarJS{}
	for _, loc := range locations {
		radars[loc.Name] = loc.Radars
	}
	assert.Equal(t, []radarJS{
		{Facility: "Large Radar System", Range: 2250, Chance: 20},
		{Facility: "Small Radar System", Range: 1500, Chance: 10},
	}, radars["Alpine HQ"])
	assert.Equal(t, []radarJS{
		{Facility: "Large Radar System", Range: 2250, Chance: 20},
	}, radars["Nortern Base"])
	assert.Nil(t, radars["UFO-0 (supply ship, Alien Supply)"])
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

func TestHandlePersonalArmour(t *testing.T) {
	loadTestGames(t, "GAME_1")

	r := mux.NewRouter()
	r.HandleFunc("/api/games/{slot}/actions/personal-armour", handlePersonalArmour)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/games/GAME_1/actions/personal-armour", nil))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	for _, s := range games["GAME_1"].sg.Soldiers() {
		if !s.IsDead() {
			assert.NotEqual(t, geoscape.NoArmor, s.Armor(), "soldier %s", s.Name())
		}
	}
}

func TestHandleResurrectSoldier(t *testing.T) {
	loadTestGames(t, "GAME_1")

	// mark the first soldier as dead (rank 0xFFFF) and reload the savegame
	path := filepath.Join(games["GAME_1"].sg.Path, "SOLDIER.DAT")
	buf, err := os.ReadFile(path)
	require.NoError(t, err)
	buf[0], buf[1] = 0xff, 0xff
	require.NoError(t, os.WriteFile(path, buf, 0644))
	sg, err := savegame.Load(games["GAME_1"].sg.Path)
	require.NoError(t, err)
	games["GAME_1"].sg = sg
	require.True(t, sg.Soldiers()[0].IsDead())

	post := func(path string) *httptest.ResponseRecorder {
		r := mux.NewRouter()
		r.HandleFunc("/api/games/{slot}/soldiers/{idx}/resurrect", handleResurrectSoldier)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
		return rec
	}

	type soldier struct {
		Rank     string `json:"rank"`
		LastRank string `json:"lastRank"`
		IsDead   bool   `json:"isDead"`
	}
	var before soldier
	getJSON(t, "/api/games/GAME_1/soldiers/0", handleGetSoldier, "/api/games/{slot}/soldiers/{idx}", &before)
	assert.Equal(t, "Dead", before.Rank)
	assert.True(t, before.IsDead)
	assert.Equal(t, sg.Soldiers()[0].LastRank().String(), before.LastRank)

	rec := post("/api/games/GAME_1/soldiers/0/resurrect")
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.False(t, sg.Soldiers()[0].IsDead())

	// the soldier comes back with the estimated rank
	var after soldier
	getJSON(t, "/api/games/GAME_1/soldiers/0", handleGetSoldier, "/api/games/{slot}/soldiers/{idx}", &after)
	assert.Equal(t, before.LastRank, after.Rank)
	assert.Equal(t, after.Rank, after.LastRank)

	// living soldiers cannot be resurrected
	assert.Equal(t, http.StatusBadRequest, post("/api/games/GAME_1/soldiers/0/resurrect").Code)
	assert.Equal(t, http.StatusNotFound, post("/api/games/GAME_1/soldiers/999/resurrect").Code)
	assert.Equal(t, http.StatusBadRequest, post("/api/games/GAME_1/soldiers/x/resurrect").Code)
}

func TestHandleGetBattleMap(t *testing.T) {
	loadTestGames(t, "GAME_1", "GAME_4")

	var m battleMapJS
	getJSON(t, "/api/games/GAME_4/battle", handleGetBattleMap, "/api/games/{slot}/battle", &m)
	assert.Equal(t, 2, m.Turn)
	assert.Equal(t, [3]int{50, 50, 2}, [3]int{m.SizeX, m.SizeY, m.SizeZ})
	assert.Len(t, m.Tiles, 50*50*2)
	assert.Empty(t, m.Crafts)
	require.Len(t, m.Units, 21)
	assert.Equal(t, battleUnitJS{
		Index: 0, Name: "Rudi Richter", Owner: "xcom", X: 41, Y: 21, Z: 0,
		Alive: true, Visible: true, Health: 40, MaxHealth: m.Units[0].MaxHealth,
		TimeUnits: m.Units[0].TimeUnits, MaxTU: m.Units[0].MaxTU, Morale: m.Units[0].Morale,
	}, m.Units[0])
	assert.Equal(t, "Sectoid Commander", m.Units[20].Name, "duplicate spaces are removed")
	assert.Equal(t, "alien", m.Units[20].Owner)
	// the tile of the first soldier has a floor and was seen by X-COM
	tile := m.Tiles[21*50+41]
	assert.NotZero(t, tile&tileGround)
	assert.NotZero(t, tile&tileSeen)

	r := mux.NewRouter()
	r.HandleFunc("/api/games/{slot}/battle", handleGetBattleMap)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/games/GAME_1/battle", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleHealSoldier(t *testing.T) {
	loadTestGames(t, "GAME_1")
	sg := games["GAME_1"].sg

	var wounded, healthy *savegame.Soldier
	for _, s := range sg.Soldiers() {
		if s.IsDead() {
			continue
		}
		if wounded == nil {
			wounded = s
		} else if healthy == nil {
			healthy = s
		}
	}
	require.NotNil(t, wounded)
	require.NotNil(t, healthy)
	wounded.SetRecoveryDays(14)
	healthy.SetRecoveryDays(0)

	post := func(path string) *httptest.ResponseRecorder {
		r := mux.NewRouter()
		r.HandleFunc("/api/games/{slot}/soldiers/{idx}/heal", handleHealSoldier)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, path, nil))
		return rec
	}
	url := func(s *savegame.Soldier) string {
		return fmt.Sprintf("/api/games/GAME_1/soldiers/%d/heal", s.Index())
	}

	rec := post(url(wounded))
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	assert.False(t, wounded.IsWounded())
	assert.Equal(t, 0, wounded.RecoveryDays())

	// healthy soldiers cannot be healed
	assert.Equal(t, http.StatusBadRequest, post(url(healthy)).Code)
	assert.Equal(t, http.StatusNotFound, post("/api/games/GAME_1/soldiers/999/heal").Code)
	assert.Equal(t, http.StatusBadRequest, post("/api/games/GAME_1/soldiers/x/heal").Code)
}
