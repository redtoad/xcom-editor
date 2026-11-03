package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestAlienFile_Unpack(t *testing.T) {
	buf := MustLoadTestFile(t, "PARTIAL", "ALIEN.DAT")
	var f geoscape.AlienFile
	err := restruct.Unpack(buf, binary.LittleEndian, &f)
	assert.NoError(t, err)
	assert.Equal(t, 12, len(f.Months))

	// Month 0 (current month)
	m0 := f.Months[0]
	assert.Equal(t, int32(2), m0.USA)
	assert.Equal(t, int32(25), m0.Russia)
	assert.Equal(t, int32(15), m0.UK)
	assert.Equal(t, int32(5), m0.France)
	assert.Equal(t, int32(38), m0.Brazil)
	assert.Equal(t, int32(3), m0.NorthAmerica)
	assert.Equal(t, int32(1), m0.Arctic)
	assert.Equal(t, int32(2), m0.Antarctica)
	assert.Equal(t, int32(195), m0.SouthAmerica)

	// Month 11 (oldest, should be zero for a new game)
	m11 := f.Months[11]
	assert.Equal(t, int32(0), m11.USA)
	assert.Equal(t, int32(0), m11.Europe)
}
