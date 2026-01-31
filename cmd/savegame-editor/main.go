package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
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
}

// games holds all loaded savegames indexed by slot name (e.g. "GAME_1").
var games map[string]*gameEntry

// OpenURL opens the specified URL in the default browser.
func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port to run server on")
	flag.Parse()

	root := flag.Arg(0)
	if root == "" {
		fmt.Fprintln(os.Stderr, "Usage: savegame-editor [-port PORT] <savegame-root-folder>")
		os.Exit(1)
	}

	log.Printf("Starting savegame editor server...")
	log.Printf("Version %s-%s %s %s", version, commit, date, builtBy)
	log.Printf("Savegame root: %s", root)

	// Scan for GAME_1 through GAME_10
	games = make(map[string]*gameEntry)
	for i := 1; i <= 10; i++ {
		slot := fmt.Sprintf("GAME_%d", i)
		gamePath := filepath.Join(root, slot)
		if _, err := os.Stat(filepath.Join(gamePath, "SAVEINFO.DAT")); os.IsNotExist(err) {
			continue
		}
		sg, err := savegame.Load(gamePath)
		if err != nil {
			log.Printf("Warning: could not load %s: %v", slot, err)
			continue
		}
		games[slot] = &gameEntry{sg: sg}
		log.Printf("Loaded %s: %s", slot, sg.Title())
	}

	if len(games) == 0 {
		log.Fatal("No savegames found in ", root)
	}

	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/games", handleListGames).Methods("GET")
	api.HandleFunc("/games/{slot}", handleGetGame).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers", handleListSoldiers).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers/{idx}", handleGetSoldier).Methods("GET")
	api.HandleFunc("/games/{slot}/soldiers/{idx}", handleUpdateSoldier).Methods("PUT")
	api.HandleFunc("/games/{slot}/bases", handleListBases).Methods("GET")
	api.HandleFunc("/games/{slot}/bases/{idx}", handleGetBase).Methods("GET")
	api.HandleFunc("/games/{slot}/bases/{idx}", handleUpdateBase).Methods("PUT")
	api.HandleFunc("/games/{slot}/craft", handleListCraft).Methods("GET")
	api.HandleFunc("/games/{slot}/transfers", handleListTransfers).Methods("GET")
	api.HandleFunc("/games/{slot}/financials", handleGetFinancials).Methods("GET")
	api.HandleFunc("/games/{slot}/financials", handleUpdateFinancials).Methods("PUT")
	api.HandleFunc("/games/{slot}/actions/heal-all", handleHealAll).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/complete-constructions", handleCompleteConstructions).Methods("POST")
	api.HandleFunc("/games/{slot}/actions/speedup-deliveries", handleSpeedupDeliveries).Methods("POST")
	api.HandleFunc("/games/{slot}/save", handleSave).Methods("POST")
	api.HandleFunc("/games/{slot}/reload", handleReload).Methods("POST")

	// SPA fallback: serve frontend for all non-API routes
	frontendContent, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal("Could not access embedded frontend: ", err)
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

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Printf("Server running on http://localhost:%s", port)
	log.Println("Opening browser...")
	_ = OpenURL("http://localhost:" + port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	log.Println("Press Ctrl+C to stop.")
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Println("Shutting down.")
	os.Exit(0)
}
