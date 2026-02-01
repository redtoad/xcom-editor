package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
)

func FinishAllConstructions(path string) {

}

func main() {

	rootPath := flag.String("path", ".", "save game path")
	flag.Parse()

	pathBasesFile := *rootPath + string(os.PathSeparator) + "BASE.DAT"
	if _, err := os.Stat(pathBasesFile); os.IsNotExist(err) {
		log.Fatalf("could not open file: %v", err)
	}

	sg, _ := savegame.Load(*rootPath + string(os.PathSeparator))
	curr := currency.USD.Amount(sg.Financials.CurrentBalance)
	p := message.NewPrinter(language.AmericanEnglish)
	p.Printf("%v\n", curr)

	for no := 0; no < len(sg.BasesData.Bases); no++ {
		base := &sg.BasesData.Bases[no]
		fmt.Printf("%d  %s (%v)\n", no, base.Name, base.Active)
		if !base.Active {
			continue
		}
		for no, cell := range base.Grid {
			if no%6 == 0 {
				println()
			}
			fmt.Print(cell.String())
		}
		println()
		fmt.Printf("%v\n", base.Grid)
		fmt.Printf("%v\n", base.DaysToCompletion)

		// complete constructions in progress
		for i := 0; i < len(base.Grid); i++ {
			if base.Grid[i] != geoscape.Empty && base.DaysToCompletion[i] > 0 {
				base.DaysToCompletion[i] = 0
			}
		}

		// increase Elirium-115
		//Elirium115 := 60
		//base.Inventory[Elirium115] = 0x7f

		AlienAlloys := 88
		base.Inventory[AlienAlloys] = 0x7f

	}

	fmt.Printf("Storing %s...\n", sg.Path)
	if err := sg.Save(); err != nil {
		log.Fatalf("could not save game: %v\n", err)
	}

}
