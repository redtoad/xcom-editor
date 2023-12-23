package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/redtoad/xcom-editor/internal/geoscape"
	"github.com/redtoad/xcom-editor/savegame"
)

func main() {

	rootPath := flag.String("path", ".", "save game path")
	flag.Parse()

	sg, err := savegame.Load(*rootPath)
	if err != nil {
		log.Fatalf("could not load file: %v", err)
	}
	fmt.Printf("[%v] [%s]\n", sg.MetaData.Time(), sg.MetaData.Name)
	//fmt.Printf("%v soldiers\n", sg.SoldierData.Soldiers)
	//fmt.Printf("%v bases\n", sg.BasesData.Bases)

	/*
				pathSaveinfoFile := path.Join(*rootPath, "SAVEINFO.DAT")
				if _, err := os.Stat(pathSaveinfoFile); os.IsNotExist(err) {
					fmt.Println(pathSaveinfoFile)
					log.Fatalf("could not open file: %v", err)
				}
				var info savegame.SAVEINFO_DAT
				data, err := os.ReadFile(pathSaveinfoFile)
				if err != nil {
					log.Fatalf("could not load file: %v", err)
				}
				restruct.Unpack(data, binary.LittleEndian, &info)
				fmt.Printf("[%v] [%s]\n", info.Time(), info.Name)

				pathSoldiersFile := path.Join(*rootPath, "SOLDIER.DAT")
			if _, err := os.Stat(pathSoldiersFile); os.IsNotExist(err) {
				log.Fatalf("could not open file: %v", err)
			}

			fmt.Printf("Loading %s...\n", pathSoldiersFile)
			var soldiers savegame.SOLDIER_DAT
			if err := savegame.LoadFile(pathSoldiersFile, &soldiers); err != nil {
				log.Fatalf("could not load file: %v", err)
			}

			if err := savegame.SaveFile(pathSoldiersFile+".bak", &soldiers); err != nil {
				log.Fatalf("could not create backup: %v\n", err)
			}

		for no := 0; no < len(sg.SoldierData.Soldiers); no++ {
			soldier := &sg.SoldierData.Soldiers[no]
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
	*/

	sg.Heal()
	sg.SpeedupDelivery()
	sg.CompleteConstructions()

	sg.BasesData.Bases[0].Inventory[geoscape.InventoryPersonalArmour] = 20
	sg.BasesData.Bases[0].Inventory[geoscape.InventoryFlyingSuit] = 5
	sg.BasesData.Bases[0].Inventory[geoscape.InventoryLaserPistol] = 20
	sg.BasesData.Bases[0].Inventory[geoscape.InventoryLaserRifle] = 20
	sg.BasesData.Bases[0].Inventory[geoscape.InventoryPlasmaBeam] = 4

	sg.Save()
	log.Fatal("")

	/*
		fmt.Printf("Storing %s...\n", pathSoldiersFile)
		if err := savegame.SaveFile(pathSoldiersFile, &soldiers); err != nil {
			log.Fatalf("could not save file: %v\n", err)
		}
	*/

	pathBasesFile := *rootPath + string(os.PathSeparator) + "BASE.DAT"
	if _, err := os.Stat(pathBasesFile); os.IsNotExist(err) {
		log.Fatalf("could not open file: %v", err)
	}

	fmt.Printf("Loading %s...\n", pathBasesFile)
	var bases geoscape.BASE_DAT
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
			if base.Grid[i] != geoscape.Empty && base.DaysToCompletion[i] > 0 {
				base.DaysToCompletion[i] = 0
			}
		}

		// increase Elirium-115
		Elirium115 := 60
		base.Inventory[Elirium115] = 0xfffe

		for i := 0; i < len(base.DaysToCompletion); i++ {
			base.DaysToCompletion[i] = 0
		}

		if err := savegame.SaveFile(pathBasesFile, &bases); err != nil {
			log.Fatalf("could not save file: %v\n", err)
		}

	}

	fmt.Printf("Storing %s...\n", pathBasesFile)
	if err := savegame.SaveFile(pathBasesFile, &bases); err != nil {
		log.Fatalf("could not save file: %v\n", err)
	}

}
