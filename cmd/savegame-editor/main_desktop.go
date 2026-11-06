//go:build desktop

package main

import (
	"flag"
	"log"
	"net/http"
	goruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// desktopApp runs the savegame editor in a native window (using Wails)
// instead of the default browser. The frontend talks to the same API, which
// is served through the asset server of the webview.
type desktopApp struct {
	app    *application.App
	window *application.WebviewWindow
	root   string
	// ready is closed as soon as the savegames have been loaded, API
	// requests are held back until then.
	ready chan struct{}
}

// startup is called by Wails once the application has been started.
func (a *desktopApp) startup(*application.ApplicationEvent) {
	// pick up savegames which are saved (or overwritten) in the game
	go watchGames(watchInterval)
	go func() {
		if a.root != "" {
			err := loadRoot(a.root)
			if err == nil {
				close(a.ready)
				return
			}
			a.showError(err)
		}
		if !a.chooseRoot() {
			a.app.Quit()
			return
		}
		close(a.ready)
	}()
}

// chooseRoot asks for a savegame folder until one could be loaded. It
// returns false if the user cancelled the dialog.
func (a *desktopApp) chooseRoot() bool {
	for {
		root, err := a.app.Dialog.OpenFile().
			SetTitle("Choose the X-COM game folder (containing GAME_1 ... GAME_10)").
			SetMessage("Choose the X-COM game folder (containing GAME_1 ... GAME_10)").
			CanChooseFiles(false).
			CanChooseDirectories(true).
			SetDirectory(a.root).
			AttachToWindow(a.window).
			PromptForSingleSelection()
		if err != nil || root == "" {
			return false
		}
		log.Printf("Savegame root: %s", root)
		if err := loadRoot(root); err != nil {
			a.showError(err)
			continue
		}
		a.root = root
		return true
	}
}

// openFolder lets the user switch to another savegame folder and reloads the
// frontend afterwards.
func (a *desktopApp) openFolder(*application.Context) {
	select {
	case <-a.ready:
	default:
		return // the startup dialog is still open
	}
	go func() {
		if a.chooseRoot() {
			a.window.Reload()
		}
	}()
}

func (a *desktopApp) showError(err error) {
	log.Print(err)
	a.app.Dialog.Error().
		SetTitle("Could not load savegames").
		SetMessage(err.Error()).
		AttachToWindow(a.window).
		Show()
}

// serveAPI returns a middleware for the asset server which passes API
// requests to api (once the savegames have been loaded) and everything else
// on to the frontend assets, which are served by the Vite dev server when
// running with `wails dev`.
func (a *desktopApp) serveAPI(api http.Handler) application.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api/") {
				next.ServeHTTP(w, r)
				return
			}
			select {
			case <-a.ready:
			case <-r.Context().Done():
				return
			}
			api.ServeHTTP(w, r)
		})
	}
}

func (a *desktopApp) menu() *application.Menu {
	m := application.NewMenu()
	if goruntime.GOOS == "darwin" {
		m.AddRole(application.AppMenu)
	}
	file := m.AddSubmenu("File")
	file.Add("Open Game Folder…").SetAccelerator("CmdOrCtrl+O").OnClick(a.openFolder)
	file.AddSeparator()
	file.Add("Quit").SetAccelerator("CmdOrCtrl+Q").OnClick(func(*application.Context) {
		a.app.Quit()
	})
	if goruntime.GOOS == "darwin" {
		// Needed for copy & paste in text fields on macOS.
		m.AddRole(application.EditMenu)
		m.AddRole(application.WindowMenu)
	}
	return m
}

func main() {
	flag.Usage = func() {
		log.Println("Usage: savegame-editor [savegame-root-folder]")
	}
	flag.Parse()

	log.Printf("Starting savegame editor...")
	log.Printf("Version %s-%s %s %s", version, commit, date, builtBy)

	d := &desktopApp{root: flag.Arg(0), ready: make(chan struct{})}

	r, err := newRouter()
	if err != nil {
		log.Fatal(err)
	}
	assets, err := frontendAssets()
	if err != nil {
		log.Fatal(err)
	}

	d.app = application.New(application.Options{
		Name:        "X-COM Savegame Editor",
		Description: "Savegame editor for X-COM: UFO Defense",
		Assets: application.AssetOptions{
			Handler:    application.AssetFileServerFS(assets),
			Middleware: d.serveAPI(r),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})
	d.app.Menu.Set(d.menu())
	d.window = d.app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "X-COM Savegame Editor",
		Width:            1280,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(0, 0, 0),
		URL:              "/",
	})
	d.app.Event.OnApplicationEvent(events.Common.ApplicationStarted, d.startup)

	if err := d.app.Run(); err != nil {
		log.Fatal(err)
	}
}
