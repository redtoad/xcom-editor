package api

//go:generate go run github.com/99designs/gqlgen

import (
	"fmt"
	"github.com/redtoad/xcom-editor/lib/savegame"
	"os"
	"path"
	"strings"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RootPath string
}

func (r *Resolver) SavegameRoot(nr int) string {
	return path.Join(r.RootPath, fmt.Sprintf("GAME_%d", nr))
}

func (r *Resolver) LoadSavegame(nr int) (*Savegame, error) {

	rootPath := r.SavegameRoot(nr)

	pathSoldiersFile := path.Join(rootPath, "SOLDIER.DAT")
	if _, err := os.Stat(pathSoldiersFile); os.IsNotExist(err) {
		return nil, err
	}

	var soldierData savegame.FileSoldier
	if err := savegame.LoadFile(pathSoldiersFile, &soldierData); err != nil {
		return nil, err
	}

	var soldiers []*Soldier
	for no := 0; no < len(soldierData.Soldiers); no++ {
		data := &soldierData.Soldiers[no]
		if strings.TrimSpace(data.Name) != "" {
			soldiers = append(soldiers, &Soldier{
				Slot:     no,
				Name:     data.Name,
				Rank:     data.Rank.String(),
				Missions: data.Missions,
				Kills:    data.Kills,
			})
		}
	}

	return &Savegame{
		Soldiers: soldiers,
	}, nil
}
