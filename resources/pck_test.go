package resources

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePCK(t *testing.T) {
	// two sprites of width 2: the second one skips a row, then has a
	// transparent pixel and two colored pixels
	data := []byte{
		0, 1, 2, 3, 4, 0xFF,
		1, 0xFE, 1, 5, 6, 0xFF,
	}
	first, err := ParsePCK(data, 2, 0)
	require.NoError(t, err)
	assert.Equal(t, []uint{1, 2, 3, 4}, first.pixels)
	assert.Equal(t, 2, first.Height())

	second, err := ParsePCK(data, 2, 6)
	require.NoError(t, err)
	assert.Equal(t, []uint{0, 0, 0, 5, 6}, second.pixels)

	_, err = ParsePCK(data, 2, 100)
	assert.Error(t, err)
	_, err = ParsePCK([]byte{0, 1}, 2, 0)
	assert.Error(t, err, "missing end of sprite")

	// LoadPCK reads the same data from a file
	path := filepath.Join(t.TempDir(), "TEST.PCK")
	require.NoError(t, os.WriteFile(path, data, 0644))
	loaded, err := LoadPCK(path, 2, 6)
	require.NoError(t, err)
	assert.Equal(t, second, loaded)
}
