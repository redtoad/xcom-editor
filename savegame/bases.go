package savegame

import (
	"fmt"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"log"
	"math"
	"path"
)

type Base struct {
	offset int
	game   *Savegame
}

// Index returns the base's index in BASE.DAT.
func (b *Base) Index() int {
	return b.offset
}

// Active returns whether the base is active.
func (b *Base) Active() bool {
	return b.game.BasesData.Bases[b.offset].Active
}

// Engineers returns the number of engineers at this base.
func (b *Base) Engineers() int {
	return b.game.BasesData.Bases[b.offset].Engineers
}

// Scientists returns the number of scientists at this base.
func (b *Base) Scientists() int {
	return b.game.BasesData.Bases[b.offset].Scientists
}

// Inventory returns the inventory quantities for this base.
func (b *Base) Inventory() [96]int {
	return b.game.BasesData.Bases[b.offset].Inventory
}

func (b *Base) Name() string {
	return b.game.BasesData.Bases[b.offset].Name.String()
}

func (b *Base) Coord() Coord {
	locationRefs := b.game.locationsForType(geoscape.XCOMBase)
	for _, ref := range locationRefs {
		if ref.Ref == b.offset {
			return ref.Coord()
		}
	}
	return Coord{-1, -1}
}

func (b *Base) Tiles() []BaseTile {
	data := b.game.BasesData.Bases[b.offset]
	tiles := make([]BaseTile, 36)
	for i := 0; i < 36; i++ {
		tiles[i] = BaseTile{
			Type:             data.Grid[i],
			DaysToCompletion: int(data.DaysToCompletion[i]),
		}
	}
	return tiles
}

func (b *Base) TileAt(x, y int) BaseTile {
	data := b.game.BasesData.Bases[b.offset]
	tileNo := x + y*6
	return BaseTile{
		Type:             data.Grid[tileNo],
		DaysToCompletion: int(data.DaysToCompletion[tileNo]),
	}
}

// Radar describes the detection capability of a radar facility.
type Radar struct {
	Facility geoscape.Facility
	// Range is the detection range in nautical miles.
	Range int
	// Chance is the chance (in percent) to detect a UFO within range at
	// each detection check.
	Chance int
}

// globeRadius is the radius of the globe in the units GEOSCAPE.EXE uses for
// distances (the amplitude of SINE.DAT and COS.DAT).
const globeRadius = 1024

// rangeToNauticalMiles converts a distance as measured by GEOSCAPE.EXE to
// nautical miles. The game measures the straight distance (chord) between
// two positions on the globe, a nautical mile is one arc minute.
func rangeToNauticalMiles(units float64) int {
	angle := 2 * math.Asin(units/(2*globeRadius)) // radians
	return int(math.Round(angle * 180 / math.Pi * 60))
}

// radars holds the detection capabilities of the radar facilities, ordered
// by range (longest first).
//
// The values are taken from the UFO detection routine of GEOSCAPE.EXE: a
// UFO is detected for sure within 800 units of a base with a hyperwave
// decoder, otherwise the base's long range chance applies within 750 units
// and its short range chance within 500 units. The chances are calculated
// when a construction completes: short range is 10 for small and 20 for
// large radars, long range is 20 for large radars and hyperwave is 100,
// facilities of the same type do not add up. (The in-game UFOpaedia states
// 300 and 450 nautical miles, other sources 1500, 2250 and 2400.)
var radars = []Radar{
	{Facility: geoscape.HyperwaveDecoder, Range: rangeToNauticalMiles(800), Chance: 100},
	{Facility: geoscape.LargeRadarSystem, Range: rangeToNauticalMiles(750), Chance: 20},
	{Facility: geoscape.SmallRadarSystem, Range: rangeToNauticalMiles(500), Chance: 10},
}

// Radars returns the detection capabilities of all completed radar
// facilities of this base, ordered by range (longest first). Multiple
// facilities of the same type are listed only once.
//
// The capabilities are derived from the facilities, since the detection
// values stored in BASE.DAT are only updated when a construction completes.
func (b *Base) Radars() []Radar {
	built := map[geoscape.Facility]bool{}
	for _, tile := range b.Tiles() {
		if tile.DaysToCompletion == 0 {
			built[tile.Type] = true
		}
	}
	result := make([]Radar, 0)
	for _, r := range radars {
		if built[r.Facility] {
			result = append(result, r)
		}
	}
	return result
}

type BaseTile struct {
	Type             geoscape.Facility
	DaysToCompletion int
}

func (game *Savegame) loadBases() error {
	filePath := path.Join(game.Path, "BASE.DAT")
	if err := internal.LoadDATFile(filePath, &game.BasesData); err != nil {
		return fmt.Errorf("could not load BASE.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveBases() error {
	filePath := path.Join(game.Path, "BASE.DAT")
	return internal.SaveDATFile(filePath, game.BasesData)
}

// CompleteConstructions will complete all ongoing constructions in all BasesData.
func (game *Savegame) CompleteConstructions() {
	for b := 0; b < len(game.BasesData.Bases); b++ {
		base := &game.BasesData.Bases[b]
		for i := 0; i < len(base.DaysToCompletion); i++ {
			if base.DaysToCompletion[i] > 0 {
				log.Printf("Complete construction of %v in %s.\n", base.Grid[i].String(), base.Name)
				base.DaysToCompletion[i] = 0
			}
		}
	}
}

// Base returns the base at the given index in BASE.DAT, or nil if there is
// no such base.
func (game *Savegame) Base(offset int) *Base {
	if offset < 0 || offset >= len(game.BasesData.Bases) {
		return nil
	}
	base := game.BasesData.Bases[offset]
	if base.Name.String() != "" {
		return &Base{
			offset: offset,
			game:   game,
		}
	}
	return nil
}

func (game *Savegame) Bases() []*Base {
	bases := make([]*Base, 0)
	for idx, base := range game.BasesData.Bases {
		if base.Name.String() != "" {
			bases = append(bases, game.Base(idx))
		}
	}
	return bases
}

func (game *Savegame) SetInventory(baseNr int, item geoscape.Inventory, amount int) {

}
