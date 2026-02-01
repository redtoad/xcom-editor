package savegame

import (
	"fmt"
	"github.com/redtoad/xcom-editor/internal"
	"path"
	"time"
)

func (game *Savegame) loadMetadata() error {
	filePath := path.Join(game.Path, "SAVEINFO.DAT")
	if err := internal.LoadDATFile(filePath, &game.meta); err != nil {
		return fmt.Errorf("could not load SAVEINFO.DAT: %w", err)
	}
	return nil
}

// Title returns the savegame title.
func (game *Savegame) Title() string {
	return game.meta.Name.String()
}

// Time returns the game time.
func (game *Savegame) Time() time.Time {
	return game.meta.Time()
}
