package main

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/savegame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mapUnit returns the unit in the given slot of GAME_4.
func mapUnit(t *testing.T, index int) savegame.MapUnit {
	t.Helper()
	return gameMapUnit(t, "GAME_4", index)
}

// gameMapUnit returns the unit in the given slot of a test game. The stats
// are a copy, so tests can change them without affecting other tests.
func gameMapUnit(t *testing.T, game string, index int) savegame.MapUnit {
	t.Helper()
	loadTestGames(t, game)
	for _, u := range games[game].sg.Battle.MapUnits() {
		if u.Index == index {
			stats := *u.Unit
			u.Unit = &stats
			return u
		}
	}
	t.Fatalf("no unit in slot %d of %s", index, game)
	return savegame.MapUnit{}
}

func TestUnitSprites(t *testing.T) {
	t.Run("soldier", func(t *testing.T) {
		u := mapUnit(t, 0) // personal armour, facing 6, male
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "XCOM_1", Frame: spriteRightArm + 6},
			{Sheet: "XCOM_1", Frame: spriteLegs + 6},
			{Sheet: "XCOM_1", Frame: 32 + 6},
			{Sheet: "XCOM_1", Frame: spriteLeftArm + 6},
		}, unitSprites(u))
	})
	t.Run("female soldier", func(t *testing.T) {
		u := mapUnit(t, 2) // facing 6, female
		assert.Contains(t, unitSprites(u), unitSpriteJS{Sheet: "XCOM_1", Frame: 267 + 6})
	})
	t.Run("kneeling", func(t *testing.T) {
		u := mapUnit(t, 0)
		u.Unit.Flags |= unitKneeling
		sprites := unitSprites(u)
		assert.Contains(t, sprites, unitSpriteJS{Sheet: "XCOM_1", Frame: spriteLegsKneel + 6})
		assert.Contains(t, sprites, unitSpriteJS{Sheet: "XCOM_1", Frame: 32 + 6, DY: kneelOffset})
	})
	t.Run("dead", func(t *testing.T) {
		u := mapUnit(t, 8)
		require.False(t, u.Alive())
		assert.Equal(t, []unitSpriteJS{{Sheet: "XCOM_1", Frame: 266}}, unitSprites(u))
	})
	t.Run("large", func(t *testing.T) {
		u := mapUnit(t, 14) // Cyberdisc, facing 4
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "CYBER", Frame: 4},
			{Sheet: "CYBER", Frame: 12, DX: 16, DY: 8},
			{Sheet: "CYBER", Frame: 20, DX: -16, DY: 8},
			{Sheet: "CYBER", Frame: 28, DY: 16},
		}, unitSprites(u))
	})
	t.Run("tank", func(t *testing.T) {
		u := gameMapUnit(t, "GAME_6", 0) // Tank/Cannon, facing 1
		require.NotZero(t, u.Unit.Tank)
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "TANKS", Frame: 1},
			{Sheet: "TANKS", Frame: 9, DX: 16, DY: 8},
			{Sheet: "TANKS", Frame: 17, DX: -16, DY: 8},
			{Sheet: "TANKS", Frame: 25, DY: 16},
			{Sheet: "TANKS", Frame: spriteTurrets + 1, DY: 16 + turretOffsetY}, // cannon turret
		}, unitSprites(u))
	})
	t.Run("hovertank", func(t *testing.T) {
		u := gameMapUnit(t, "GAME_6", 0)
		u.Unit.MovementType = 2 // fly
		u.Unit.TurretImage = 3  // plasma
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "TANKS", Frame: spriteHovertank + 1},
			{Sheet: "TANKS", Frame: spritePropulsion, DX: 16, DY: 8},
			{Sheet: "TANKS", Frame: spriteHovertank + 9, DX: 16, DY: 8},
			{Sheet: "TANKS", Frame: spritePropulsion + 8, DX: -16, DY: 8},
			{Sheet: "TANKS", Frame: spriteHovertank + 17, DX: -16, DY: 8},
			{Sheet: "TANKS", Frame: spritePropulsion + 16, DY: 16},
			{Sheet: "TANKS", Frame: spriteHovertank + 25, DY: 16},
			{Sheet: "TANKS", Frame: spriteTurrets + 3*8 + 1, DX: -7, DY: 16 + turretOffsetY - 3},
		}, unitSprites(u))
	})
	t.Run("alien", func(t *testing.T) {
		u := mapUnit(t, 32) // Sectoid, facing 6
		assert.Contains(t, unitSprites(u), unitSpriteJS{Sheet: "SECTOID", Frame: 32 + 6})
	})
	t.Run("floater", func(t *testing.T) {
		u := mapUnit(t, 32) // Sectoid, facing 6
		u.Unit.Type = 8     // Floater: no legs, right arm before left arm
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "FLOATER", Frame: spriteFloaterRightArm + 6},
			{Sheet: "FLOATER", Frame: 16 + 6},
			{Sheet: "FLOATER", Frame: spriteFloaterLeftArm + 6},
		}, unitSprites(u))
	})
	t.Run("dead floater", func(t *testing.T) {
		u := mapUnit(t, 8) // dead
		u.Unit.Type = 8
		assert.Equal(t, []unitSpriteJS{{Sheet: "FLOATER", Frame: 66}}, unitSprites(u))
	})
	t.Run("reaper", func(t *testing.T) {
		u := mapUnit(t, 14) // Cyberdisc, facing 4
		u.Unit.Type = 12
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "X_REAP", Frame: 4},
			{Sheet: "X_REAP", Frame: 12, DX: 16, DY: 8},
			{Sheet: "X_REAP", Frame: 20, DX: -16, DY: 8},
			{Sheet: "X_REAP", Frame: 28, DY: 16},
		}, unitSprites(u))
	})
	t.Run("dead reaper", func(t *testing.T) {
		u := mapUnit(t, 14) // Cyberdisc
		u.Unit.Type = 12
		u.Unit.Health = 0
		// the corpse (OBDATA.DAT reaper corpse 1-4) lies on the 4 tiles
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "FLOOROB", Frame: 52},
			{Sheet: "FLOOROB", Frame: 55, DX: 16, DY: 8},
			{Sheet: "FLOOROB", Frame: 54, DX: -16, DY: 8},
			{Sheet: "FLOOROB", Frame: 53, DY: 16},
		}, unitSprites(u))
	})
	t.Run("dead cyberdisc", func(t *testing.T) {
		u := mapUnit(t, 14) // Cyberdisc
		u.Unit.Health = 0
		// the corpse (OBDATA.DAT cyber 1-4) lies on the 4 tiles
		assert.Equal(t, []unitSpriteJS{
			{Sheet: "FLOOROB", Frame: 56},
			{Sheet: "FLOOROB", Frame: 57, DX: 16, DY: 8},
			{Sheet: "FLOOROB", Frame: 58, DX: -16, DY: 8},
			{Sheet: "FLOOROB", Frame: 59, DY: 16},
		}, unitSprites(u))
	})
	t.Run("dead large unit without corpse", func(t *testing.T) {
		u := gameMapUnit(t, "GAME_6", 0) // Tank/Cannon
		u.Unit.Health = 0
		assert.Nil(t, unitSprites(u))
	})
	t.Run("unknown type", func(t *testing.T) {
		u := mapUnit(t, 26)
		u.Unit.Type = 6 // Ethereal, not supported
		assert.Nil(t, unitSprites(u))
	})
}

