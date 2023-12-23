package savegame

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type Savegame struct {
	Path          string
	MetaData      geoscape.SavegameInfo
	FinancialData geoscape.LIGLOB_DAT
	BasesData     geoscape.BASE_DAT
	SoldierData   geoscape.SOLDIER_DAT
	TransferData  geoscape.TRANSFER_DAT
}

// loadFile loads a single data file from disk. name is the path
// inside the game directory.
func (game Savegame) loadFile(name string, obj interface{}) error {
	fp := path.Join(game.Path, name)
	buf, err := os.ReadFile(fp)
	if err != nil {
		return err
	}
	if err = restruct.Unpack(buf, binary.LittleEndian, obj); err != nil {
		return err
	}
	return nil
}

// saveFile saves a single data to its original location on disk file
// if the content has changed. name is the path inside the game directory.
func (game Savegame) saveFile(name string, obj interface{}) error {

	saveData, err := restruct.Pack(binary.LittleEndian, obj)
	if err != nil {
		return err
	}

	fp := path.Join(game.Path, name)
	original, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	if bytes.Equal(saveData, original) {
		log.Printf("not saving %s. content has not changed.\n", fp)
		return nil
	}
	if err = os.WriteFile(fp, saveData, os.ModePerm); err != nil {
		return fmt.Errorf("could not save file %s: %w", fp, err)
	}
	return nil
}

// Save saves the entire savegame on disk at its original location.
func (game Savegame) Save() error {
	files := map[string]interface{}{
		"LIGLOB.DAT":   &game.FinancialData,
		"BASE.DAT":     &game.BasesData,
		"SOLDIER.DAT":  &game.SoldierData,
		"TRANSFER.DAT": &game.TransferData,
	}
	for name, obj := range files {
		log.Printf("saving %s...\n", name)
		if err := game.saveFile(name, obj); err != nil {
			return err
		}
	}
	return nil
}

// Heal will restore all soldiers back to health.
func (game Savegame) Heal() {
	for no := 0; no < len(game.SoldierData.Soldiers); no++ {
		soldier := &game.SoldierData.Soldiers[no]
		// resurrect solders
		if soldier.Rank == geoscape.DeadOrUnused && strings.TrimSpace(soldier.Name) != "" {
			fmt.Printf("Resurrect %s from the dead\n", soldier.Name)
			soldier.Rank = geoscape.Squaddie
		}
		soldier.RecoveryDays = 0
		if soldier.Craft == 0xffff {
			soldier.Craft = soldier.CraftBefore
			soldier.CraftBefore = 0xffff
		}
	}
}

// SpeedupDelivery will reducce delivery time for all outstanding deliveries to 1 hour.
func (game Savegame) SpeedupDelivery() {
	for no := 0; no < len(game.TransferData.Transfers); no++ {
		transfer := &game.TransferData.Transfers[no]
		if transfer.HoursLeft > 0 {
			transfer.HoursLeft = 1
		}
	}
}

// CompleteConstructions will complete all ongoing constructions in all bases.
func (game Savegame) CompleteConstructions() {
	for b := 0; b < len(game.BasesData.Bases); b++ {
		base := &game.BasesData.Bases[b]
		for i := 0; i < len(base.DaysToCompletion); i++ {
			if base.DaysToCompletion[i] > 0 {
				log.Printf("Complete construction of %v in %s.\n", base.Grid[i].Tile(), base.Name)
				base.DaysToCompletion[i] = 0
			}
		}
	}
}

// Load loads a savegame from disk. This includes loading all required data files
// one by one.
func Load(root string) (Savegame, error) {
	game := Savegame{Path: root}
	files := map[string]interface{}{
		"SAVEINFO.DAT": &game.MetaData,
		"LIGLOB.DAT":   &game.FinancialData,
		"BASE.DAT":     &game.BasesData,
		"SOLDIER.DAT":  &game.SoldierData,
		"TRANSFER.DAT": &game.TransferData,
	}
	for name, obj := range files {
		log.Printf("loading %s...\n", name)
		if err := game.loadFile(name, obj); err != nil {
			return game, fmt.Errorf("could not load file %s: %w", name, err)
		}
	}
	return game, nil
}

// LoadFile will load a data file into a struct.
// Deprecated: This will be replaced with savegame.Load()
func LoadFile(path string, obj interface{}) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = restruct.Unpack(buf, binary.LittleEndian, obj); err != nil {
		return err
	}
	return nil
}

// SaveFile will save a strct into a data file.
// Deprecated: This will be replaced with savegame.Savegame.Save()
func SaveFile(path string, obj interface{}) error {
	buf, err := restruct.Pack(binary.LittleEndian, &obj)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, buf, os.ModePerm); err != nil {
		return err
	}
	return nil
}
