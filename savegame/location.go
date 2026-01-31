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

// Location holds exported data for a single location entry from LOC.DAT.
type Location struct {
	Type           geoscape.LocationType
	TableReference int
	X, Y           int
	CountSuffix    int
}

// Coord converts the location's game coordinates to GPS coordinates.
func (loc *Location) Coord() Coord {
	return NewCoord(loc.X, loc.Y)
}

// Locations returns all active locations, skipping Unused and Waypoint entries.
func (game *Savegame) Locations() []Location {
	locs := make([]Location, 0)
	for _, obj := range game.locationFile.Objects {
		if obj.Type == geoscape.Unused || obj.Type == geoscape.Waypoint {
			continue
		}
		locs = append(locs, Location{
			Type:           obj.Type,
			TableReference: int(obj.TableReference),
			X:              obj.X,
			Y:              obj.Y,
			CountSuffix:    obj.CountSuffix,
		})
	}
	return locs
}
