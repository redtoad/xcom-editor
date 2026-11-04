package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/redtoad/xcom-editor/savegame"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

//go:embed frontend/dist
var frontendFS embed.FS

// gameEntry holds a savegame and its synchronization mutex.
type gameEntry struct {
	sg *savegame.Savegame
	mu sync.RWMutex
	// stamp is the latest modification time of the savegame files when the
	// savegame was loaded or saved by the editor.
	stamp time.Time
	// stale is true if the files were changed since (e.g. by the game).
	stale bool
}

// games holds all loaded savegames indexed by slot name (e.g. "GAME_1").
var games map[string]*gameEntry

// currentRoot is the folder the savegames were loaded from, it is watched
// for new and changed savegames (see watchGames).
var currentRoot string

// stateMu guards games, currentRoot, world, fonts, baseBits, soldierImages,
// terrain and unitSheets. API requests hold a read lock for their whole
// lifetime (see lockState), loadRoot replaces all of them at once under the
// write lock.
var stateMu sync.RWMutex

// loadRoot loads all savegames (GAME_1 through GAME_10) and the optional
// game resources (globe terrain, fonts, base graphics, soldier images and
// Battlescape graphics) from root. If no
// savegame can be loaded, the current state is kept.
func loadRoot(root string) error {
	loaded := make(map[string]*gameEntry)
	for i := 1; i <= maxSlots; i++ {
		slot := fmt.Sprintf("GAME_%d", i)
		gamePath := filepath.Join(root, slot)
		if _, err := os.Stat(filepath.Join(gamePath, "SAVEINFO.DAT")); os.IsNotExist(err) {
			continue
		}
		entry, err := loadGame(gamePath)
		if err != nil {
			log.Printf("Warning: could not load %s: %v", slot, err)
			continue
		}
		loaded[slot] = entry
		log.Printf("Loaded %s: %s", slot, entry.sg.Title())
	}

	if len(loaded) == 0 {
		return fmt.Errorf("no savegames found in %s", root)
	}

	// The globe terrain and base graphics are optional, they are only
	// available if root is the game folder (with GEODATA and GEOGRAPH).
	w, err := loadWorld(root)
	if err != nil {
		log.Printf("Globe terrain not available (%v), using a generic map instead", err)
		w = nil
	} else {
		log.Printf("Loaded globe terrain: %d polygons", len(w.Polygons))
	}
	f, err := loadFonts(root)
	if err != nil {
		log.Printf("Game fonts not available (%v)", err)
		f = map[string][]byte{}
	} else {
		log.Printf("Loaded game fonts")
	}
	bits, err := loadBaseBits(root)
	if err != nil {
		log.Printf("Base graphics not available (%v)", err)
		bits = nil
	} else {
		log.Printf("Loaded base graphics: %d sprites", len(bits))
	}
	images, err := loadSoldierImages(root)
	if err != nil {
		log.Printf("Soldier images not available (%v)", err)
		images = nil
	} else {
		log.Printf("Loaded soldier images: %d images", len(images))
	}

	ter, err := loadTerrain(root)
	if err != nil {
		log.Printf("Battlescape graphics not available (%v)", err)
		ter = nil
	} else {
		log.Printf("Loaded Battlescape graphics: %d terrain sets", len(ter))
	}

	units, err := loadUnitSheets(root)
	if err != nil {
		log.Printf("Unit graphics not available (%v)", err)
		units = nil
	} else {
		log.Printf("Loaded unit graphics: %d sprite files", len(units))
	}

	stateMu.Lock()
	defer stateMu.Unlock()
	games, currentRoot, world, fonts, baseBits, soldierImages, terrain, unitSheets = loaded, root, w, f, bits, images, ter, units
	return nil
}

// lockState holds a read lock on the loaded savegames and resources while
// a request is handled, so loadRoot cannot replace them in the meantime.
func lockState(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		stateMu.RLock()
		defer stateMu.RUnlock()
		next.ServeHTTP(w, r)
	})
}

