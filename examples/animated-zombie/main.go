package main

import (
	"fmt"
	"image/gif"
	"os"

	"github.com/redtoad/xcom-editor/resources"
)

func main() {
	palettes, err := resources.LoadPalettes("GEODATA/PALETTES.DAT")
	if err != nil {
		panic(err)
	}

	offsets, err := resources.LoadTAB("UNITS/ZOMBIE.TAB", 2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Offsets: %v", offsets)

	collection, err := resources.LoadImageCollectionFromPCK("UNITS/ZOMBIE.PCK", 32, offsets)
	if err != nil {
		panic(err)
	}

	animated := collection.Animated(25, 40, palettes[4])

	f, _ := os.Create("test.gif")
	err = gif.EncodeAll(f, animated)
	if err != nil {
		panic(err)
	}

}
