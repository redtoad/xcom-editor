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

// data returns a pointer to the underlying soldier data for direct access.
func (s *Soldier) data() *geoscape.SoldierData {
	return &s.game.soldierFile.Soldiers[s.offset]
}

// SetName sets the soldier's name.
func (s *Soldier) SetName(name string) {
	s.game.soldierFile.Soldiers[s.offset].Name = internal.NewNullString(name, 25)
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

// SetRecoveryDays sets the number of wound recovery days remaining.
func (s *Soldier) SetRecoveryDays(days int) {
	s.game.soldierFile.Soldiers[s.offset].RecoveryDays = days
}

// TimeUnits returns total time units (initial + improvement).
func (s *Soldier) TimeUnits() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialTimeUnits) + int(d.TimeUnitImprovement)
}

// Health returns total health (initial + improvement).
func (s *Soldier) Health() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialHealth) + int(d.HealthImprovement)
}

// Energy returns total energy/stamina (initial + improvement).
func (s *Soldier) Energy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialEnergy) + int(d.EnergyImprovement)
}

// Reactions returns total reactions (initial + improvement).
func (s *Soldier) Reactions() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialReactions) + int(d.ReactionsImprovement)
}

// Strength returns total strength (initial + improvement).
func (s *Soldier) Strength() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialStrength) + int(d.StrengthImprovement)
}

// FiringAccuracy returns total firing accuracy (initial + improvement).
func (s *Soldier) FiringAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialFiringAccuracy) + int(d.FiringAccuracyImprovement)
}

// ThrowingAccuracy returns total throwing accuracy (initial + improvement).
func (s *Soldier) ThrowingAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialThrowingAccuracy) + int(d.ThrowingAccuracyImprovement)
}

// MeleeAccuracy returns total melee accuracy (initial + improvement).
func (s *Soldier) MeleeAccuracy() int {
	d := s.game.soldierFile.Soldiers[s.offset]
	return int(d.InitialMeleeAccuracy) + int(d.MeleeAccuracyImprovement)
}

// PsionicStrength returns psionic strength (never changes).
func (s *Soldier) PsionicStrength() int {
	return int(s.game.soldierFile.Soldiers[s.offset].InitialPsionicStrength)
}

