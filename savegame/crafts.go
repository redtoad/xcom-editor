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

// Crafts returns all existing X-COM craft and UFOs, ordered by their index
// in CRAFT.DAT.
//
// Like the game, only crafts referenced by a location in LOC.DAT are
// considered to exist. Unused entries in CRAFT.DAT are not reliably marked
// as unused (some savegames contain zeroed entries, which would be read as
// Skyrangers).
func (game *Savegame) Crafts() []*Craft {
	used := make(map[int]bool)
	for _, obj := range game.locationFile.Objects {
		if obj.Type != geoscape.XCOMShip && obj.Type != geoscape.AlienShip {
			continue
		}
		used[int(obj.TableReference)] = true
	}
	crafts := make([]*Craft, 0)
	for idx, data := range game.craftsFile.Crafts {
		if !used[idx] || data.Type == geoscape.EntryNotUsed {
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

// Speed returns the craft's current speed.
func (c *Craft) Speed() int {
	return c.game.craftsFile.Crafts[c.offset].Speed
}

// Destination returns the GPS coordinate of the craft's destination,
// or nil if the craft is stationary (speed 0 or no destination set).
func (c *Craft) Destination() *Coord {
	data := c.game.craftsFile.Crafts[c.offset]
	if data.Speed == 0 || data.FlightMode == geoscape.NoDestination {
		return nil
	}
	idx := data.Destination
	if idx < 0 || idx >= len(c.game.locationFile.Objects) {
		return nil
	}
	dest := c.game.locationFile.Objects[idx]
	coord := NewCoord(dest.X, dest.Y)
	return &coord
}

// NextWaypoint returns the GPS coordinate of the craft's next UFO waypoint,
// or nil if the craft is stationary or has no waypoint set.
func (c *Craft) NextWaypoint() *Coord {
	data := c.game.craftsFile.Crafts[c.offset]
	if data.Speed == 0 || data.FlightMode == geoscape.NoDestination {
		return nil
	}
	if data.NextUFOWaypointLon == 0 && data.NextUFOWaypointLat == 0 {
		return nil
	}
	coord := NewCoord(data.NextUFOWaypointLon, data.NextUFOWaypointLat)
	return &coord
}

// MissionType returns the mission type string for this craft.
func (c *Craft) MissionType() string {
	return c.game.craftsFile.Crafts[c.offset].MissionType.String()
}

// MissionZone returns the mission zone string for this craft.
func (c *Craft) MissionZone() string {
	return c.game.craftsFile.Crafts[c.offset].MissionZone.String()
}

func (game *Savegame) Craft(offset int) *Craft {
	if offset < 0 || offset >= len(game.craftsFile.Crafts) {
		return nil
	}
	craft := game.craftsFile.Crafts[offset]
	if craft.Type == geoscape.EntryNotUsed {
		return nil
	}
	return &Craft{
		offset: offset,
		game:   game,
	}
}

// Altitude returns the craft's current altitude.
func (c *Craft) Altitude() geoscape.Altitude {
	return c.game.craftsFile.Crafts[c.offset].Altitude
}

// alienRaces maps CraftData.AlienRace to the race names.
var alienRaces = []string{"Sectoid", "Snakeman", "Ethereal", "Muton", "Floater", "Mixed"}

// AlienRace returns the alien race on board of a UFO, or "" for X-COM craft.
func (c *Craft) AlienRace() string {
	if c.Type() < geoscape.SmallScout {
		return ""
	}
	race := c.game.craftsFile.Crafts[c.offset].AlienRace
	if race < 0 || race >= len(alienRaces) {
		return ""
	}
	return alienRaces[race]
}

// DestinationLocation returns the location (LOC.DAT) the craft is heading
// to, or nil if the craft is stationary.
func (c *Craft) DestinationLocation() *Location {
	data := c.game.craftsFile.Crafts[c.offset]
	if data.Speed == 0 || data.FlightMode == geoscape.NoDestination {
		return nil
	}
	idx := data.Destination
	if idx < 0 || idx >= len(c.game.locationFile.Objects) {
		return nil
	}
	obj := c.game.locationFile.Objects[idx]
	return &Location{
		Type:           obj.Type,
		TableReference: int(obj.TableReference),
		X:              obj.X,
		Y:              obj.Y,
		CountSuffix:    obj.CountSuffix,
	}
}

// Damaged returns true if the craft exists (it is referenced by LOC.DAT) but
// its entry in CRAFT.DAT is (nearly) empty, e.g. because the file was
// overwritten with zeros. The values of such a craft are meaningless.
func (c *Craft) Damaged() bool {
	data, err := internal.Marshall(&c.game.craftsFile.Crafts[c.offset])
	if err != nil {
		return true
	}
	nonZero := 0
	for _, b := range data {
		if b != 0 {
			nonZero++
		}
	}
	// valid entries always contain several values (weapons, fuel, base,
	// location, ...), damaged ones only leftovers.
	return nonZero <= 2
}
