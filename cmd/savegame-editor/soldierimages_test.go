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
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSoldierImageName(t *testing.T) {
	for _, tt := range []struct {
		armor      geoscape.Armor
		gender     geoscape.Gender
		appearance geoscape.Appearance
		want       string
	}{
		{geoscape.NoArmor, geoscape.Male, geoscape.Blonde, "MAN_0M0"},
		{geoscape.NoArmor, geoscape.Female, geoscape.African, "MAN_0F3"},
		{geoscape.PersonalArmour, geoscape.Female, geoscape.Oriental, "MAN_1F2"},
		{geoscape.PersonalArmour, geoscape.Male, geoscape.BrownHair, "MAN_1M1"},
		{geoscape.PowerSuit, geoscape.Female, geoscape.Oriental, "MAN_2"},
		{geoscape.FlyingSuit, geoscape.Male, geoscape.Blonde, "MAN_3"},
		{geoscape.SectoidArmor, geoscape.Male, geoscape.Blonde, ""},
		{geoscape.NoArmor, geoscape.Male, geoscape.Appearance(7), ""},
	} {
		assert.Equal(t, tt.want, soldierImageName(tt.armor, tt.gender, tt.appearance))
	}
	assert.Len(t, soldierImageNames(), 18)
}

// spkFile returns an SPK image (320x200) with a single row of opaque pixels
// of colour 1 at x=100..103 in line y=50.
func spkFile() []byte {
	var buf bytes.Buffer
	w := func(v uint16) { _ = binary.Write(&buf, binary.LittleEndian, v) }
	w(0xffff)
	w((50*320 + 100) / 2) // skip pixels (in units of 2)
	w(0xfffe)
	w(2) // draw 4 pixels
	buf.Write([]byte{1, 1, 1, 1})
	w(0xffff)
	w((320*200 - 50*320 - 104) / 2)
	w(0xfffd)
	return buf.Bytes()
}

// writeFakeSoldierImages writes all soldier images (and a palette, unless
// there already is one) into root.
func writeFakeSoldierImages(t *testing.T, root string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "GEODATA"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "UFOGRAPH"), 0755))
	// keep the palettes of other fake game files
	palettes := filepath.Join(root, "GEODATA", "PALETTES.DAT")
	if _, err := os.Stat(palettes); os.IsNotExist(err) {
		require.NoError(t, os.WriteFile(palettes, make([]byte, 5*(768+6)), 0644))
	}
	for _, name := range soldierImageNames() {
		require.NoError(t, os.WriteFile(filepath.Join(root, "UFOGRAPH", name+".SPK"), spkFile(), 0644))
	}
}

func TestLoadSoldierImages(t *testing.T) {
	root := t.TempDir()
	writeFakeSoldierImages(t, root)

	images, err := loadSoldierImages(root)
	require.NoError(t, err)
	require.Len(t, images, 18)

	// the image is cropped to the opaque pixels
	img, err := png.Decode(bytes.NewReader(images["MAN_1F2"]))
	require.NoError(t, err)
	assert.Equal(t, 4, img.Bounds().Dx())
	assert.Equal(t, 1, img.Bounds().Dy())
}

func TestLoadSoldierImages_Missing(t *testing.T) {
	root := t.TempDir()
	writeFakeSoldierImages(t, root)
	require.NoError(t, os.Remove(filepath.Join(root, "UFOGRAPH", "MAN_3.SPK")))

	_, err := loadSoldierImages(root)
	assert.Error(t, err)
}

func TestHandleGetSoldierImage(t *testing.T) {
	root := t.TempDir()
	writeFakeSoldierImages(t, root)
	images, err := loadSoldierImages(root)
	require.NoError(t, err)
	soldierImages = images
	t.Cleanup(func() { soldierImages = nil })

	r := mux.NewRouter()
	r.HandleFunc("/api/soldier-images/{name}.png", handleGetSoldierImage)

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/soldier-images/MAN_0M1.png", nil))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "image/png", rec.Header().Get("Content-Type"))

	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/soldier-images/MAN_9.png", nil))
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestHandleGetSoldier_Images(t *testing.T) {
	loadTestGames(t, "GAME_1")

	type soldier struct {
		Index      int               `json:"index"`
		Gender     string            `json:"gender"`
		Appearance string            `json:"appearance"`
		Images     map[string]string `json:"images"`
	}

	// without game files there are no images
	var s soldier
	getJSON(t, "/api/games/GAME_1/soldiers/0", handleGetSoldier, "/api/games/{slot}/soldiers/{idx}", &s)
	assert.Nil(t, s.Images)
	assert.NotEmpty(t, s.Gender)
	assert.NotEmpty(t, s.Appearance)

	soldierImages = map[string][]byte{}
	t.Cleanup(func() { soldierImages = nil })
	getJSON(t, "/api/games/GAME_1/soldiers/0", handleGetSoldier, "/api/games/{slot}/soldiers/{idx}", &s)
	sg := games["GAME_1"].sg.Soldiers()[0]
	assert.Equal(t, map[string]string{
		"None":            soldierImageName(geoscape.NoArmor, sg.Gender(), sg.Appearance()),
		"Personal Armour": soldierImageName(geoscape.PersonalArmour, sg.Gender(), sg.Appearance()),
		"Power Suit":      "MAN_2",
		"Flying Suit":     "MAN_3",
	}, s.Images)
}
