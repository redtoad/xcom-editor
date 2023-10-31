package main

import (
	"image/color"
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

	font, err := resources.LoadFont("../../GAME/GEODATA/BIGLETS.DAT", 16, 16)
	//font, err := resources.LoadFont("../../GAME/GEODATA/SMALLSET.DAT", 8, 9)
	if err != nil {
		panic(err)
	}

	palette := color.Palette{color.Transparent, color.Gray16{0xfffe}, color.Gray16{0xffcc}, color.Gray16{0xcccc}, color.Gray16{0x9999}, color.Gray16{0x3333}}
	img, err := resources.Text("Hello world! 11$%/&", font, palette)
	if err != nil {
		panic(err)
	}

	f, _ := os.Create("text.png")
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

}
