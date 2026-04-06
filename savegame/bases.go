package savegame

import (
	"fmt"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"log"
	"path"
)

type Base struct {
	offset int
	game   *Savegame
}

func (b *Base) Name() string {
	return b.game.baseFile.Bases[b.offset].Name.String()
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
	data := b.game.baseFile.Bases[b.offset]
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
	data := b.game.baseFile.Bases[b.offset]
	tileNo := x + y*6
	return BaseTile{
		Type:             data.Grid[tileNo],
		DaysToCompletion: int(data.DaysToCompletion[tileNo]),
	}
}

type BaseTile struct {
	Type             geoscape.Facility
	DaysToCompletion int
}

func (game *Savegame) loadBases() error {
	filePath := path.Join(game.Path, "BASE.DAT")
	if err := internal.LoadDATFile(filePath, &game.baseFile); err != nil {
		return fmt.Errorf("could not load BASE.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveBases() error {
	filePath := path.Join(game.Path, "BASE.DAT")
	return internal.SaveDATFile(filePath, game.baseFile)
}

// CompleteConstructions will complete all ongoing constructions in all baseFile.
func (game *Savegame) CompleteConstructions() {
	for b := 0; b < len(game.baseFile.Bases); b++ {
		base := &game.baseFile.Bases[b]
		for i := 0; i < len(base.DaysToCompletion); i++ {
			if base.DaysToCompletion[i] > 0 {
				log.Printf("Complete construction of %v in %s.\n", base.Grid[i].String(), base.Name)
				base.DaysToCompletion[i] = 0
			}
		}
	}
}

func (game *Savegame) Base(offset int) *Base {
	base := game.baseFile.Bases[offset]
	if len(base.Name) > 0 {
		return &Base{
			offset: offset,
			game:   game,
		}
	}
	return nil
}

func (game *Savegame) Bases() []*Base {
	bases := make([]*Base, 0)
	for idx, base := range game.baseFile.Bases {
		if len(base.Name) > 0 {
			bases = append(bases, game.Base(idx))
		}
	}
	return bases
}

func (game *Savegame) SetInventory(baseNr int, item geoscape.Inventory, amount int) {

}
