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
	"github.com/redtoad/xcom-editor/internal/geodata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeFakeGameFiles writes minimal (synthetic) game files needed for the
// globe into root.
func writeFakeGameFiles(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEODATA"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEOGRAPH"), 0755))

	// WORLD.DAT: a quad and a triangle
	var world bytes.Buffer
	for _, p := range [][9]int32{
		{0, 0, 80, 0, 80, -80, 0, -80, 3},       // quad, terrain 3
		{1440, 0, 1520, 0, 1440, 80, -1, 0, 12}, // triangle, terrain 12
	} {
		for _, v := range p[:8] {
			require.NoError(t, binary.Write(&world, binary.LittleEndian, int16(v)))
		}
		require.NoError(t, binary.Write(&world, binary.LittleEndian, p[8]))
	}
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEODATA", "WORLD.DAT"), world.Bytes(), 0644))

	// TEXTURE.DAT: polar seas (set 0) mostly uses color 5
	textures := make([]byte, geodata.NumTextureSets*geodata.NumTerrains*geodata.TextureSize*geodata.TextureSize)
	start := polarSeas * geodata.TextureSize * geodata.TextureSize
	for i := 0; i < 600; i++ {
		textures[start+i] = 5
	}
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEOGRAPH", "TEXTURE.DAT"), textures, 0644))

	// PALETTES.DAT: 5 palettes, color 5 of palette 0 is (4, 8, 12) * 4
	palettes := make([]byte, 5*(768+6))
	copy(palettes[5*3:], []byte{4, 8, 12})
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEODATA", "PALETTES.DAT"), palettes, 0644))
}

func TestLoadWorld(t *testing.T) {
	root := t.TempDir()
	writeFakeGameFiles(t, root)

	w, err := loadWorld(root)
	require.NoError(t, err)

	require.Len(t, w.Polygons, 2)
	assert.Equal(t, worldPolygon{
		Terrain: 3,
		Points:  [][2]float32{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
	}, w.Polygons[0])
	assert.Equal(t, 12, w.Polygons[1].Terrain)
	assert.Len(t, w.Polygons[1].Points, 3, "triangle")
	assert.Equal(t, "#102030", w.Ocean)

	img, err := png.Decode(bytes.NewReader(w.Textures[1][polarSeas]))
	require.NoError(t, err)
	assert.Equal(t, geodata.TextureSize, img.Bounds().Dx())
}

func TestLoadWorld_MissingFiles(t *testing.T) {
	_, err := loadWorld(t.TempDir())
	assert.Error(t, err)
}

func TestHandleGetWorld(t *testing.T) {
	world = nil
	var result struct {
		Available bool           `json:"available"`
		Ocean     string         `json:"ocean"`
		Polygons  []worldPolygon `json:"polygons"`
	}
	getJSON(t, "/api/world", handleGetWorld, "/api/world", &result)
	assert.False(t, result.Available)

	root := t.TempDir()
	writeFakeGameFiles(t, root)
	w, err := loadWorld(root)
	require.NoError(t, err)
	world = w
	t.Cleanup(func() { world = nil })

	getJSON(t, "/api/world", handleGetWorld, "/api/world", &result)
	assert.True(t, result.Available)
	assert.Equal(t, "#102030", result.Ocean)
	assert.Len(t, result.Polygons, 2)
}

func TestHandleGetTexture(t *testing.T) {
	root := t.TempDir()
	writeFakeGameFiles(t, root)
	w, err := loadWorld(root)
	require.NoError(t, err)
	world = w
	t.Cleanup(func() { world = nil })

	r := mux.NewRouter()
	r.HandleFunc("/api/world/textures/{set}/{terrain}.png", handleGetTexture)
	for path, status := range map[string]int{
		"/api/world/textures/0/3.png":  http.StatusOK,
		"/api/world/textures/2/12.png": http.StatusOK,
		"/api/world/textures/3/0.png":  http.StatusNotFound,
		"/api/world/textures/0/13.png": http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, status, rec.Code, path)
		if status == http.StatusOK {
			assert.Equal(t, "image/png", rec.Header().Get("Content-Type"), path)
		}
	}
}
