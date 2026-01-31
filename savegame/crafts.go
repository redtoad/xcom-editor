package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type Craft struct {
	offset int
	game   *Savegame
}

func (c *Craft) Name() string {
	var suffix = 0
	for _, object := range c.game.locationFile.Objects {
		if object.Type == geoscape.XCOMShip && int(object.TableReference) == c.offset {
			suffix = object.CountSuffix
		}
	}
	prefixMapping := map[geoscape.CraftType]string{
		geoscape.Skyranger:   "SKYRANGER",
		geoscape.Lightning:   "LIGHTNING",
		geoscape.Avenger:     "AVENGER",
		geoscape.Interceptor: "INTERCEPTOR",
		geoscape.Firestorm:   "FIRESTORM",
		geoscape.SmallScout:  "UFO",
		geoscape.MediumScout: "UFO",
		geoscape.LargeScout:  "UFO",
		geoscape.Harvester:   "UFO",
		geoscape.Abductor:    "UFO",
		geoscape.TerrorShip:  "UFO",
		geoscape.Battleship:  "UFO",
		geoscape.SupplyShip:  "UFO",
	}
	prefix := prefixMapping[c.game.craftsFile.Crafts[c.offset].Type]
	return fmt.Sprintf("%s-%d", prefix, suffix)
}

func (game *Savegame) loadCrafts() error {
	filePath := path.Join(game.Path, "CRAFT.DAT")
	if err := internal.LoadDATFile(filePath, &game.craftsFile); err != nil {
		return fmt.Errorf("could not load CRAFT.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveCrafts() error {
	filePath := path.Join(game.Path, "CRAFT.DAT")
	return internal.SaveDATFile(filePath, game.craftsFile)
}

// Crafts returns all active craft (skipping EntryNotUsed slots).
func (game *Savegame) Crafts() []*Craft {
	crafts := make([]*Craft, 0)
	for idx, data := range game.craftsFile.Crafts {
		if data.Type == geoscape.EntryNotUsed {
			continue
		}
		crafts = append(crafts, &Craft{
			offset: idx,
			game:   game,
		})
	}
	return crafts
}

// Index returns the craft's index in CRAFT.DAT.
func (c *Craft) Index() int {
	return c.offset
}

// Type returns the craft type.
func (c *Craft) Type() geoscape.CraftType {
	return c.game.craftsFile.Crafts[c.offset].Type
}

// Status returns the craft status.
func (c *Craft) Status() geoscape.CraftStatus {
	return c.game.craftsFile.Crafts[c.offset].Status
}

// Damage returns the current damage amount.
func (c *Craft) Damage() int {
	return c.game.craftsFile.Crafts[c.offset].Damage
}

// Fuel returns the remaining fuel amount.
func (c *Craft) Fuel() int {
	return c.game.craftsFile.Crafts[c.offset].FuelAmountRemaining
}

// Base returns the base this craft is stationed at.
func (c *Craft) Base() *Base {
	baseRef := int(c.game.craftsFile.Crafts[c.offset].BaseReference)
	return c.game.Base(baseRef)
}

func (game *Savegame) Craft(offset int) *Craft {
	craft := game.craftsFile.Crafts[offset]
	if craft.Type == geoscape.EntryNotUsed {
		return nil
	}
	return &Craft{
		offset: offset,
		game:   game,
	}
}