// frontendAssets returns the embedded build of the frontend.
func frontendAssets() (fs.FS, error) {
	assets, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		return nil, fmt.Errorf("could not access embedded frontend: %w", err)
	}
	return assets, nil
}

// newRouter returns the handler serving the API and the embedded frontend.
func newRouter() (*mux.Router, error) {
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(lockState)
	api.HandleFunc("/games", handleListGames).Methods("GET")
	api.HandleFunc("/games/{slot}", handleGetGame).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers", handleListSoldiers).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers/{idx}", handleGetSoldier).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers/{idx}", handleUpdateSoldier).Methods("PUT")
	api.HandleFunc("/games/{slot}/soldiers/{idx}/resurrect", handleResurrectSoldier).Methods("POST")
	api.HandleFunc("/games/{slot}/soldiers/{idx}/heal", handleHealSoldier).Methods("POST")
	api.HandleFunc("/games/{slot}/bases", handleListBases).Methods("GET")
	api.HandleFunc("/games/{slot}/bases/{idx}", handleGetBase).Methods("GET")
	api.HandleFunc("/games/{slot}/bases/{idx}", handleUpdateBase).Methods("PUT")
	api.HandleFunc("/games/{slot}/craft", handleListCraft).Methods("GET")
	api.HandleFunc("/games/{slot}/transfers", handleListTransfers).Methods("GET")
	api.HandleFunc("/games/{slot}/research", handleListResearch).Methods("GET")
	api.HandleFunc("/games/{slot}/research/{base}/{topic}/speedup", handleSpeedupResearchProject).Methods("POST")
	api.HandleFunc("/games/{slot}/production", handleListProduction).Methods("GET")
	api.HandleFunc("/games/{slot}/production/{base}/{product}/speedup", handleSpeedupProductionProject).Methods("POST")
	api.HandleFunc("/games/{slot}/locations", handleListLocations).Methods("GET")
	api.HandleFunc("/games/{slot}/battle", handleGetBattleMap).Methods("GET")
	api.HandleFunc("/games/{slot}/battle/terrain", handleGetBattleTerrain).Methods("GET")
	api.HandleFunc("/games/{slot}/financials", handleGetFinancials).Methods("GET")
	api.HandleFunc("/games/{slot}/financials", handleUpdateFinancials).Methods("PUT")
	api.HandleFunc("/games/{slot}/actions/heal-all", handleHealAll).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/personal-armour", handlePersonalArmour).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/repair-refuel-crafts", handleRepairRefuelCrafts).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/complete-constructions", handleCompleteConstructions).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/speedup-deliveries", handleSpeedupDeliveries).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/speedup-research", handleSpeedupResearch).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/speedup-production", handleSpeedupProduction).Methods("POST")
	api.HandleFunc("/games/{slot}/save", handleSave).Methods("POST")
	api.HandleFunc("/games/{slot}/reload", handleReload).Methods("POST")
	api.HandleFunc("/world", handleGetWorld).Methods("GET")
	api.HandleFunc("/world/textures/{set}/{terrain}.png", handleGetTexture).Methods("GET")
	api.HandleFunc("/soldier-images/{name}.png", handleGetSoldierImage).Methods("GET")
	api.HandleFunc("/basebits/{sprite}.png", handleGetBaseBit).Methods("GET")
	api.HandleFunc("/terrain/{name}.png", handleGetTerrainSheet).Methods("GET")
	api.HandleFunc("/unitsprites/{name}.png", handleGetUnitSheet).Methods("GET")
	api.HandleFunc("/fonts/{name}.ttf", handleGetFont).Methods("GET")

	// SPA fallback: serve frontend for all non-API routes
	frontendContent, err := frontendAssets()
	if err != nil {
		return nil, err
	}
	fileServer := http.FileServer(http.FS(frontendContent))
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the path has no extension, serve index.html (SPA routing)
		p := r.URL.Path
		if p != "/" && !strings.Contains(path.Base(p), ".") {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
	return r, nil
}
