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
