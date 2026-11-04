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
	loadTestGames(t, "GAME_4")
	for _, u := range games["GAME_4"].sg.Battle.MapUnits() {
		if u.Index == index {
			return u
		}
	}
	t.Fatalf("no unit in slot %d", index)
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
	t.Run("alien", func(t *testing.T) {
		u := mapUnit(t, 26) // Sectoid, facing 6
		assert.Contains(t, unitSprites(u), unitSpriteJS{Sheet: "SECTOID", Frame: 32 + 6})
	})
	t.Run("unknown type", func(t *testing.T) {
		u := mapUnit(t, 26)
		u.Unit.Type = 8 // Floater, not supported
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
	// all units of GAME_4 have sprites (soldiers, Sectoids and a Cyberdisc)
	assert.Len(t, units, 21)
	assert.Equal(t, 14, units[14].Index)
}
