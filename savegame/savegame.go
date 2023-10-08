package savegame

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"path"

	"github.com/go-restruct/restruct"
)

type Savegame struct {
	Path         string
	MetaData     SAVEINFO_DAT
	BasesData    BASE_DAT
	SoldiersData SOLDIER_DAT
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
		return err
	}
	return nil
}

// Save saves the entire savegame on disk at its original location.
func (game Savegame) Save() error {
	files := map[string]interface{}{
		"SAVEINFO.DAT": &game.MetaData,
		"BASE.DAT":     &game.BasesData,
		"SOLDIER.DAT":  &game.SoldiersData,
	}
	for name, obj := range files {
		log.Printf("saving %s...\n", name)
		if err := game.saveFile(name, obj); err != nil {
			return err
		}
	}
	return nil
}

// Load loads a savegame from disk. This includeds loading all data files
// individually.
func Load(root string) (Savegame, error) {
	game := Savegame{Path: root}
	files := map[string]interface{}{
		"SAVEINFO.DAT": &game.MetaData,
		"BASE.DAT":     &game.BasesData,
		"SOLDIER.DAT":  &game.SoldiersData,
	}
	for name, obj := range files {
		log.Printf("loading %s...\n", name)
		if err := game.loadFile(name, obj); err != nil {
			return game, err
		}
	}
	return game, nil
}

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
