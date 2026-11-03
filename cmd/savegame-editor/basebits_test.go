package main

import (
	"bytes"
	"encoding/json"
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
	tiles[12].Type = geoscape.LaserDefense

	layers := baseLayers(tiles, []hangarCraft{{Type: geoscape.Interceptor}})
	require.Len(t, layers, 5)
	ground, shapes, corridors, interiors, crafts := layers[0], layers[1], layers[2], layers[3], layers[4]

	assert.Len(t, ground, 36)
	assert.Equal(t, baseSprite{spriteGround, 160, 160}, ground[35])

	assert.ElementsMatch(t, []baseSprite{
		{2, 0, 0},                                           // access lift (octagon)
		{1, 32, 0},                                          // living quarters (square)
		{1 + 3, 32, 32},                                     // laboratory under construction (dotted square)
		{9, 64, 0}, {10, 96, 0}, {11, 64, 32}, {12, 96, 32}, // hangar
		{3, 0, 64}, // laser defense (cross)
	}, shapes)

	assert.ElementsMatch(t, []baseSprite{
		{spriteFacilityStart, 0, 0},                                 // access lift
		{spriteFacilityStart + int(geoscape.LivingQuarters), 32, 0}, // living quarters
		{spriteFacilityStart + int(geoscape.LaserDefense), 0, 64},   // laser defense
	}, interiors, "no interior for hangars and facilities under construction")

	// no corridors within the hangar or to the laboratory under construction
	assert.ElementsMatch(t, []baseSprite{
		{spriteCorridorH, 16, 0}, // access lift - living quarters
		{spriteCorridorH, 48, 0}, // living quarters - hangar
	}, corridors)

	// the craft is drawn in the middle of the hangar
	assert.Equal(t, []baseSprite{{spriteCraftStart + int(geoscape.Interceptor), 64 + 18, 12}}, crafts)
}

func TestBaseLayers_Crafts(t *testing.T) {
	tiles := emptyBase()
	for _, hangar := range []int{0, 2, 4} {
		tiles[hangar].Type = geoscape.HangarTopLeft
		tiles[hangar+1].Type = geoscape.HangarTopRight
		tiles[hangar+6].Type = geoscape.HangarBottomLeft
		tiles[hangar+7].Type = geoscape.HangarBottomRight
	}
	tiles[2].DaysToCompletion = 3 // second hangar is under construction

	crafts := baseLayers(tiles, []hangarCraft{
		{Type: geoscape.Skyranger, Out: true}, // occupies the first hangar
		{Type: geoscape.Avenger},
		{Type: geoscape.Firestorm}, // no hangar left
	})[4]

	assert.Equal(t, []baseSprite{{spriteCraftStart + int(geoscape.Avenger), 128 + 18, 12}}, crafts)
}

func TestBaseLayers_EmptyBase(t *testing.T) {
	layers := baseLayers(emptyBase(), []hangarCraft{{Type: geoscape.Skyranger}})
	assert.Len(t, layers[0], 36)
	for _, layer := range layers[1:] {
		assert.Empty(t, layer)
	}
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
	assert.Len(t, result["layers"], 5)
}

// Empty layers must be encoded as [] (not null), otherwise the frontend
// cannot iterate over them.
func TestBaseLayers_EmptyLayersAreNotNull(t *testing.T) {
	buf, err := json.Marshal(baseLayers(emptyBase(), nil))
	require.NoError(t, err)
	assert.NotContains(t, string(buf), "null")
}
