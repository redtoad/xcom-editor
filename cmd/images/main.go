package main

import (
	"fmt"
	"github.com/redtoad/xcom-editor/resources"
	"image"
	"image/png"
	"os"
)

func main() {
	palettes, err := resources.LoadPalettes("/Users/srahlf/Desktop/privat/X-COM UFO Defense/GEODATA/PALETTES.DAT")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", palettes)

	width := 256
	height := 5

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, palettes[y].Colors[x])
		}
	}

	// Encode as PNG.
	f, _ := os.Create("palettes.png")
	png.Encode(f, img)

	//offsets, err := resources.LoadTAB(
	//	//"/Users/srahlf/Desktop/privat/X-COM UFO Defense/UFOGRAPH/X1.TAB",
	//	"/Users/srahlf/Desktop/privat/X-COM UFO Defense/UNITS/XCOM_0.TAB",
	//	2)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("Offsets: %v", offsets)
	//
	//for no, offset := range offsets {
	//	sprite, err := resources.LoadPCK(
	//		//"/Users/srahlf/Desktop/privat/X-COM UFO Defense/UFOGRAPH/X1.PCK",
	//		"/Users/srahlf/Desktop/privat/X-COM UFO Defense/UNITS/XCOM_0.PCK",
	//		32, offset)
	//	if err != nil {
	//		panic(err)
	//	}
	//	img := sprite.Image(palettes[4])
	//	f, _ := os.Create(fmt.Sprintf("xcom_%d.png", no))
	//	err = png.Encode(f, img)
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	sprite, err := resources.LoadSCR("/Users/srahlf/Desktop/privat/X-COM UFO Defense/GEOGRAPH/GEOBORD.SCR", 320)
	if err != nil {
		panic(err)
	}

	img2 := sprite.Image(palettes[4])
	f, _ = os.Create("back06.png")
	err = png.Encode(f, img2)

	sprite2, err := resources.LoadSPK("/Users/srahlf/Desktop/privat/X-COM UFO Defense/UFOGRAPH/MAN_0F3.SPK")
	if err != nil {
		panic(err)
	}

	f, err = os.Create("MAN_03F.png")
	if err != nil {
		panic(err)
	}
	_ = png.Encode(f, sprite2.Image(palettes[4]))

}
