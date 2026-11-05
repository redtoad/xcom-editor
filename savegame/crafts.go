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

// Base returns the base this craft is stationed at, or nil if unknown.
func (c *Craft) Base() *Base {
	// the base reference is an index into LOC.DAT
	idx := int(c.game.craftsFile.Crafts[c.offset].BaseReference)
	if idx < 0 || idx >= len(c.game.locationFile.Objects) {
		return nil
	}
	loc := c.game.locationFile.Objects[idx]
	if loc.Type != geoscape.XCOMBase {
		return nil
	}
	return c.game.Base(int(loc.TableReference))
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

// maxFuel maps the X-COM craft types to their fuel capacity.
var maxFuel = map[geoscape.CraftType]int{
	geoscape.Skyranger:   1500,
	geoscape.Lightning:   30,
	geoscape.Avenger:     60,
	geoscape.Interceptor: 1000,
	geoscape.Firestorm:   20,
}

// craftWeapon describes the ammunition of a craft weapon.
type craftWeapon struct {
	// maxAmmo is the number of rounds of a fully armed weapon.
	maxAmmo int
	// ammo is the item in the base stores used to rearm the weapon, or -1
	// if the weapon needs no ammunition (laser and plasma weapons).
	ammo geoscape.Inventory
	// roundsPerItem is the number of rounds one ammo item contains.
	roundsPerItem int
}

var craftWeapons = map[geoscape.WeaponType]craftWeapon{
	geoscape.Stingray:    {6, geoscape.InventoryStingrayMissile, 1},
	geoscape.Avalanche:   {3, geoscape.InventoryAvalancheMissile, 1},
	geoscape.Cannon:      {200, geoscape.InventoryCannonRounds, 50}, // Cannon Rounds (x50)
	geoscape.FusionBall:  {2, geoscape.InventoryFusionBalls, 1},
	geoscape.LaserCannon: {99, -1, 0},
	geoscape.PlasmaBeam:  {100, -1, 0},
}

// MaxFuel returns the fuel capacity of an X-COM craft, or 0 for UFOs.
func (c *Craft) MaxFuel() int {
	return maxFuel[c.Type()]
}

// RepairAndRefuel repairs all damage of an X-COM craft, fills up its fuel
// and rearms its weapons. Values above the capacity (e.g. from earlier edits)
// are kept. A craft waiting at its base for repairs, refuelling or rearming
// becomes ready. The ammunition is not taken from the base stores (see
// RepairAndRefuelAllCrafts). It returns false for UFOs and damaged entries
// (see Damaged), which are left unchanged.
func (c *Craft) RepairAndRefuel() bool {
	capacity := c.MaxFuel()
	if capacity == 0 || c.Damaged() {
		return false
	}
	data := &c.game.craftsFile.Crafts[c.offset]
	data.Damage = 0
	data.FuelAmountRemaining = max(data.FuelAmountRemaining, capacity)
	if w, ok := craftWeapons[data.LeftWeapon]; ok {
		data.LeftAmmo = max(data.LeftAmmo, w.maxAmmo)
	}
	if w, ok := craftWeapons[data.RightWeapon]; ok {
		data.RightAmmo = max(data.RightAmmo, w.maxAmmo)
	}
	switch data.Status {
	case geoscape.Repairs, geoscape.Refueling, geoscape.Rearming:
		data.Status = geoscape.Ready
	}
	return true
}

// reloadItems returns the ammo items (and their quantities) needed to fully
// rearm the craft's weapons once.
func (c *Craft) reloadItems() map[geoscape.Inventory]int {
	items := make(map[geoscape.Inventory]int)
	data := c.game.craftsFile.Crafts[c.offset]
	for _, weapon := range []geoscape.WeaponType{data.LeftWeapon, data.RightWeapon} {
		w, ok := craftWeapons[weapon]
		if !ok || w.ammo < 0 {
			continue
		}
		items[w.ammo] += (w.maxAmmo + w.roundsPerItem - 1) / w.roundsPerItem
	}
	return items
}

// RepairAndRefuelAllCrafts repairs, refuels and rearms all X-COM craft (see
// RepairAndRefuel) and returns the number of craft that were changed.
//
// The base stores are not used for rearming. Instead, every base gets enough
// ammunition in its stores to fully rearm all of its craft once more (items
// already in the stores count towards that), so the game can rearm the craft
// after their next mission.
func (game *Savegame) RepairAndRefuelAllCrafts() int {
	n := 0
	needed := make(map[int]map[geoscape.Inventory]int)
	for _, craft := range game.Crafts() {
		if !craft.RepairAndRefuel() {
			continue
		}
		n++
		base := craft.Base()
		if base == nil {
			continue
		}
		if needed[base.offset] == nil {
			needed[base.offset] = make(map[geoscape.Inventory]int)
		}
		for item, qty := range craft.reloadItems() {
			needed[base.offset][item] += qty
		}
	}
	for baseIdx, items := range needed {
		inventory := &game.BasesData.Bases[baseIdx].Inventory
		for item, qty := range items {
			inventory[item] = max(inventory[item], qty)
		}
	}
	return n
}
