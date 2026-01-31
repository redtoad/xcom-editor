package savegame

import (
	"fmt"
	"path"

	"github.com/redtoad/xcom-editor/internal"
)

// Transfer wraps a single transfer entry.
type Transfer struct {
	offset int
	game   *Savegame
}

// Index returns the transfer's index in TRANSFER.DAT.
func (t *Transfer) Index() int {
	return t.offset
}

// Origin returns the base index the item is coming from (255 = purchased).
func (t *Transfer) Origin() int {
	return int(t.game.transferFile.Transfers[t.offset].Origin)
}

// Destination returns the destination base index.
func (t *Transfer) Destination() int {
	return int(t.game.transferFile.Transfers[t.offset].Destination)
}

// HoursLeft returns the number of hours until delivery.
func (t *Transfer) HoursLeft() int {
	return int(t.game.transferFile.Transfers[t.offset].HoursLeft)
}

// Type returns the item type being transferred.
func (t *Transfer) Type() int {
	return int(t.game.transferFile.Transfers[t.offset].Type)
}

// ReferenceNumber returns the reference number (meaning depends on Type).
func (t *Transfer) ReferenceNumber() int {
	return t.game.transferFile.Transfers[t.offset].ReferenceNumber
}

// Quantity returns the number of items being transferred.
func (t *Transfer) Quantity() int {
	return int(t.game.transferFile.Transfers[t.offset].Quantity)
}

// Transfers returns all active transfers (where Quantity > 0).
func (game *Savegame) Transfers() []*Transfer {
	transfers := make([]*Transfer, 0)
	for idx, data := range game.transferFile.Transfers {
		if data.Quantity > 0 {
			transfers = append(transfers, &Transfer{
				offset: idx,
				game:   game,
			})
		}
	}
	return transfers
}

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
