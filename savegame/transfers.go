package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
)

func (game *Savegame) loadTransfers() error {
	filePath := path.Join(game.Path, "TRANSFER.DAT")
	if err := internal.LoadDATFile(filePath, &game.transferFile); err != nil {
		return fmt.Errorf("could not load TRANSFER.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveTransfers() error {
	filePath := path.Join(game.Path, "TRANSFER.DAT")
	return internal.SaveDATFile(filePath, game.transferFile)
}

// SpeedupDelivery will reduce delivery time for all outstanding deliveries to 1 hour.
func (game *Savegame) SpeedupDelivery() {
	for no := 0; no < len(game.transferFile.Transfers); no++ {
		transfer := &game.transferFile.Transfers[no]
		if transfer.HoursLeft > 0 {
			transfer.HoursLeft = 1
		}
	}
}
