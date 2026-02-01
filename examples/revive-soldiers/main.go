package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/redtoad/xcom-editor/savegame"
)

var Reset = "\033[0m"
var Red = "\033[31m"

func main() {

	rootPath := flag.String("path", ".", "save game path")
	flag.Parse()

	sg, err := savegame.Load(*rootPath)
	if err != nil {
		log.Fatalf("could not load savegame: %v", err)
	}
	fmt.Printf("Savegame %s: time=%v title=%s\n",
		*rootPath, sg.Time(), sg.Title())

	for _, soldier := range sg.Soldiers() {
		var craftName = ""
		if craft := soldier.Craft(); craft == nil {
			craftName = "no craft"
		} else {
			craftName = craft.Name()
		}
		line := fmt.Sprintf("%s (%s) @ %v / %s",
			soldier.Name(), soldier.Rank(),
			craftName, soldier.Base().Name())
		if soldier.IsDead() {
			soldier.Heal()
			fmt.Println(Red + line + " revived as Squaddie" + Reset)
		} else if soldier.IsWounded() {
			soldier.Heal()
			fmt.Println(Red + line + " healed" + Reset)
		} else {
			fmt.Println(line)
		}
	}

	if err := sg.Save(); err != nil {
		panic(err)
	}

}
