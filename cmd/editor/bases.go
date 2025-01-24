package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/redtoad/xcom-editor/internal/geoscape"
)

func tileTemplate(tile geoscape.Facility) string {
	switch tile {
	case geoscape.AccessLift:
		return "access-lift"
	case geoscape.LivingQuarters:
		return "quarters"
	case geoscape.Laboratory:
		return "laboratory"
	case geoscape.Workshop:
		return "workshop"
	case geoscape.SmallRadarSystem:
		return "radar-small"
	case geoscape.LargeRadarSystem:
		return "radar-big"
	case geoscape.MissileDefense:
		return "missile-defenses"
	case geoscape.GeneralStores:
		return "general-stores"
	case geoscape.AlienContainment:
		return "alien-containment"
	case geoscape.LaserDefense:
		return "lase-defenses"
	case geoscape.PlasmaDefense:
		return "plasma-defenses"
	case geoscape.FusionBallDefense:
		return "fusion-ball-defenses"
	case geoscape.GravShield:
		return "grav-shield"
	case geoscape.MindShield:
		return "mind-shield"
	case geoscape.PsionicLaboratory:
		return "laboratory"
	case geoscape.HyperwaveDecoder:
		return "hyperwave-decoder"
	case geoscape.HangarTopLeft:
		return "hangar-tl"
	case geoscape.HangarTopRight:
		return "hangar-tr"
	case geoscape.HangarBottomLeft:
		return "hangar-bl"
	case geoscape.HangarBottomRight:
		return "hangar-br"
	case geoscape.Empty:
		return "empty"
	}
	log.Printf("unknown tile: %v", tile)
	return "empty"
}

func Bases(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	files := []string{
		"templates/bases.html",
	}
	tpl := template.New("bases")
	tpl.Funcs(template.FuncMap{
		"tile": func(tile geoscape.Facility) (ret template.HTML, err error) {
			buf := bytes.NewBuffer([]byte{})
			err = tpl.ExecuteTemplate(buf, tileTemplate(tile), nil)
			ret = template.HTML(buf.String())
			return
		},
		"mod": func(i, j, k int) bool { return i%j == k },
	})

	_, err := tpl.ParseFS(Templates, files...)

	if err != nil {
		log.Printf("could not load templates: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	bases := currentSavegame.Bases()
	err = tpl.ExecuteTemplate(w, "bases", bases)
	if err != nil {
		log.Printf("error: %v", err)
	}
}
