package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/redtoad/xcom-editor/resources"
	"github.com/redtoad/xcom-editor/savegame"
)

const DefaultPort = 8080

var loader *resources.ResourceLoader

// Let's face it. This will never be used in any concurrent setting with mutliple users at the same
// time. So we'll not bother with an elaborate service architecture but simply use global variables
// which we pass to functions during initialisation.

var rootPath *string
var currentSavegame *savegame.Savegame

// CheckGameRoot returns an error if root is not pointing at the root directory of the game.
func CheckGameRoot(root string) error {
	exists := func(p string) error {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %w", err)
		}
		return nil
	}
	paths := []string{
		root,
		path.Join(root, "UFOEXE", "GEOSCAPE.EXE"),
		path.Join(root, "GEODATA", "PALETTES.DAT"),
		path.Join(root, "UNITS"),
	}
	for _, p := range paths {
		if err := exists(p); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	rootPath = flag.String("path", ".", "game path")
	port := *flag.Int("port", DefaultPort, fmt.Sprintf("network port (default: %d", DefaultPort))
	flag.Parse()

	if err := CheckGameRoot(*rootPath); err != nil {
		log.Fatalf("game not found: %v", err)
	}

	// TODO select and load savegame
	savegamePath := path.Join(*rootPath, "GAME_1")
	log.Printf("loading savegame %s...", savegamePath)
	currentSavegame, _ = savegame.Load(savegamePath)
	curr := currency.USD.Amount(currentSavegame.FinancialData.CurrentBalance)
	p := message.NewPrinter(language.AmericanEnglish)
	p.Printf("%v\n", curr)

	var err error
	loader, err = resources.NewResourceLoader(*rootPath)
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./cmd/editor/static"))

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.PathPrefix("/resource").HandlerFunc(ServeImage)
	r.PathPrefix("/bases").HandlerFunc(Bases)
	r.PathPrefix("/locations").HandlerFunc(locations)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%d", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Printf("Starting server...\n")
	log.Printf("Game root: %s\n", *rootPath+"/..")
	log.Printf("Try opening http://localhost:%d/bases\n", port)

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	log.Println("Press Ctrl+C to stop.")

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	_ = srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

	/*
		for no := 0; no < len(sg.Bases); no++ {
			base := &sg.Bases[no]
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
			//Elirium115 := 60
			//base.Inventory[Elirium115] = 0x7f

			AlienAlloys := 88
			base.Inventory[AlienAlloys] = 0x7f

		}
	*/

	fmt.Printf("Storing %s...\n", currentSavegame.Path)
	if err := currentSavegame.Save(); err != nil {
		log.Fatalf("could not save game: %v\n", err)
	}

}
