package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"path"

	"github.com/mmcloughlin/globe"
	"github.com/redtoad/xcom-editor/internal"
	"github.com/redtoad/xcom-editor/internal/geodata"
	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
)

func main() {

	pth := flag.String("path", "./", "path to savegame dir")
	flag.Parse()

	fmt.Print("Loading WORLD.DAT...\n")
	world := geodata.WorldData{}
	if err := internal.LoadDATFile(path.Join(*pth, "..", "GEODATA", "WORLD.DAT"), &world); err != nil {
		log.Fatalf("could not open WORLD.DAT: %s", err)
	}

	locations := geoscape.LOC_DAT{}
	if err := internal.LoadDATFile(path.Join(*pth, "LOC.DAT"), &locations); err != nil {
		log.Fatalf("could not read data from LOC.DAT: %s", err)
	}

	//green := color.NRGBA{0x00, 0x64, 0x3c, 192}
	red := color.NRGBA{0xff, 0x0, 0x0, 192}
	blue := color.NRGBA{0x0, 0x0, 0xff, 192}

	g := globe.New()
	g.DrawGraticule(10.0)
	g.DrawLandBoundaries()

	for _, loc := range locations.Objects {
		coord := savegame.NewCoord(loc.X, loc.Y)
		x, y := float64(coord.Lat), float64(coord.Lon)
		switch loc.Type {
		case geoscape.XCOMBase:
			log.Printf("Base: %s", coord)
			//g.DrawDot(x, y, 0.1, globe.Color(green))
		case geoscape.XCOMShip:
			log.Printf("Ship: %s", coord)
			g.DrawDot(x, y, 0.1, globe.Color(blue))
		case geoscape.AlienShip:
			log.Printf("UFO: %s", coord)
			g.DrawDot(x, y, 0.05, globe.Color(red))
		case geoscape.CrashSite:
			log.Printf("Crash site: %s", coord)
			g.DrawDot(x, y, 0.5, globe.Color(red))
		case geoscape.AlienBase:
			log.Printf("Alien base: %s", coord)
			g.DrawDot(x, y, 0.5, globe.Color(red))
		case geoscape.LandedUFO:
			log.Printf("UFO landed: %s", coord)
		case geoscape.Waypoint:
			log.Printf("Waypoint: %s", coord)

		}
	}

	g.CenterOn(51.453349, -2.588323)
	g.SavePNG("land.png", 400)

	fmt.Printf("%v\n", locations)
	fmt.Printf("%v\n", world.Polygons[0])

	/*
		fmt.Printf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
			<svg width="768px" height="447px" viewBox="0 -720 2879 720" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
			<g>`)
		for _, poly := range world.Polygons {
			points := fmt.Sprintf("%d %d %d %d %d %d %d %d", poly.X0, poly.Y0, poly.X1, poly.Y1, poly.X2, poly.Y2, poly.X3, poly.Y3)
			if poly.Type() == resources.Triangle {
				points = fmt.Sprintf("%d %d %d %d %d %d", poly.X0, poly.Y0, poly.X1, poly.Y1, poly.X2, poly.Y2)
			}
			fmt.Printf(`<polygon fill="%s" points="%s"></polygon>`, resources.TerrainHexColorsXCom[poly.Terrain], points)
		}
		fmt.Printf(`</g></svg>`)
	*/
}