func TestLoadUnitSheets(t *testing.T) {
	root := t.TempDir()
	_, err := loadUnitSheets(root)
	assert.Error(t, err, "no game files")

	writePalettes(t, root)
	dir := filepath.Join(root, "UNITS")
	require.NoError(t, os.MkdirAll(dir, 0755))
	for _, name := range unitSheetFiles {
		// two sprites with a single red row
		var pck, tab []byte
		for range 2 {
			tab = binary.LittleEndian.AppendUint16(tab, uint16(len(pck)))
			pck = append(pck, 0)
			pck = append(pck, bytes.Repeat([]byte{1}, unitSpriteWidth)...)
			pck = append(pck, 0xFF)
		}
		require.NoError(t, os.WriteFile(filepath.Join(dir, name+".PCK"), pck, 0644))
		require.NoError(t, os.WriteFile(filepath.Join(dir, name+".TAB"), tab, 0644))
	}

	sheets, err := loadUnitSheets(root)
	require.NoError(t, err)
	require.Len(t, sheets, len(unitSheetFiles))
	img, err := png.Decode(bytes.NewReader(sheets["SECTOID"]))
	require.NoError(t, err)
	assert.Equal(t, 2*unitSpriteWidth, img.Bounds().Dx())
	assert.Equal(t, unitSpriteHeight, img.Bounds().Dy())

	unitSheets = sheets
	defer func() { unitSheets = nil }()
	r := mux.NewRouter()
	r.HandleFunc("/api/unitsprites/{name}.png", handleGetUnitSheet)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/unitsprites/sectoid.png", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"), "the sprites depend on the game folder")
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/unitsprites/nope.png", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestBattleUnitSprites(t *testing.T) {
	loadTestGames(t, "GAME_4")
	units := battleUnitSprites(games["GAME_4"].sg.Battle)
	// all units of GAME_4 have sprites (soldiers, Sectoids and Cyberdiscs)
	assert.Len(t, units, 27)
	assert.Equal(t, 14, units[14].Index)
}
