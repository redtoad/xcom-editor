package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CRAFT_DAT_Unpack(t *testing.T) {
	buf := MustLoadTestFile(t, "PARTIAL", "CRAFT.DAT")
	var obj geoscape.CraftFile
	remaining, err := obj.Unpack(buf, binary.LittleEndian)
	require.NoError(t, err)
	assert.Empty(t, remaining, "expected no remaining bytes")
	assert.EqualValues(t, 50, len(obj.Crafts), "should have all 50 slots")
	// First 8 entries have valid craft types, rest are EntryNotUsed
	for i := 0; i < 8; i++ {
		assert.NotEqual(t, geoscape.EntryNotUsed, obj.Crafts[i].Type, "craft %d should be active", i)
	}
	for i := 8; i < 50; i++ {
		assert.Equal(t, geoscape.CraftType(geoscape.EntryNotUsed), obj.Crafts[i].Type, "craft %d should be unused", i)
	}
}
