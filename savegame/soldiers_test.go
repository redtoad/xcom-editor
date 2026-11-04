package savegame

import (
	"fmt"
	"testing"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

func TestSoldierCraft_NoCraftSentinel(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 1)
	game.craftsFile.Crafts[0].Type = geoscape.Avenger
	game.soldierFile.Soldiers[0].Craft = -1
	game.soldierFile.Soldiers[0].CraftBefore = -1

	s := &Soldier{game: game, offset: 0}
	if got := s.Craft(); got != nil {
		t.Fatalf("expected no craft, got %#v", got)
	}
}

func TestSoldierCraft_UsesCurrentCraftWhenPresent(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 1)
	game.craftsFile.Crafts[0].Type = geoscape.Avenger
	game.soldierFile.Soldiers[0].Craft = 1 // 1-based index into LOC.DAT/Craft references
	game.soldierFile.Soldiers[0].CraftBefore = -1

	s := &Soldier{game: game, offset: 0}
	if got := s.Craft(); got == nil {
		t.Fatal("expected craft, got nil")
	}
}

func TestSavegameCraft_OutOfBoundsReturnsNil(t *testing.T) {
	game := &Savegame{}
	if got := game.Craft(-2); got != nil {
		t.Fatalf("expected nil for negative offset, got %#v", got)
	}
	if got := game.Craft(len(game.craftsFile.Crafts)); got != nil {
		t.Fatalf("expected nil for offset beyond bounds, got %#v", got)
	}
}

func TestSoldierStats_AboveSignedByteRange(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 1)
	game.soldierFile.Soldiers[0].InitialTimeUnits = 200
	game.soldierFile.Soldiers[0].TimeUnitImprovement = 100

	s := &Soldier{game: game, offset: 0}
	if got := s.InitialTimeUnits(); got != 200 {
		t.Fatalf("InitialTimeUnits() = %d, want 200", got)
	}
	// The total must not wrap around at 255.
	if got := s.TimeUnits(); got != 300 {
		t.Fatalf("TimeUnits() = %d, want 300", got)
	}
}

func TestSoldierSetInitialStat_ClampsToByteRange(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 1)
	s := &Soldier{game: game, offset: 0}

	for _, tt := range []struct{ in, want int }{
		{135, 135},
		{255, 255},
		{300, 255},
		{-5, 0},
	} {
		s.SetInitialFiringAccuracy(tt.in)
		if got := s.InitialFiringAccuracy(); got != tt.want {
			t.Errorf("SetInitialFiringAccuracy(%d): got %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestEquipPersonalArmour(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 4)
	for i, armor := range []geoscape.Armor{geoscape.NoArmor, geoscape.PowerSuit, geoscape.NoArmor, geoscape.NoArmor} {
		game.soldierFile.Soldiers[i].Name = internal.NewNullString(fmt.Sprintf("Soldier %d", i), 25)
		game.soldierFile.Soldiers[i].Rank = geoscape.Squaddie
		game.soldierFile.Soldiers[i].Armor = armor
	}
	// soldier 2 is dead, soldier 3 is an unused entry
	game.soldierFile.Soldiers[2].Rank = geoscape.DeadOrUnused
	game.soldierFile.Soldiers[3].Name = internal.NewNullString("", 25)

	if got := game.EquipPersonalArmour(); got != 1 {
		t.Errorf("EquipPersonalArmour() = %d, want 1", got)
	}
	want := []geoscape.Armor{geoscape.PersonalArmour, geoscape.PowerSuit, geoscape.NoArmor, geoscape.NoArmor}
	for i, w := range want {
		if got := game.soldierFile.Soldiers[i].Armor; got != w {
			t.Errorf("soldier %d: armor = %v, want %v", i, got, w)
		}
	}
}

func TestSoldierResurrect(t *testing.T) {
	game := &Savegame{}
	game.soldierFile.Soldiers = make([]geoscape.SoldierData, 2)
	for i := range game.soldierFile.Soldiers {
		game.soldierFile.Soldiers[i].Name = internal.NewNullString(fmt.Sprintf("Soldier %d", i), 25)
		game.soldierFile.Soldiers[i].Rank = geoscape.Captain
	}
	game.soldierFile.Soldiers[0].Rank = geoscape.DeadOrUnused
	game.soldierFile.Soldiers[0].RecoveryDays = 12

	dead := &Soldier{game: game, offset: 0}
	alive := &Soldier{game: game, offset: 1}

	if !dead.Resurrect() {
		t.Fatal("Resurrect() of a dead soldier = false, want true")
	}
	if dead.IsDead() || dead.Rank() != Squaddie || dead.RecoveryDays() != 0 {
		t.Errorf("resurrected soldier: dead=%v rank=%v recovery=%d, want alive Squaddie without recovery",
			dead.IsDead(), dead.Rank(), dead.RecoveryDays())
	}
	if alive.Resurrect() {
		t.Error("Resurrect() of a living soldier = true, want false")
	}
	if alive.Rank() != Captain {
		t.Errorf("living soldier: rank = %v, want Captain", alive.Rank())
	}
}
