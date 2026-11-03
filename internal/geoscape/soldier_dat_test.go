package geoscape_test

import (
	"encoding/binary"
	"testing"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/stretchr/testify/assert"
)

func TestUnpackSoldier(t *testing.T) {
	tests := []struct {
		hexData  string
		expected geoscape.SoldierData
	}{
		{
			`00000000 0100FFFF 00000000 00001400 47756472
			756E2055 6E676572 00000000 00000000 00000000
			00003919 32211635 391A1800 08000000 00000000
			00000000 00000100`,
			geoscape.SoldierData{
				Rank:                        geoscape.Rookie,
				Base:                        0,
				Craft:                       1,
				CraftBefore:                 -1,
				Missions:                    0,
				Kills:                       0,
				RecoveryDays:                0,
				SoldierValue:                20,
				Name:                        internal.NewNullString("Gudrun Unger", 25),
				DestinationBase:             0,
				InitialTimeUnits:            57,
				InitialHealth:               25,
				InitialEnergy:               50,
				InitialReactions:            33,
				InitialStrength:             22,
				InitialFiringAccuracy:       53,
				InitialThrowingAccuracy:     57,
				InitialMeleeAccuracy:        26,
				InitialPsionicStrength:      24,
				InitialPsionicSkill:         0,
				InitialBravery:              8,
				TimeUnitImprovement:         0,
				HealthImprovement:           0,
				EnergyImprovement:           0,
				ReactionsImprovement:        0,
				StrengthImprovement:         0,
				FiringAccuracyImprovement:   0,
				ThrowingAccuracyImprovement: 0,
				MeleeAccuracyImprovement:    0,
				BraveryImprovement:          0,
				Armor:                       geoscape.NoArmor,
				MostRecentPsiLabTraining:    0,
				InPsiLabTraining:            false,
				Promotion:                   false,
				Sex:                         geoscape.Female,
				Appearance:                  geoscape.Blonde,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.hexData, func(t *testing.T) {
			data := MustLoadFromHex(tt.hexData)
			var soldier geoscape.SoldierData
			err := restruct.Unpack(data, binary.LittleEndian, &soldier)
			assert.NoError(t, err, "could not unpack test data: %v", err)

			assert.Equal(t, tt.expected.Name.String(), soldier.Name.String())
			// Name carries null padding from binary data; normalize before struct comparison.
			soldier.Name = tt.expected.Name
			assert.Equal(t, tt.expected, soldier)

		})
	}
}

// TestSoldierStatsAreUnsigned ensures that stat bytes above 127 (e.g. from
// hacked or heavily trained soldiers) are read as unsigned values and written
// back unchanged.
func TestSoldierStatsAreUnsigned(t *testing.T) {
	data := make([]byte, 68)
	// Bytes 42-51: initial stats, 53-60: stat improvements, 63: psi lab increase.
	for _, offset := range []int{42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 53, 54, 55, 56, 57, 58, 59, 60, 63} {
		data[offset] = 0x87 // 135
	}

	var soldier geoscape.SoldierData
	assert.NoError(t, restruct.Unpack(data, binary.LittleEndian, &soldier))

	for name, got := range map[string]int{
		"InitialTimeUnits":            int(soldier.InitialTimeUnits),
		"InitialHealth":               int(soldier.InitialHealth),
		"InitialEnergy":               int(soldier.InitialEnergy),
		"InitialReactions":            int(soldier.InitialReactions),
		"InitialStrength":             int(soldier.InitialStrength),
		"InitialFiringAccuracy":       int(soldier.InitialFiringAccuracy),
		"InitialThrowingAccuracy":     int(soldier.InitialThrowingAccuracy),
		"InitialMeleeAccuracy":        int(soldier.InitialMeleeAccuracy),
		"InitialPsionicStrength":      int(soldier.InitialPsionicStrength),
		"InitialPsionicSkill":         int(soldier.InitialPsionicSkill),
		"TimeUnitImprovement":         int(soldier.TimeUnitImprovement),
		"HealthImprovement":           int(soldier.HealthImprovement),
		"EnergyImprovement":           int(soldier.EnergyImprovement),
		"ReactionsImprovement":        int(soldier.ReactionsImprovement),
		"StrengthImprovement":         int(soldier.StrengthImprovement),
		"FiringAccuracyImprovement":   int(soldier.FiringAccuracyImprovement),
		"ThrowingAccuracyImprovement": int(soldier.ThrowingAccuracyImprovement),
		"MeleeAccuracyImprovement":    int(soldier.MeleeAccuracyImprovement),
		"MostRecentPsiLabTraining":    int(soldier.MostRecentPsiLabTraining),
	} {
		assert.Equal(t, 135, got, name)
	}

	packed, err := restruct.Pack(binary.LittleEndian, &soldier)
	assert.NoError(t, err)
	assert.Equal(t, data, packed)
}
