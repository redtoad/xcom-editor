package savegame

//go:generate stringer -type=Rank -output=soldiers_string.go -linecomment

import (
	"fmt"
	"path"
	"strings"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

type Soldier struct {
	game   *Savegame
	offset int // offset in file SOLDIER.DAT
}

func (s *Soldier) Name() string {
	return s.game.soldierFile.Soldiers[s.offset].Name.String()
}

type Rank int

const (
	Rookie Rank = iota
	Squaddie
	Sergeant
	Captain
	Colonel
	Commander
)

func (s *Soldier) Rank() Rank {
	mapping := map[geoscape.Rank]Rank{
		geoscape.Rookie:    Rookie,
		geoscape.Squaddie:  Squaddie,
		geoscape.Sergeant:  Sergeant,
		geoscape.Captain:   Captain,
		geoscape.Colonel:   Colonel,
		geoscape.Commander: Commander,
	}
	data := s.game.soldierFile.Soldiers[s.offset]
	return mapping[data.Rank]
}

func (s *Soldier) Base() *Base {
	baseNo := s.game.soldierFile.Soldiers[s.offset].Base
	return s.game.Base(baseNo)
}

const NoCraft = 0xffff

func (s *Soldier) Craft() *Craft {
	currCraftNo := s.game.soldierFile.Soldiers[s.offset].Craft
	prevCraftNo := s.game.soldierFile.Soldiers[s.offset].CraftBefore
	if currCraftNo != NoCraft {
		return s.game.Craft(currCraftNo - 1)
	}
	if prevCraftNo != NoCraft {
		return s.game.Craft(prevCraftNo - 1)
	}
	return nil
}

func (s *Soldier) IsDead() bool {
	data := s.game.soldierFile.Soldiers[s.offset]
	return data.Rank == geoscape.DeadOrUnused &&
		strings.TrimSpace(data.Name.String()) != ""
}

func (s *Soldier) IsWounded() bool {
	data := s.game.soldierFile.Soldiers[s.offset]
	return !s.IsDead() && data.RecoveryDays > 0
}

func (s *Soldier) Heal() {
	data := &s.game.soldierFile.Soldiers[s.offset]
	// Dead soldiers are resurrected as Squaddies.
	if s.IsDead() {
		data.Rank = geoscape.Squaddie
	}
	data.RecoveryDays = 0
	// Soldiers are returned to their original craft.
	if data.Craft == NoCraft {
		data.Craft = data.CraftBefore
		data.CraftBefore = NoCraft
	}
}

func (game *Savegame) loadSoldiers() error {
	filePath := path.Join(game.Path, "SOLDIER.DAT")
	if err := internal.LoadDATFile(filePath, &game.soldierFile); err != nil {
		return fmt.Errorf("could not load SOLDIER.DAT: %w", err)
	}
	return nil
}

func (game *Savegame) saveSoldiers() error {
	filePath := path.Join(game.Path, "SOLDIER.DAT")
	return internal.SaveDATFile(filePath, game.soldierFile)
}

func (game *Savegame) Soldiers() []*Soldier {
	soldiers := make([]*Soldier, 0)
	for idx, data := range game.soldierFile.Soldiers {
		// When soldiers die, their Rank is set to geoscape.DeadOrUnused any can
		// be overwritten. Therefore we assume that any non empty name marks a
		// soldier's entry.
		if data.Name.String() == "" {
			continue
		}
		soldiers = append(soldiers, &Soldier{
			offset: idx,
			game:   game,
		})
	}
	return soldiers
}

// HealAllSoldiers will restore all soldiers back to health.
func (game *Savegame) HealAllSoldiers() {
	for _, soldier := range game.Soldiers() {
		soldier.Heal()
	}
}
