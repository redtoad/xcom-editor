package savegame_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResearchProjects(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)

	projects := sg.ResearchProjects()
	require.Len(t, projects, 1)
	p := projects[0]
	assert.Equal(t, 0, p.Base().Index())
	assert.Equal(t, 15, p.Topic())
	assert.Equal(t, "Elerium-115", p.Name())
	assert.Equal(t, 107, p.DaysRemaining())
	assert.Equal(t, 57, p.Scientists())
}

func TestResearchProject_SpeedUp(t *testing.T) {
	dir := copySavegame(t, "GAME_1")
	sg, err := savegame.Load(dir)
	require.NoError(t, err)

	p := sg.ResearchProjects()[0]
	p.SpeedUp()
	assert.Equal(t, 1, p.DaysRemaining())
	assert.Equal(t, 57, p.Scientists())

	// the change is saved
	require.NoError(t, sg.Save())
	sg, err = savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, 1, sg.ResearchProjects()[0].DaysRemaining())
}

func TestSpeedupResearch(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)
	// research two more topics in another base
	sg.Research.Bases[1].DaysRemaining[3] = 1
	sg.Research.Bases[1].DaysRemaining[4] = 20

	assert.Equal(t, 3, sg.SpeedupResearch())
	for _, p := range sg.ResearchProjects() {
		assert.Equal(t, 1, p.DaysRemaining(), p.Name())
	}
	// topics not being researched are not started
	assert.Equal(t, uint16(0), sg.Research.Bases[0].DaysRemaining[0])
}

func TestResearchTopicNames(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_3"))
	require.NoError(t, err)
	for i, topic := range sg.ResearchTopics.Topics {
		name := geoscape.ResearchTopicNames[i]
		require.NotEmpty(t, name, "topic %d", i)
		// topics for items in stores are named after the item
		if item := topic.Topic[0]; item != 0 {
			assert.Equal(t, normalize(geoscape.Inventory(item).String()), normalize(name), "topic %d", i)
		}
		// topics for live aliens are named after race and rank
		if race := geoscape.AlienRace(topic.Topic[1]); race != geoscape.AlienNone {
			assert.Equal(t, race.String()+" "+geoscape.AlienRank(topic.Topic[2]).String(), name, "topic %d", i)
		}
	}
}

func TestProductionProjects(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)

	projects := sg.ProductionProjects()
	require.Len(t, projects, 1)
	p := projects[0]
	assert.Equal(t, 0, p.Base().Index())
	assert.Equal(t, 12, p.Product())
	assert.Equal(t, "Medi-Kit", p.Name())
	assert.Equal(t, 420, p.HoursPerUnit())
	assert.Equal(t, 380, p.HoursDone())
	assert.Equal(t, 10, p.Engineers())
	assert.Equal(t, 10, p.Total())
	assert.Equal(t, 3, p.Produced())
}

func TestProductionProject_SpeedUp(t *testing.T) {
	dir := copySavegame(t, "GAME_1")
	sg, err := savegame.Load(dir)
	require.NoError(t, err)

	p := sg.ProductionProjects()[0]
	p.SpeedUp()
	// the next unit is finished within the next hour
	assert.Equal(t, 419, p.HoursDone())
	assert.Equal(t, 3, p.Produced())
	assert.Equal(t, 10, p.Total())

	require.NoError(t, sg.Save())
	sg, err = savegame.Load(dir)
	require.NoError(t, err)
	assert.Equal(t, 419, sg.ProductionProjects()[0].HoursDone())
}

func TestSpeedupProduction(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_1"))
	require.NoError(t, err)
	// a second project which is almost done already is not set back
	b := &sg.Production.Bases[1]
	b.HoursStored[8] = 299
	b.Engineers[8] = 5
	b.Total[8] = 2

	assert.Equal(t, 2, sg.SpeedupProduction())
	projects := sg.ProductionProjects()
	require.Len(t, projects, 2)
	assert.Equal(t, 419, projects[0].HoursDone())
	assert.Equal(t, "Laser Pistol", projects[1].Name())
	assert.Equal(t, 299, projects[1].HoursDone())
}

func TestProductNames(t *testing.T) {
	sg, err := savegame.Load(copySavegame(t, "GAME_3"))
	require.NoError(t, err)
	for i, product := range sg.Products.Products {
		name := geoscape.ProductNames[i]
		if product.Item >= 200 {
			// crafts
			assert.Equal(t, geoscape.CraftType(product.Item-200).String(), name, "product %d", i)
		} else {
			assert.Equal(t, normalize(geoscape.Inventory(product.Item).String()), normalize(name), "product %d", i)
		}
	}
}

// normalize makes item names comparable by ignoring case, punctuation and a
// plural "s" (e.g. "Elerium-115" and "Elerium115", "Fusion Ball" and
// "FusionBalls").
func normalize(name string) string {
	name = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, name)
	return strings.TrimSuffix(name, "s")
}
