package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/image/font/sfnt"
)

// writeFakeFonts writes (synthetic) font files with one glyph ('!') into root.
func writeFakeFonts(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEODATA"), 0755))
	small := make([]byte, 8*9)
	small[0] = 1
	big := make([]byte, 16*16)
	big[0] = 1
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEODATA", "SMALLSET.DAT"), small, 0644))
	require.NoError(t, os.WriteFile(filepath.Join(root, "GEODATA", "BIGLETS.DAT"), big, 0644))
}

func TestLoadFonts(t *testing.T) {
	root := t.TempDir()
	writeFakeFonts(t, root)

	result, err := loadFonts(root)
	require.NoError(t, err)
	require.Contains(t, result, "small")
	require.Contains(t, result, "big")

	f, err := sfnt.Parse(result["big"])
	require.NoError(t, err)
	name, err := f.Name(nil, sfnt.NameIDFamily)
	require.NoError(t, err)
	assert.Equal(t, "X-COM big", name)

	_, err = loadFonts(t.TempDir())
	assert.Error(t, err, "missing game files")
}

func TestHandleGetFont(t *testing.T) {
	fonts = map[string][]byte{"small": []byte("font")}
	t.Cleanup(func() { fonts = map[string][]byte{} })

	r := mux.NewRouter()
	r.HandleFunc("/api/fonts/{name}.ttf", handleGetFont)
	for path, status := range map[string]int{
		"/api/fonts/small.ttf": http.StatusOK,
		"/api/fonts/big.ttf":   http.StatusNotFound,
	} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		assert.Equal(t, status, rec.Code, path)
		if status == http.StatusOK {
			assert.Equal(t, "font/ttf", rec.Header().Get("Content-Type"))
			assert.Equal(t, "font", rec.Body.String())
		}
	}
}
