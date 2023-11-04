package main

import (
	"image"
	"image/png"
	"os"

	"github.com/redtoad/xcom-editor/lib/resources"
)

func main() {
	if _, err := resources.LoadFont("../../GAME/GEODATA/SMALLSET.DAT", 8, 9); err != nil {
		panic(err)
	}
	if _, err := resources.LoadFont("../../GAME/GEODATA/BIGLETS.DAT", 16, 16); err != nil {
		panic(err)
	}

	xcomFont, err := resources.LoadFont("../../GAME/GEODATA/BIGLETS.DAT", 16, 16)
	//xcomFont, err := resources.LoadFont("../../GAME/GEODATA/SMALLSET.DAT", 8, 9)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("xcom-font.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, xcomFont.Mask)
	if err != nil {
		panic(err)
	}

	//text := "The fox jumps over the lazy alien! 11$%/&?=)(/&%$§\"!^°)"
	text := "A"

	var img image.Image
	//palette := color.Palette{color.Transparent, color.Gray16{0xfffe}, color.Gray16{0xffcc}, color.Gray16{0xcccc}, color.Gray16{0x9999}, color.Gray16{0x3333}}
	img, err = xcomFont.Draw(text, resources.YellowFontPalette)
	if err != nil {
		panic(err)
	}

	f, err = os.Create("text.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

}
