package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type LocationRef struct {
	X, Y int
	Ref  int
}

func (loc *LocationRef) Coord() Coord {
	return NewCoord(loc.X, loc.Y)
}

func (game *Savegame) loadLocations() error {
	filePath := path.Join(game.Path, "LOC.DAT")
	if err := internal.LoadDATFile(filePath, &game.locationFile); err != nil {
		return fmt.Errorf("could not load LOC.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) savelocations() error {
	filePath := path.Join(game.Path, "LOC.DAT")
	return internal.SaveDATFile(filePath, game.locationFile)
}

func (game *Savegame) locationsForType(type_ geoscape.LocationType) []*LocationRef {
	data := make([]*LocationRef, 0)
	for _, loc := range game.locationFile.Objects {
		if loc.Type == type_ {
			data = append(data, &LocationRef{
				loc.X, loc.Y, int(loc.TableReference),
			})
		}
	}
	return data
}
