package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/redtoad/xcom-editor/lib/savegame"
)

func main() {

	rootPath := flag.String("path", ".", "save game path")
	flag.Parse()

	pathSoldiersFile := path.Join(*rootPath, "SOLDIER.DAT")
	if _, err := os.Stat(pathSoldiersFile); os.IsNotExist(err) {
		log.Fatalf("could not open file: %v", err)
	}

	fmt.Printf("Loading %s...\n", pathSoldiersFile)
	var soldiers savegame.FileSoldier
	if err := savegame.LoadFile(pathSoldiersFile, &soldiers); err != nil {
		log.Fatalf("could not load file: %v", err)
	}

	if err := savegame.SaveFile(pathSoldiersFile+".bak", &soldiers); err != nil {
		log.Fatalf("could not create backup: %v\n", err)
	}

	for no := 0; no < len(soldiers.Soldiers); no++ {
		soldier := &soldiers.Soldiers[no]
		// resurrect solders
		if soldier.Rank == savegame.DeadOrUnused && strings.TrimSpace(soldier.Name) != "" {
			fmt.Printf("Resurrect %s from the dead\n", soldier.Name)
			soldier.Rank = savegame.Squaddie
		}
		if soldier.Rank != savegame.DeadOrUnused {
			fmt.Printf("%d  %v %s (%v)\n", no, soldier.Rank, soldier.Name, soldier.Armor)
			soldier.Armor = savegame.PersonalArmor
			soldier.InitialFiringAccuracy += 10
			soldier.InitialTimeUnits += 10
			soldier.InitialReactions += 10
			soldier.InitialBravery = 0
			soldier.InitialEnergy += 10
			soldier.RecoveryDays = 0
		}
	}

	fmt.Printf("Storing %s...\n", pathSoldiersFile)
	if err := savegame.SaveFile(pathSoldiersFile, &soldiers); err != nil {
		log.Fatalf("could not save file: %v\n", err)
	}

	pathBasesFile := *rootPath + string(os.PathSeparator) + "BASE.DAT"
	if _, err := os.Stat(pathBasesFile); os.IsNotExist(err) {
		log.Fatalf("could not open file: %v", err)
	}

	fmt.Printf("Loading %s...\n", pathBasesFile)
	var bases savegame.FileBase
	if err := savegame.LoadFile(pathBasesFile, &bases); err != nil {
		log.Fatalf("could not load file: %v", err)
	}

	if err := savegame.SaveFile(pathBasesFile+".bak", &bases); err != nil {
		log.Fatalf("could not create backup: %v\n", err)
	}

	for no := 0; no < len(bases.Bases); no++ {
		base := &bases.Bases[no]
		fmt.Printf("%d  %s (%v)\n", no, base.Name, base.Active)
		if !base.Active {
			continue
		}
		for no, cell := range base.Grid {
			if no%6 == 0 {
				println()
			}
			fmt.Print(cell.Tile())
		}
		println()
		fmt.Printf("%v\n", base.Grid)
		fmt.Printf("%v\n", base.DaysToCompletion)

		// complete constructions in progress
		for i := 0; i < len(base.Grid); i++ {
			if base.Grid[i] != savegame.Empty && base.DaysToCompletion[i] > 0 {
				base.DaysToCompletion[i] = 0
			}
		}

		// increase Elirium-115
		Elirium115 := 60
		base.Inventory[Elirium115] = 0xfffe
	}

	fmt.Printf("Storing %s...\n", pathBasesFile)
	if err := savegame.SaveFile(pathBasesFile, &bases); err != nil {
		log.Fatalf("could not save file: %v\n", err)
	}

}
