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
	"github.com/redtoad/xcom-editor/internal/battlescape"
	"github.com/redtoad/xcom-editor/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeTerrainSet writes a terrain set to root/TERRAIN with one record per
// tile type (the record of tile type n uses sprite n) and one red sprite
// per record.
func writeTerrainSet(t *testing.T, root, name string, tileTypes []uint8) {
	t.Helper()
	dir := filepath.Join(root, "TERRAIN")
	require.NoError(t, os.MkdirAll(dir, 0755))
	mcd := make([]byte, len(tileTypes)*resources.MCDRecordSize)
	var pck, tab []byte
	for i, tt := range tileTypes {
		rec := mcd[i*resources.MCDRecordSize:]
		rec[0] = uint8(i)  // first frame
		rec[49] = uint8(i) // PLevel
		rec[53] = tt
		tab = binary.LittleEndian.AppendUint16(tab, uint16(len(pck)))
		// a single row of colored pixels (color 1)
		pck = append(pck, 0)
		for range terrainSpriteWidth {
			pck = append(pck, 1)
		}
		pck = append(pck, 0xFF)
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, name+".MCD"), mcd, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name+".PCK"), pck, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, name+".TAB"), tab, 0644))
}

// writePalettes writes PALETTES.DAT with 5 palettes in which color 1 is red.
func writePalettes(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(root, "GEODATA")
	require.NoError(t, os.MkdirAll(dir, 0755))
	buf := make([]byte, 5*(768+6))
	for p := 0; p < 5; p++ {
		buf[p*(768+6)+3] = 63 // color 1: red
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "PALETTES.DAT"), buf, 0644))
}

func TestLoadTerrain(t *testing.T) {
	root := t.TempDir()
	_, err := loadTerrain(root)
	assert.Error(t, err, "no game files")

	writePalettes(t, root)
	writeTerrainSet(t, root, "BLANKS", []uint8{0, 0})
	writeTerrainSet(t, root, "TEST", []uint8{0, 1, 2, 3})

	sets, err := loadTerrain(root)
	require.NoError(t, err)
	require.Len(t, sets, 2)
	set := sets["TEST"]
	require.NotNil(t, set)
	assert.Len(t, set.records, 4)
	assert.Equal(t, 4, set.frames)

	sheet, err := png.Decode(bytes.NewReader(set.sheet))
	require.NoError(t, err)
	assert.Equal(t, 4*terrainSpriteWidth, sheet.Bounds().Dx())
	assert.Equal(t, terrainSpriteHeight, sheet.Bounds().Dy())
	r, _, _, a := sheet.At(terrainSpriteWidth, 0).RGBA()
	assert.Equal(t, [2]uint32{0xFCFC, 0xFFFF}, [2]uint32{r, a}, "first row of the second sprite is red")
	_, _, _, a = sheet.At(0, 1).RGBA()
	assert.Zero(t, a, "the rest of the sprite is transparent")
}

func TestHandleGetBattleTerrain(t *testing.T) {
	loadTestGames(t, "GAME_1", "GAME_4")
	defer func() { terrain = nil }()

	get := func(path, route string, handler http.HandlerFunc) *httptest.ResponseRecorder {
		r := mux.NewRouter()
		r.HandleFunc(route, handler)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		return rec
	}
	const route = "/api/games/{slot}/battle/terrain"

	terrain = nil
	assert.Equal(t, http.StatusNotFound, get("/api/games/GAME_4/battle/terrain", route, handleGetBattleTerrain).Code,
		"no game graphics")

	root := t.TempDir()
	writePalettes(t, root)
	writeTerrainSet(t, root, "BLANKS", []uint8{0, 0})
	writeTerrainSet(t, root, "TEST", []uint8{0, 1, 2, 3})
	sets, err := loadTerrain(root)
	require.NoError(t, err)
	terrain = sets

	assert.Equal(t, http.StatusNotFound, get("/api/games/GAME_1/battle/terrain", route, handleGetBattleTerrain).Code,
		"not in a mission")
	assert.Equal(t, http.StatusUnprocessableEntity, get("/api/games/GAME_4/battle/terrain", route, handleGetBattleTerrain).Code,
		"the sets do not match the map")

	// a small map built from both sets
	b := games["GAME_4"].sg.Battle
	b.Globals.SizeX, b.Globals.SizeY, b.Globals.SizeZ, b.Globals.MCDRecords = 2, 1, 1, 6
	b.Map.Tiles = []battlescape.MapTile{{Ground: 2, WestWall: 3, NorthWall: 4, Object: 5}, {Ground: 1}}

	var result battleTerrainJS
	getJSON(t, "/api/games/GAME_4/battle/terrain", handleGetBattleTerrain, route, &result)
	assert.Equal(t, []string{"BLANKS", "TEST"}, result.Sets)
	require.Len(t, result.Parts, 6)
	assert.Equal(t, terrainPartJS{Set: 1, Frame: 3, Y: 3}, result.Parts[5])
	assert.Equal(t, []int{2, 3, 4, 5, 1, 0, 0, 0}, result.Tiles)

	rec := get("/api/terrain/test.png", "/api/terrain/{name}.png", handleGetTerrainSheet)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusNotFound, get("/api/terrain/nope.png", "/api/terrain/{name}.png", handleGetTerrainSheet).Code)
}