// PsionicSkill returns current psionic skill.
func (s *Soldier) PsionicSkill() int {
	return int(s.game.soldierFile.Soldiers[s.offset].InitialPsionicSkill)
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

// SetArmor sets the soldier's armor type.
func (s *Soldier) SetArmor(armor geoscape.Armor) {
	s.game.soldierFile.Soldiers[s.offset].Armor = armor
}

// Gender returns the soldier's gender.
func (s *Soldier) Gender() geoscape.Gender {
	return s.data().Sex
}

// Appearance returns the soldier's appearance.
func (s *Soldier) Appearance() geoscape.Appearance {
	return s.data().Appearance
}

// InitialTimeUnits returns the soldier's initial time units (before improvements).
func (s *Soldier) InitialTimeUnits() int { return int(s.data().InitialTimeUnits) }

// InitialHealth returns the soldier's initial health (before improvements).
func (s *Soldier) InitialHealth() int { return int(s.data().InitialHealth) }

// InitialEnergy returns the soldier's initial energy (before improvements).
func (s *Soldier) InitialEnergy() int { return int(s.data().InitialEnergy) }

// InitialReactions returns the soldier's initial reactions (before improvements).
func (s *Soldier) InitialReactions() int { return int(s.data().InitialReactions) }

// InitialStrength returns the soldier's initial strength (before improvements).
func (s *Soldier) InitialStrength() int { return int(s.data().InitialStrength) }

// InitialFiringAccuracy returns the soldier's initial firing accuracy (before improvements).
func (s *Soldier) InitialFiringAccuracy() int { return int(s.data().InitialFiringAccuracy) }

// InitialThrowingAccuracy returns the soldier's initial throwing accuracy (before improvements).
func (s *Soldier) InitialThrowingAccuracy() int { return int(s.data().InitialThrowingAccuracy) }

// InitialMeleeAccuracy returns the soldier's initial melee accuracy (before improvements).
func (s *Soldier) InitialMeleeAccuracy() int { return int(s.data().InitialMeleeAccuracy) }

// InitialPsionicStrength returns the soldier's initial psionic strength.
func (s *Soldier) InitialPsionicStrength() int { return int(s.data().InitialPsionicStrength) }

// InitialPsionicSkill returns the soldier's initial psionic skill.
func (s *Soldier) InitialPsionicSkill() int { return int(s.data().InitialPsionicSkill) }

// InitialBravery returns the soldier's initial bravery value.
func (s *Soldier) InitialBravery() int { return s.data().InitialBravery }

// SetInitialTimeUnits sets the soldier's initial time units.
func (s *Soldier) SetInitialTimeUnits(v int) { s.data().InitialTimeUnits = statByte(v) }

// SetInitialHealth sets the soldier's initial health.
func (s *Soldier) SetInitialHealth(v int) { s.data().InitialHealth = statByte(v) }

// SetInitialEnergy sets the soldier's initial energy.
func (s *Soldier) SetInitialEnergy(v int) { s.data().InitialEnergy = statByte(v) }

// SetInitialReactions sets the soldier's initial reactions.
func (s *Soldier) SetInitialReactions(v int) { s.data().InitialReactions = statByte(v) }

// SetInitialStrength sets the soldier's initial strength.
func (s *Soldier) SetInitialStrength(v int) { s.data().InitialStrength = statByte(v) }

// SetInitialFiringAccuracy sets the soldier's initial firing accuracy.
func (s *Soldier) SetInitialFiringAccuracy(v int) { s.data().InitialFiringAccuracy = statByte(v) }

// SetInitialThrowingAccuracy sets the soldier's initial throwing accuracy.
func (s *Soldier) SetInitialThrowingAccuracy(v int) { s.data().InitialThrowingAccuracy = statByte(v) }

// SetInitialMeleeAccuracy sets the soldier's initial melee accuracy.
func (s *Soldier) SetInitialMeleeAccuracy(v int) { s.data().InitialMeleeAccuracy = statByte(v) }

// SetInitialPsionicStrength sets the soldier's initial psionic strength.
func (s *Soldier) SetInitialPsionicStrength(v int) { s.data().InitialPsionicStrength = statByte(v) }

// SetInitialPsionicSkill sets the soldier's initial psionic skill.
func (s *Soldier) SetInitialPsionicSkill(v int) { s.data().InitialPsionicSkill = statByte(v) }

type Rank int

const (
	Rookie Rank = iota
	Squaddie
	Sergeant
	Captain
	Colonel
	Commander
	// Dead is the rank of soldiers killed in action, the game does not keep
	// their former rank (see LastRank).
	Dead Rank = -1
)

// Rank returns the soldier's rank, Dead for soldiers killed in action.
func (s *Soldier) Rank() Rank {
	if s.IsDead() {
		return Dead
	}
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

// LastRank returns the rank of a soldier before being killed in action (the
// current rank for living soldiers).
//
// The game overwrites the rank of dead soldiers, so it is estimated from the
// soldier value, which is 20 + missions + a rank bonus (Sergeant +1,
// Captain +3, Colonel +6, Commander +10, see geoscape.SoldierData). Rookies
// and Squaddies both have no bonus and are returned as Rookie. Edited
// savegames may have soldier values that do not follow this rule.
func (s *Soldier) LastRank() Rank {
	if !s.IsDead() {
		return s.Rank()
	}
	data := s.data()
	switch bonus := data.SoldierValue - 20 - data.Missions; {
	case bonus >= 10:
		return Commander
	case bonus >= 6:
		return Colonel
	case bonus >= 3:
		return Captain
	case bonus >= 1:
		return Sergeant
	}
	return Rookie
}

// Base returns the base the soldier is at, or nil if unknown (e.g. while the
// soldier is being transferred).
func (s *Soldier) Base() *Base {
	// the base reference is an index into LOC.DAT
	return s.game.baseAtLocation(s.game.soldierFile.Soldiers[s.offset].Base)
}

const NoCraft = -1

func (s *Soldier) Craft() *Craft {
	currCraftNo := s.game.soldierFile.Soldiers[s.offset].Craft
	prevCraftNo := s.game.soldierFile.Soldiers[s.offset].CraftBefore
	if currCraftNo > 0 {
		return s.game.Craft(currCraftNo - 1)
	}
	if prevCraftNo > 0 {
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

// Resurrect brings a dead soldier back to life, healthy and with the rank
// they had before (see LastRank). It returns false if the soldier is not
// dead.
func (s *Soldier) Resurrect() bool {
	if !s.IsDead() {
		return false
	}
	s.Heal()
	return true
}

// HealWounds lets a wounded soldier recover immediately and returns them to
// their original craft. It returns false if the soldier is not wounded (dead
// soldiers are not resurrected, see Resurrect).
func (s *Soldier) HealWounds() bool {
	if !s.IsWounded() {
		return false
	}
	s.Heal()
	return true
}

func (s *Soldier) Heal() {
	data := &s.game.soldierFile.Soldiers[s.offset]
	// Dead soldiers are resurrected with the rank they had before.
	if s.IsDead() {
		data.Rank = geoscape.Rank(s.LastRank())
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

// EquipPersonalArmour gives every living soldier without armour a personal
// armour. Soldiers already wearing better armour keep it. It returns the
// number of soldiers that were equipped.
func (game *Savegame) EquipPersonalArmour() int {
	n := 0
	for _, soldier := range game.Soldiers() {
		if soldier.IsDead() || soldier.Armor() != geoscape.NoArmor {
			continue
		}
		soldier.SetArmor(geoscape.PersonalArmour)
		n++
	}
	return n
}

// statByte converts a stat value to the unsigned byte stored in SOLDIER.DAT,
// clamping it to the valid range 0-255.
func statByte(v int) uint8 {
	return uint8(min(max(v, 0), 255))
}
