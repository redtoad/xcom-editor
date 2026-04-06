package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestLiglobFile_Unpack(t *testing.T) {
	buf := MustLoadFromBase64(testFile_LIGLOB_DAT)
	var lf geoscape.LiglobFile
	err := restruct.Unpack(buf, binary.LittleEndian, &lf)
	assert.NoError(t, err)

	assert.Equal(t, int32(2_107_656_887), lf.CurrentBalance)

	assert.Equal(t, 12, len(lf.Expenditure))
	assert.Equal(t, int32(5_903_230), lf.Expenditure[0])
	assert.Equal(t, int32(786_000), lf.Expenditure[7])
	assert.Equal(t, int32(0), lf.Expenditure[8])

	assert.Equal(t, 12, len(lf.Maintenance))
	assert.Equal(t, int32(3_662_000), lf.Maintenance[0])

	assert.Equal(t, 12, len(lf.Balance))
	assert.Equal(t, int32(-1), lf.Balance[0])
	assert.Equal(t, int32(2_108_442_887), lf.Balance[7])
	assert.Equal(t, int32(0), lf.Balance[8])
}

// base64 -i GAME_1/LIGLOB.DAT -b 120 | pbcopy
const testFile_LIGLOB_DAT = `
t0qgfX4TWgD0o/0AiHWnAPBeVgC8DIwAqOWZAGiNfQBQ/gsAAAAAAAAAAAAAAAAAoGEcALDgNwC44m4AmP6GAJBERABYsk8AeBB3ANiJdQAAAAAAAAAAAAAA
AAAAAAAAoGEcAP/////rf7d/7zHQfh/Ke34/j31+d50KfsdGvH0HSax9AAAAAAAAAAAAAAAAAAAAAA==
`

