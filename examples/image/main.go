package main

import (
	"fmt"
	"image"
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

	var sprites []*image.Paletted
	for _, offset := range offsets {
		sprite, err := resources.LoadPCK("UNITS/ZOMBIE.PCK", 32, offset)
		if err != nil {
			panic(err)
		}
		sprites = append(sprites, sprite.Paletted(palettes[4]))
	}

	delay := make([]int, len(sprites))
	disposal := make([]byte, len(sprites))
	for i := 0; i < len(delay); i++ {
		delay[i] = 25
		disposal[i] = gif.DisposalBackground
	}
	animated := &gif.GIF{
		Image:           sprites,
		Delay:           delay,
		LoopCount:       0,
		Disposal:        disposal,
		BackgroundIndex: 0,
		Config:          image.Config{Width: 32, Height: 40},
	}

	f, _ := os.Create("test.gif")
	err = gif.EncodeAll(f, animated)
	if err != nil {
		panic(err)
	}

}
