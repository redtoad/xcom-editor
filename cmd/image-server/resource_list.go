package main

import "fmt"

type ImageEntry struct {
	PaletteNr int
	Width     int
	Height    int
	TabFile   string
	TabOffset int
}

var images = map[string]ImageEntry{

	// ...
	"GEOGRAPH/BACK01.SCR":   {0, 320, 200, "", 0},
	"GEOGRAPH/BACK02.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK03.SCR":   {3, 320, 200, "", 0},
	"GEOGRAPH/BACK04.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK05.SCR":   {3, 320, 200, "", 0},
	"GEOGRAPH/BACK06.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK07.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK08.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK09.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK10.SCR":   {3, 320, 200, "", 0},
	"GEOGRAPH/BACK11.SCR":   {0, 320, 200, "", 0},
	"GEOGRAPH/BACK12.SCR":   {4, 320, 200, "", 0},
	"GEOGRAPH/BACK13.SCR":   {0, 320, 200, "", 0},
	"GEOGRAPH/BACK14.SCR":   {3, 320, 200, "", 0},
	"GEOGRAPH/BACK15.SCR":   {3, 320, 200, "", 0},
	"GEOGRAPH/BACK16.SCR":   {0, 320, 200, "", 0},
	"GEOGRAPH/BACK17.SCR":   {0, 320, 200, "", 0},
	"GEOGRAPH/BASEBITS.PCK": {1, 32, 40, "GEOGRAPH/BASEBITS.TAB", 2},
	// ...
	"GEOGRAPH/GEOBORD.SCR": {0, 320, 200, "", 0},
	"GEOGRAPH/INTICON.PCK": {0, 32, 40, "GEOGRAPH/INTICON.TAB", 2},
	// ...
	"TERRAIN/XBASE1.PCK": {4, 32, 40, "TERRAIN/XBASE1.TAB", 2},
	"TERRAIN/XBASE2.PCK": {4, 32, 40, "TERRAIN/XBASE2.TAB", 2},
	// ...
	"UFOGRAPH/X1.PCK": {1, 128, 40, "UFOGRAPH/X1.TAB", 2},
	// ...
	"UNITS/BIGOBS.PCK": {4, 32, 48, "UNITS/BIGOBS.TAB", 2},
	"UNITS/XCOM_0.PCK": {4, 32, 40, "UNITS/XCOM_0.TAB", 2},
	"UNITS/XCOM_1.PCK": {4, 32, 40, "UNITS/XCOM_1.TAB", 2},
	"UNITS/XCOM_2.PCK": {4, 32, 40, "UNITS/XCOM_2.TAB", 2},
	"UNITS/X_REAP.PCK": {1, 32, 40, "UNITS/X_REAP.TAB", 2},
	"UNITS/X_ROB.PCK":  {1, 32, 40, "UNITS/X_ROB.TAB", 2},
	// ...
	"UNITS/ZOMBIE.PCK": {1, 32, 40, "UNITS/ZOMBIE.TAB", 2},

	// INTERWIN.DAT - 10 images 160px wide no compression
	// LANG1.DAT - GeoScape control panel overlay in German. Direct palette indexes, no compression, 64x154.
	// LANG2.DAT - GeoScape control panel overlay in French. Direct palette indexes, no compression, 64x154.

}

func init() {

	// add all 42 UFOPedia backgrounds
	for i := 1; i <= 42; i++ {
		images[fmt.Sprintf("GEOGRAPH/UP_%03d.SPK", i)] = ImageEntry{3, 320, 200, "", 0}
	}

	// add all 56 BIGOBs
	for i := 0; i <= 56; i++ {
		images[fmt.Sprintf("UFOGRAPH/BIGOB_%02d.PCK", i)] = ImageEntry{1, 32, 40, "", 4}
	}

	// add soldier images
	for armour := 0; armour <= 1; armour++ {
		for _, sex := range []string{"F", "M"} {
			for race := 0; race <= 3; race++ {
				images[fmt.Sprintf("UFOGRAPH/MAN_%d%s%d.SPK", armour, sex, race)] = ImageEntry{4, 32, 40, "", 4}
			}
		}
	}
	// soldiers in armoured suits
	images["UFOGRAPH/MAN_2.SPK"] = ImageEntry{4, 32, 40, "", 4}
	images["UFOGRAPH/MAN_3.SPK"] = ImageEntry{4, 32, 40, "", 4}
}
