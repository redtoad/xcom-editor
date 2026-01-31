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

// Index returns the soldier's index in SOLDIER.DAT.
func (s *Soldier) Index() int {
	return s.offset
}

func (s *Soldier) Name() string {
	return s.game.soldierFile.Soldiers[s.offset].Name.String()
}

// Data returns a pointer to the underlying soldier data for direct access.
func (s *Soldier) Data() *geoscape.SoldierData {
	return &s.game.soldierFile.Soldiers[s.offset]
}

// SetName sets the soldier's name.
func (s *Soldier) SetName(name string) {
	s.game.soldierFile.Soldiers[s.offset].Name = internal.NullString(name)
}

// Missions returns the number of missions completed.
func (s *Soldier) Missions() int {
	return s.game.soldierFile.Soldiers[s.offset].Missions
}

// Kills returns the number of kills.
func (s *Soldier) Kills() int {
	return s.game.soldierFile.Soldiers[s.offset].Kills
}

// RecoveryDays returns the number of wound recovery days remaining.
func (s *Soldier) RecoveryDays() int {
	return s.game.soldierFile.Soldiers[s.offset].RecoveryDays
}

// TimeUnits returns total time units (initial + improvement).
func (s *Soldier) TimeUnits() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialTimeUnits + d.TimeUnitImprovement
}

// Health returns total health (initial + improvement).
func (s *Soldier) Health() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialHealth + d.HealthImprovement
}

// Energy returns total energy/stamina (initial + improvement).
func (s *Soldier) Energy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialEnergy + d.EnergyImprovement
}

// Reactions returns total reactions (initial + improvement).
func (s *Soldier) Reactions() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialReactions + d.ReactionsImprovement
}

// Strength returns total strength (initial + improvement).
func (s *Soldier) Strength() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialStrength + d.StrengthImprovement
}

// FiringAccuracy returns total firing accuracy (initial + improvement).
func (s *Soldier) FiringAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialFiringAccuracy + d.FiringAccuracyImprovement
}

// ThrowingAccuracy returns total throwing accuracy (initial + improvement).
func (s *Soldier) ThrowingAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialThrowingAccuracy + d.ThrowingAccuracyImprovement
}

// MeleeAccuracy returns total melee accuracy (initial + improvement).
func (s *Soldier) MeleeAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return d.InitialMeleeAccuracy + d.MeleeAccuracyImprovement
}

// PsionicStrength returns psionic strength (never changes).
func (s *Soldier) PsionicStrength() int {
	return s.game.soldierFile.Soldiers[s.offset].InitialPsionicStrength
}

// PsionicSkill returns current psionic skill.
func (s *Soldier) PsionicSkill() int {
	return s.game.soldierFile.Soldiers[s.offset].InitialPsionicSkill
}

// Bravery returns bravery value (110 - 10*(initial - improvement)).
func (s *Soldier) Bravery() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return 110 - 10*(d.InitialBravery-d.BraveryImprovement)
}

// Armor returns the soldier's current armor type.
func (s *Soldier) Armor() geoscape.Armor {
	return s.game.soldierFile.Soldiers[s.offset].Armor
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

const NoCraft = -1 //0xffff

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
	return data.Rank == geoscape.DeadOrUnused && strings.TrimSpace(data.Name.String()) != ""
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
