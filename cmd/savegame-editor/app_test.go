package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// copySlot copies a savegame from testdata/ (at the repository root) into
// root.
func copySlot(t *testing.T, root, slot string) {
	t.Helper()
	src := filepath.Join("..", "..", "testdata", slot)
	dst := filepath.Join(root, slot)
	require.NoError(t, os.MkdirAll(dst, 0755))
	entries, err := os.ReadDir(src)
	require.NoError(t, err)
	for _, e := range entries {
		buf, err := os.ReadFile(filepath.Join(src, e.Name()))
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(filepath.Join(dst, e.Name()), buf, 0644))
	}
}

func TestLoadRoot(t *testing.T) {
	root := t.TempDir()
	copySlot(t, root, "GAME_1")
	writeFakeGameFiles(t, root)

	require.NoError(t, loadRoot(root))
	assert.Len(t, games, 1)
	assert.NotNil(t, world)

	// a folder without savegames keeps the current state
	assert.Error(t, loadRoot(t.TempDir()))
	assert.Len(t, games, 1)
	assert.NotNil(t, world)
}

// TestLoadRoot_ConcurrentRequests switches between a folder with and one
// without game files while API requests are running (run with -race).
func TestLoadRoot_ConcurrentRequests(t *testing.T) {
	withGameFiles := t.TempDir()
	copySlot(t, withGameFiles, "GAME_1")
	writeFakeGameFiles(t, withGameFiles)
	savegamesOnly := t.TempDir()
	copySlot(t, savegamesOnly, "GAME_2")
	require.NoError(t, loadRoot(withGameFiles))
	t.Cleanup(func() { world = nil; baseBits = nil; fonts = map[string][]byte{} })

	r, err := newRouter()
	require.NoError(t, err)

	var wg sync.WaitGroup
	done := make(chan struct{})
	for _, path := range []string{"/api/games", "/api/world", "/api/world/textures/1/0.png", "/api/basebits/0.png"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-done:
					return
				default:
				}
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
			}
		}()
	}
	for i := 0; i < 20; i++ {
		require.NoError(t, loadRoot(savegamesOnly))
		require.NoError(t, loadRoot(withGameFiles))
	}
	close(done)
	wg.Wait()
}
