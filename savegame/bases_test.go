package savegame_test

import (
	"testing"

	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBase_AdjustInventory(t *testing.T) {
	dir := copySavegame(t, "GAME_1")
	sg, err := savegame.Load(dir)
	require.NoError(t, err)
	base := sg.Base(0)
	before := base.Inventory()[geoscape.InventoryMediKit]

	n, err := base.AdjustInventory(geoscape.InventoryMediKit, 5)
	require.NoError(t, err)
	assert.Equal(t, before+5, n)
	n, err = base.AdjustInventory(geoscape.InventoryMediKit, -2)
	require.NoError(t, err)
	assert.Equal(t, before+3, n)

	// the change is saved
	require.NoError(t, sg.Save())
	sg, err = savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, before+3, sg.Base(0).Inventory()[geoscape.InventoryMediKit])
}

func TestBase_AdjustInventory_Limits(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)
	base := sg.Base(0)

	// quantities are stored as int16 and cannot be negative
	n, err := base.AdjustInventory(geoscape.InventoryLaserRifle, -1000)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	n, err = base.AdjustInventory(geoscape.InventoryLaserRifle, 40000)
	require.NoError(t, err)
	assert.Equal(t, savegame.MaxInventory, n)

	// unused slots and indices out of range are rejected
	for _, item := range []geoscape.Inventory{43, 62, -1, 96} {
		_, err = base.AdjustInventory(item, 1)
		assert.Error(t, err, "item %d", item)
	}
}

func TestInventoryItems(t *testing.T) {
	items := savegame.InventoryItems()
	assert.Len(t, items, 90)
	assert.Equal(t, geoscape.InventoryStingrayLauncher, items[0])
	assert.NotContains(t, items, geoscape.Inventory(43))
	assert.Contains(t, items, geoscape.InventoryHWPFusionBomb)
}
