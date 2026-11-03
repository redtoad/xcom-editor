package geoscape_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/require"
)

// savegameDirs are real savegames (geoscape) stored in testdata/ at the
// repository root.
var savegameDirs = []string{"GAME_1", "GAME_2", "GAME_3"}

// MustLoadTestFile loads a DAT file from a savegame in testdata/.
func MustLoadTestFile(t *testing.T, game, name string) []byte {
	t.Helper()
	buf, err := os.ReadFile(filepath.Join("..", "..", "testdata", game, name))
	require.NoError(t, err)
	return buf
}

// TestSavegameFilesRoundTrip loads every supported file of all real savegames
// and checks that packing it again results in identical bytes.
func TestSavegameFilesRoundTrip(t *testing.T) {
	files := []struct {
		name string
		obj  func() any
	}{
		{"ACTS.DAT", func() any { return &geoscape.ActsFile{} }},
		{"ALIEN.DAT", func() any { return &geoscape.AlienFile{} }},
		{"BASE.DAT", func() any { return &geoscape.BaseFile{} }},
		{"CRAFT.DAT", func() any { return &geoscape.CraftFile{} }},
		{"IGLOB.DAT", func() any { return &geoscape.IglobFile{} }},
		{"INTER.DAT", func() any { return &geoscape.InterFile{} }},
		{"LIGLOB.DAT", func() any { return &geoscape.LiglobFile{} }},
		{"LOC.DAT", func() any { return &geoscape.LocFile{} }},
		{"SAVEINFO.DAT", func() any { return &geoscape.SaveinfoFile{} }},
		{"SOLDIER.DAT", func() any { return &geoscape.SoldierFile{} }},
		{"TRANSFER.DAT", func() any { return &geoscape.TransferFile{} }},
		{"UIGLOB.DAT", func() any { return &geoscape.UiglobFile{} }},
		{"XBASES.DAT", func() any { return &geoscape.XbasesFile{} }},
		{"XCOM.DAT", func() any { return &geoscape.XcomFile{} }},
		{"ZONAL.DAT", func() any { return &geoscape.ZonalFile{} }},
	}

	for _, game := range savegameDirs {
		for _, f := range files {
			t.Run(game+"/"+f.name, func(t *testing.T) {
				buf := MustLoadTestFile(t, game, f.name)
				obj := f.obj()

				require.NoError(t, restruct.Unpack(buf, binary.LittleEndian, obj))
				packed, err := restruct.Pack(binary.LittleEndian, obj)
				require.NoError(t, err)
				if !bytes.Equal(buf, packed) {
					t.Errorf("re-packed %s does not match original data", f.name)
				}
			})
		}
	}
}
