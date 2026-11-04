//go:build desktop || dev

package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	goruntime "runtime"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// desktopApp runs the savegame editor in a native window (using Wails)
// instead of the default browser. The frontend talks to the same API, which
// is served through the asset server of the webview.
type desktopApp struct {
	ctx  context.Context
	root string
	// ready is closed as soon as the savegames have been loaded, API
	// requests are held back until then.
	ready chan struct{}
}

// startup is called by Wails once the window has been created.
func (a *desktopApp) startup(ctx context.Context) {
	a.ctx = ctx
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
			wailsruntime.Quit(ctx)
			return
		}
		close(a.ready)
	}()
}

// chooseRoot asks for a savegame folder until one could be loaded. It
// returns false if the user cancelled the dialog.
func (a *desktopApp) chooseRoot() bool {
	for {
		root, err := wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
			Title:            "Choose the X-COM game folder (containing GAME_1 ... GAME_10)",
			DefaultDirectory: a.root,
		})
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
func (a *desktopApp) openFolder(*menu.CallbackData) {
	select {
	case <-a.ready:
	default:
		return // the startup dialog is still open
	}
	if a.chooseRoot() {
		wailsruntime.WindowReloadApp(a.ctx)
	}
}

func (a *desktopApp) showError(err error) {
	log.Print(err)
	_, _ = wailsruntime.MessageDialog(a.ctx, wailsruntime.MessageDialogOptions{
		Type:    wailsruntime.ErrorDialog,
		Title:   "Could not load savegames",
		Message: err.Error(),
	})
}

// serveAPI returns a middleware for the asset server which passes API
// requests to api (once the savegames have been loaded) and everything else
// on to the frontend assets, which are served by the Vite dev server when
// running with `wails dev`.
func (a *desktopApp) serveAPI(api http.Handler) assetserver.Middleware {
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

func (a *desktopApp) menu() *menu.Menu {
	m := menu.NewMenu()
	if goruntime.GOOS == "darwin" {
		m.Append(menu.AppMenu())
	}
	file := m.AddSubmenu("File")
	file.AddText("Open Game Folder…", keys.CmdOrCtrl("o"), a.openFolder)
	file.AddSeparator()
	file.AddText("Quit", keys.CmdOrCtrl("q"), func(*menu.CallbackData) {
		wailsruntime.Quit(a.ctx)
	})
	if goruntime.GOOS == "darwin" {
		// Needed for copy & paste in text fields on macOS.
		m.Append(menu.EditMenu())
		m.Append(menu.WindowMenu())
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

	app := &desktopApp{root: flag.Arg(0), ready: make(chan struct{})}

	r, err := newRouter()
	if err != nil {
		log.Fatal(err)
	}
	assets, err := frontendAssets()
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            "X-COM Savegame Editor",
		Width:            1280,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 255},
		AssetServer: &assetserver.Options{
			Assets:     assets,
			Middleware: app.serveAPI(r),
		},
		Menu:      app.menu(),
		OnStartup: app.startup,
	})
	if err != nil {
		log.Fatal(err)
	}
}
