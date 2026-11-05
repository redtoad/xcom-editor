package savegame

import "github.com/redtoad/xcom-editor/internal/geoscape"

// ResearchProject is a topic being researched at a base (PROJECT.DAT).
type ResearchProject struct {
	base  int
	topic int
	game  *Savegame
}

func (p *ResearchProject) data() *geoscape.BaseResearch {
	return &p.game.Research.Bases[p.base]
}

// Base returns the base where the research takes place.
func (p *ResearchProject) Base() *Base {
	return &Base{offset: p.base, game: p.game}
}

// Topic returns the index of the topic in RESEARCH.DAT.
func (p *ResearchProject) Topic() int {
	return p.topic
}

// Name returns the name of the research topic.
func (p *ResearchProject) Name() string {
	return geoscape.ResearchTopicNames[p.topic]
}

// DaysRemaining returns the scientist days needed to complete the research.
func (p *ResearchProject) DaysRemaining() int {
	return int(p.data().DaysRemaining[p.topic])
}

// Scientists returns the number of scientists assigned to the research.
func (p *ResearchProject) Scientists() int {
	return int(p.data().Scientists[p.topic])
}

// SpeedUp reduces the remaining research time to a single scientist day, so
// the research is completed on the next day (if scientists are assigned).
func (p *ResearchProject) SpeedUp() {
	if days := &p.data().DaysRemaining[p.topic]; *days > 1 {
		*days = 1
	}
}

// ResearchProjects returns all topics currently being researched in all bases.
func (game *Savegame) ResearchProjects() []*ResearchProject {
	projects := make([]*ResearchProject, 0)
	for b, base := range game.Research.Bases {
		for topic, days := range base.DaysRemaining {
			if days > 0 {
				projects = append(projects, &ResearchProject{base: b, topic: topic, game: game})
			}
		}
	}
	return projects
}

// SpeedupResearch speeds up all ongoing research (see ResearchProject.SpeedUp)
// and returns the number of research projects.
func (game *Savegame) SpeedupResearch() int {
	projects := game.ResearchProjects()
	for _, p := range projects {
		p.SpeedUp()
	}
	return len(projects)
}

// ProductionProject is an item being manufactured at a base (BPROD.DAT).
type ProductionProject struct {
	base    int
	product int
	game    *Savegame
}

func (p *ProductionProject) data() *geoscape.BaseProduction {
	return &p.game.Production.Bases[p.base]
}

// Base returns the base where the item is manufactured.
func (p *ProductionProject) Base() *Base {
	return &Base{offset: p.base, game: p.game}
}

// Product returns the index of the item in PRODUCT.DAT.
func (p *ProductionProject) Product() int {
	return p.product
}

// Name returns the name of the manufactured item.
func (p *ProductionProject) Name() string {
	return geoscape.ProductNames[p.product]
}

// HoursPerUnit returns the engineer hours needed to build one unit.
func (p *ProductionProject) HoursPerUnit() int {
	return int(p.game.Products.Products[p.product].Hours)
}

// HoursDone returns the engineer hours of work done towards the next unit.
func (p *ProductionProject) HoursDone() int {
	return int(p.data().HoursStored[p.product])
}

// Engineers returns the number of engineers allocated.
func (p *ProductionProject) Engineers() int {
	return int(p.data().Engineers[p.product])
}

// Total returns the total number of units to produce.
func (p *ProductionProject) Total() int {
	return int(p.data().Total[p.product])
}

// Produced returns the number of units produced so far.
func (p *ProductionProject) Produced() int {
	return int(p.data().Produced[p.product])
}

// SpeedUp sets the work done towards the next unit to one hour short of
// completion, so the unit is finished within the next hour (if engineers are
// allocated). The remaining units are built at the normal pace.
func (p *ProductionProject) SpeedUp() {
	hours := &p.data().HoursStored[p.product]
	if almostDone := int32(p.HoursPerUnit()) - 1; *hours < almostDone {
		*hours = almostDone
	}
}

// ProductionProjects returns all items currently being manufactured in all
// bases.
func (game *Savegame) ProductionProjects() []*ProductionProject {
	projects := make([]*ProductionProject, 0)
	for b, base := range game.Production.Bases {
		for product := 0; product < geoscape.NumProducts; product++ {
			// completed projects are reset to 0 in all fields
			if base.HoursStored[product] != 0 || base.Engineers[product] != 0 || base.Total[product] != 0 {
				projects = append(projects, &ProductionProject{base: b, product: product, game: game})
			}
		}
	}
	return projects
}

// SpeedupProduction speeds up all ongoing manufacturing (see
// ProductionProject.SpeedUp) and returns the number of production projects.
func (game *Savegame) SpeedupProduction() int {
	projects := game.ProductionProjects()
	for _, p := range projects {
		p.SpeedUp()
	}
	return len(projects)
}
