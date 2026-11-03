package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
)

// geoscapeFiles maps simple data files (which need no further processing) to
// the fields of the savegame they are loaded into.
func (game *Savegame) geoscapeFiles() map[string]any {
	return map[string]any{
		"ALIEN.DAT":  &game.AlienActivity,
		"XCOM.DAT":   &game.XcomActivity,
		"UIGLOB.DAT": &game.Counters,
		"ZONAL.DAT":  &game.ZonePriorities,
		"ACTS.DAT":   &game.AlienActs,
		"XBASES.DAT": &game.AlienBaseTargets,
	}
}

func (game *Savegame) loadGeoscapeFiles() error {
	for name, obj := range game.geoscapeFiles() {
		if err := internal.LoadDATFile(path.Join(game.Path, name), obj); err != nil {
			return fmt.Errorf("could not load %s: %w", name, err)
		}
	}
	return nil
}

func (game *Savegame) saveGeoscapeFiles() error {
	for name, obj := range game.geoscapeFiles() {
		if err := internal.SaveDATFile(path.Join(game.Path, name), obj); err != nil {
			return fmt.Errorf("could not save %s: %w", name, err)
		}
	}
	return nil
}
