package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/go-restruct/restruct"
	"github.com/redtoad/xcom-editor/lib/resources"
)

func main() {

	buf, err := os.ReadFile("./GEODATA/WORLD.DAT")
	if err != nil {
		log.Fatalf("could not open WORLD.DAT: %w", err)
	}

	var world resources.WORLD_DAT
	if err = restruct.Unpack(buf, binary.LittleEndian, &world); err != nil {
		log.Fatalf("could not parse WORLD.DAT: %w", err)
	}

	fmt.Printf(`<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<svg width="768px" height="447px" viewBox="0 -720 2879 720" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
<g>`)
	for _, poly := range world.Polygons {
		points := fmt.Sprintf("%d %d %d %d %d %d %d %d", poly.X0, poly.Y0, poly.X1, poly.Y1, poly.X2, poly.Y2, poly.X3, poly.Y3)
		if poly.Type() == resources.Triangle {
			points = fmt.Sprintf("%d %d %d %d %d %d", poly.X0, poly.Y0, poly.X1, poly.Y1, poly.X2, poly.Y2)
		}
		fmt.Printf(`  <polygon fill="%s" points="%s"></polygon>`, resources.TerrainHexColorsXCom[poly.Terrain], points)

	}
	fmt.Printf(`</g></svg>
`)

}
