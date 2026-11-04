package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/resources"
)

// inventoryPalette is the palette used for the inventory screen (see
// PALETTES.DAT).
const inventoryPalette = 4

// soldierImages holds the PNG encoded soldier images of the inventory
// screen (UFOGRAPH/MAN_*.SPK) by name (e.g. "MAN_0M1"), nil if the game files
// could not be loaded.
var soldierImages map[string][]byte

// soldierImageName returns the name of the inventory image (UFOGRAPH/*.SPK)
// that shows a soldier wearing armor, or "" if there is none.
//
// Without armour and with personal armour the image depends on the gender
// and appearance of the soldier, e.g. MAN_1F2 (personal armour, female,
// oriental). Power suits (MAN_2) and flying suits (MAN_3) look the same for
// everyone.
func soldierImageName(armor geoscape.Armor, gender geoscape.Gender, appearance geoscape.Appearance) string {
	switch armor {
	case geoscape.NoArmor, geoscape.PersonalArmour:
		if appearance < geoscape.Blonde || appearance > geoscape.African {
			return ""
		}
		sex := "M"
		if gender == geoscape.Female {
			sex = "F"
		}
		return fmt.Sprintf("MAN_%d%s%d", armor, sex, appearance)
	case geoscape.PowerSuit, geoscape.FlyingSuit:
		return fmt.Sprintf("MAN_%d", armor)
	}
	return ""
}

// soldierImageNames returns all inventory images of soldiers.
func soldierImageNames() []string {
	var names []string
	for _, armor := range []geoscape.Armor{geoscape.NoArmor, geoscape.PersonalArmour} {
		for _, gender := range []geoscape.Gender{geoscape.Male, geoscape.Female} {
			for appearance := geoscape.Blonde; appearance <= geoscape.African; appearance++ {
				names = append(names, soldierImageName(armor, gender, appearance))
			}
		}
	}
	return append(names, soldierImageName(geoscape.PowerSuit, 0, 0), soldierImageName(geoscape.FlyingSuit, 0, 0))
}

// soldierImagesFor returns the image names of a soldier for each armour
// (by its display name, see armorToString), or nil if no images are loaded.
func soldierImagesFor(gender geoscape.Gender, appearance geoscape.Appearance) map[string]string {
	if soldierImages == nil {
		return nil
	}
	result := make(map[string]string)
	for armor, label := range armorToString {
		if name := soldierImageName(armor, gender, appearance); name != "" {
			result[label] = name
		}
	}
	return result
}

// loadSoldierImages loads the soldier images of the inventory screen
// (UFOGRAPH/MAN_*.SPK) from the game root folder. The images are cropped to
// the soldier, leaving out the transparent area around it.
func loadSoldierImages(root string) (map[string][]byte, error) {
	palettes, err := resources.LoadPalettes(filepath.Join(root, "GEODATA", "PALETTES.DAT"))
	if err != nil {
		return nil, fmt.Errorf("could not load PALETTES.DAT: %w", err)
	}
	result := make(map[string][]byte)
	for _, name := range soldierImageNames() {
		spk, err := resources.LoadSPK(filepath.Join(root, "UFOGRAPH", name+".SPK"))
		if err != nil {
			return nil, fmt.Errorf("could not load %s.SPK: %w", name, err)
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, cropTransparent(spk.Image(palettes[inventoryPalette]))); err != nil {
			return nil, fmt.Errorf("could not encode %s: %w", name, err)
		}
		result[name] = buf.Bytes()
	}
	return result, nil
}

// cropTransparent returns the smallest part of img that contains all of its
// non-transparent pixels (or img itself if it is completely transparent).
func cropTransparent(img image.Image) image.Image {
	b := img.Bounds()
	crop := image.Rectangle{Min: b.Max, Max: b.Min}
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			if _, _, _, a := img.At(x, y).RGBA(); a == 0 {
				continue
			}
			crop.Min.X, crop.Min.Y = min(crop.Min.X, x), min(crop.Min.Y, y)
			crop.Max.X, crop.Max.Y = max(crop.Max.X, x+1), max(crop.Max.Y, y+1)
		}
	}
	if crop.Empty() {
		return img
	}
	sub, ok := img.(interface {
		SubImage(image.Rectangle) image.Image
	})
	if !ok {
		return img
	}
	return sub.SubImage(crop)
}

// GET /api/soldier-images/{name}.png
func handleGetSoldierImage(w http.ResponseWriter, r *http.Request) {
	buf, found := soldierImages[mux.Vars(r)["name"]]
	if !found {
		writeError(w, http.StatusNotFound, "image not found")
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=86400")
	_, _ = w.Write(buf)
}
