package savegame

import (
	"testing"

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
