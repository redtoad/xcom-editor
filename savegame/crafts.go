package savegame

import (
	"fmt"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"path"
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
