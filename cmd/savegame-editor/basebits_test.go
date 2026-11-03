package main

import (
	"bytes"
	"image/png"
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

// emptyBase returns the tiles of a base without any facilities.
func emptyBase() []savegame.BaseTile {
	tiles := make([]savegame.BaseTile, 36)
	for i := range tiles {
		tiles[i].Type = geoscape.Empty
	}
	return tiles
}

func TestBaseLayers(t *testing.T) {
	tiles := emptyBase()
	tiles[0].Type = geoscape.AccessLift
	tiles[1].Type = geoscape.LivingQuarters
	tiles[7].Type = geoscape.Laboratory // below the living quarters
	tiles[7].DaysToCompletion = 5
	tiles[2].Type = geoscape.HangarTopLeft
	tiles[3].Type = geoscape.HangarTopRight
	tiles[8].Type = geoscape.HangarBottomLeft
	tiles[9].Type = geoscape.HangarBottomRight

	layers := baseLayers(tiles)
	require.Len(t, layers, 4)
	ground, shapes, corridors, interiors := layers[0], layers[1], layers[2], layers[3]

	assert.Len(t, ground, 36)
	assert.Equal(t, baseSprite{spriteGround, 10, 10}, ground[35])

	assert.ElementsMatch(t, []baseSprite{
		{2, 0, 0},                                     // access lift (octagon)
		{1, 2, 0},                                     // living quarters (square)
		{1 + 3, 2, 2},                                 // laboratory under construction (dotted square)
		{9, 4, 0}, {10, 6, 0}, {11, 4, 2}, {12, 6, 2}, // hangar
	}, shapes)

	assert.ElementsMatch(t, []baseSprite{
		{spriteFacilityStart, 0, 0},                                // access lift
		{spriteFacilityStart + int(geoscape.LivingQuarters), 2, 0}, // living quarters
		{33, 4, 0}, {34, 6, 0}, {35, 4, 2}, {36, 6, 2}, // hangar
	}, interiors, "facilities under construction have no interior")

	// no corridors within the hangar or to the laboratory under construction
	assert.ElementsMatch(t, []baseSprite{
		{spriteCorridorH, 1, 0}, // access lift - living quarters
		{spriteCorridorH, 3, 0}, // living quarters - hangar
	}, corridors)
}

func TestBaseLayers_EmptyBase(t *testing.T) {
	layers := baseLayers(emptyBase())
	assert.Len(t, layers[0], 36)
	assert.Empty(t, layers[1])
	assert.Empty(t, layers[2])
	assert.Empty(t, layers[3])
}

// writeFakeBaseBits writes a (synthetic) BASEBITS.PCK with two sprites and
// a palette into root.
func writeFakeBaseBits(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEODATA"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEOGRAPH"), 0755))

	// sprite 0: one row of color 1, sprite 1: skip a row, then 31 transparent pixels and one pixel of color 2
	sprite0 := append([]byte{0}, append(bytes.Repeat([]byte{1}, 32), 0xFF)...)
	sprite1 := []byte{1, 0xFE, 31, 2, 0xFF}
	pck := append(sprite0, sprite1...)
	tab := []byte{0, 0, byte(len(sprite0)), 0}
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEOGRAPH", "BASEBITS.PCK"), pck, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEOGRAPH", "BASEBITS.TAB"), tab, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEODATA", "PALETTES.DAT"), make([]byte, 5*(768+6)), 0644))
}

func TestLoadBaseBits(t *testing.T) {
	root := t.TempDir()
	writeFakeBaseBits(t, root)

	bits, err := loadBaseBits(root)
	require.NoError(t, err)
	require.Len(t, bits, 2)

	img, err := png.Decode(bytes.NewReader(bits[1]))
	require.NoError(t, err)
	assert.Equal(t, 32, img.Bounds().Dx())
	assert.Equal(t, 2, img.Bounds().Dy())

	_, err = loadBaseBits(t.TempDir())
	assert.Error(t, err, "missing game files")
}

func TestHandleGetBaseBit(t *testing.T) {
	root := t.TempDir()
	writeFakeBaseBits(t, root)
	bits, err := loadBaseBits(root)
	require.NoError(t, err)
	baseBits = bits
	t.Cleanup(func() { baseBits = nil })

	r := mux.NewRouter()
	r.HandleFunc("/api/basebits/{sprite}.png", handleGetBaseBit)
	for path, status := range map[string]int{
		"/api/basebits/0.png": http.StatusOK,
		"/api/basebits/1.png": http.StatusOK,
		"/api/basebits/2.png": http.StatusNotFound,
		"/api/basebits/x.png": http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, status, rec.Code, path)
	}
}

func TestHandleGetBase_Layers(t *testing.T) {
	loadTestGames(t, "GAME_1")

	var result map[string]any
	getJSON(t, "/api/games/GAME_1/bases/0", handleGetBase, "/api/games/{slot}/bases/{idx}", &result)
	assert.NotContains(t, result, "layers", "no layers without game graphics")

	baseBits = [][]byte{{}}
	t.Cleanup(func() { baseBits = nil })
	getJSON(t, "/api/games/GAME_1/bases/0", handleGetBase, "/api/games/{slot}/bases/{idx}", &result)
	require.Contains(t, result, "layers")
	assert.Len(t, result["layers"], 4)
}
