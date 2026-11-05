package savegame

import (
	"testing"

	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
)

func TestRepairAndRefuelAllCrafts(t *testing.T) {
	const noWeapon = geoscape.WeaponType(-1) // 255 read as signed byte
	const baseLoc = 10                       // index of the base in LOC.DAT

	game := &Savegame{}
	game.BasesData.Bases = make([]geoscape.BaseData, 1)
	game.BasesData.Bases[0].Name = internal.NewNullString("Base", 16)
	game.BasesData.Bases[0].Inventory[geoscape.InventoryAvalancheMissile] = 10
	game.BasesData.Bases[0].Inventory[geoscape.InventoryCannonRounds] = 1
	game.locationFile.Objects[baseLoc] = geoscape.LocationData{Type: geoscape.XCOMBase, TableReference: 0}

	crafts := []geoscape.CraftData{
		// damaged interceptor waiting for repairs
		{Type: geoscape.Interceptor, LeftWeapon: geoscape.Avalanche, LeftAmmo: 1,
			RightWeapon: geoscape.Cannon, RightAmmo: 50,
			Damage: 40, FuelAmountRemaining: 200, Status: geoscape.Repairs, BaseReference: baseLoc},
		// unarmed skyranger refuelling
		{Type: geoscape.Skyranger, LeftWeapon: noWeapon, RightWeapon: noWeapon,
			FuelAmountRemaining: 100, Status: geoscape.Refueling, BaseReference: baseLoc},
		// airborne avenger with more fuel than its capacity
		{Type: geoscape.Avenger, LeftWeapon: geoscape.PlasmaBeam, LeftAmmo: 20,
			RightWeapon: geoscape.Stingray, RightAmmo: 0,
			Damage: 10, FuelAmountRemaining: 1000, Status: geoscape.Out, Speed: 500, BaseReference: baseLoc},
		// UFOs are not touched
		{Type: geoscape.Battleship, Damage: 100, FuelAmountRemaining: 1, Speed: 100, AlienRace: 3},
		// damaged (zeroed) entries are not touched
		{Type: geoscape.Interceptor},
	}
	for i, data := range crafts {
		game.craftsFile.Crafts[i] = data
		objType := geoscape.XCOMShip
		if data.Type >= geoscape.SmallScout {
			objType = geoscape.AlienShip
		}
		game.locationFile.Objects[i] = geoscape.LocationData{Type: objType, TableReference: uint(i)}
	}

	if got := game.RepairAndRefuelAllCrafts(); got != 3 {
		t.Errorf("RepairAndRefuelAllCrafts() = %d, want 3", got)
	}

	tests := []struct {
		name        string
		damage      int
		fuel        int
		left, right int
		status      geoscape.CraftStatus
	}{
		{"interceptor", 0, 1000, 3, 200, geoscape.Ready},
		{"skyranger", 0, 1500, 0, 0, geoscape.Ready},
		{"avenger", 0, 1000, 100, 6, geoscape.Out},
		{"ufo", 100, 1, 0, 0, geoscape.Ready},
		{"damaged entry", 0, 0, 0, 0, geoscape.Ready},
	}
	for i, tt := range tests {
		c := game.craftsFile.Crafts[i]
		if c.Damage != tt.damage || c.FuelAmountRemaining != tt.fuel ||
			c.LeftAmmo != tt.left || c.RightAmmo != tt.right || c.Status != tt.status {
			t.Errorf("%s: damage=%d fuel=%d ammo=%d/%d status=%v, want damage=%d fuel=%d ammo=%d/%d status=%v",
				tt.name, c.Damage, c.FuelAmountRemaining, c.LeftAmmo, c.RightAmmo, c.Status,
				tt.damage, tt.fuel, tt.left, tt.right, tt.status)
		}
	}

	// the stores hold enough ammo to rearm all craft once more
	inventory := game.BasesData.Bases[0].Inventory
	for item, want := range map[geoscape.Inventory]int{
		geoscape.InventoryAvalancheMissile: 10, // enough already
		geoscape.InventoryCannonRounds:     4,  // 200 rounds = 4 x 50
		geoscape.InventoryStingrayMissile:  6,
		geoscape.InventoryFusionBalls:      0,
	} {
		if got := inventory[item]; got != want {
			t.Errorf("stores: %v = %d, want %d", item, got, want)
		}
	}
}

func TestCraftMaxFuel(t *testing.T) {
	game := &Savegame{}
	game.craftsFile.Crafts[0].Type = geoscape.Lightning
	game.craftsFile.Crafts[1].Type = geoscape.SmallScout
	if got := (&Craft{game: game, offset: 0}).MaxFuel(); got != 30 {
		t.Errorf("Lightning MaxFuel() = %d, want 30", got)
	}
	if got := (&Craft{game: game, offset: 1}).MaxFuel(); got != 0 {
		t.Errorf("UFO MaxFuel() = %d, want 0", got)
	}
}
