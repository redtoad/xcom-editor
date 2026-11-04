
# xcom-editor

[![Go](https://github.com/redtoad/xcom-editor/actions/workflows/go.yml/badge.svg)](https://github.com/redtoad/xcom-editor/actions/workflows/go.yml)

This is a library to load image resources for the game [UFO: Enemy Unknown](https://en.wikipedia.org/wiki/UFO:_Enemy_Unknown). 

## Game images

All images are stored in a custom format (see https://www.ufopaedia.org/index.php/Image_Formats for more information).

Here is an example how you can extract the images from `ZOMBIE.PCK` to create an animated GIF:

```go
package main

import (
    "image/gif"
    "os"
    
    "github.com/redtoad/xcom-editor/resources"
)

func main() {
    // load all game palettes 
    palettes, _ := resources.LoadPalettes("GEODATA/PALETTES.DAT")

    // load image data from .PCK (all offsets are stored in .TAB file)
    offsets, _ := resources.LoadTAB("UNITS/ZOMBIE.TAB", 2)
    collection, _ := resources.LoadImageCollectionFromPCK("UNITS/ZOMBIE.PCK", 32, offsets)
    
    // convert to animated gif
    animated := collection.Animated(25, 40, palettes[4])
    f, _ := os.Create("zombie.gif")
    _ = gif.EncodeAll(f, animated)
}
```

The result is:

![](zombie.gif)






## Savegames

The `savegame` package loads and saves savegames (`GAME_1` ... `GAME_10`). All 27
Geoscape data files are supported, as well as the 14 additional Battlescape files of
savegames saved during a mission (available as `Savegame.Battle`). Every file is
written back byte-identical to the original. See [`internal/geoscape`](internal/geoscape)
and [`internal/battlescape`](internal/battlescape) for the file formats.

```go
sg, _ := savegame.Load("GAME/GAME_1")
sg.Countries.Countries[0].Satisfaction = 3 // make the USA happy
_ = sg.Save()
```

The real savegames in `testdata/` are used for round-trip tests.

## Savegame editor

`cmd/savegame-editor` is a web-based editor for savegames:

```sh
go run ./cmd/savegame-editor path/to/game   # folder containing GAME_1 ... GAME_10
```

It shows whether a savegame was saved during a mission (with turn, surviving
units and score). Edits to soldiers, bases and the balance are applied with
"Apply Changes" (or "Update"); "Save" writes the savegame to disk and also
applies edits in an open form that were not applied yet, while "Reload"
refreshes open forms with the values read from disk. It has a map of both hemispheres of the globe with all
bases, X-COM craft and UFOs, their flight paths and details on hover. If the
folder is the game folder (containing `GEODATA` and `GEOGRAPH`), the map is
drawn with the terrain and textures of the game (`WORLD.DAT`, `TEXTURE.DAT`)
and bases are shown like on the base screen of the game, including the crafts
in their hangars (`BASEBITS.PCK`), and each soldier is shown like on the
inventory screen, without armour next to the selected armour (`MAN_*.SPK`),
otherwise a generic world map and a simple facility grid are used.
The editor watches the savegame folder: savegames that are saved in the game
appear in the sidebar, and if a loaded savegame is overwritten (e.g. saved
again in the game), it is marked as changed and can be reloaded.
Each tab offers quick actions: on *Soldiers* "Heal All" and "All With Personal
Armour" (equips every living soldier without armour), on *Bases* "Complete
Constructions", on *Transfers* "Speed Up Deliveries", and on *Financials* a
"Max" button sets the current balance to its maximum ($2,147,483,647). The user
interface mimics the game and uses its fonts (`SMALLSET.DAT`, `BIGLETS.DAT`),
which the editor converts to TrueType (see `resources.BitmapFont.TTF`). The
frontend (Svelte) is in `cmd/savegame-editor/frontend`; run `npm run build`
there after changing it, the build output is embedded into the binary.

### Desktop application

The editor can also be built as a native desktop application with
[Wails](https://wails.io) (build tag `desktop`). It uses the same frontend and
API, but shows them in a native window instead of the browser. If no folder
is given on the command line, it asks for the game folder on startup; another
folder can be opened via *File → Open Game Folder…*.

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@latest
cd cmd/savegame-editor
wails build -skipbindings   # creates build/bin/ (an .app bundle on macOS)
wails dev -skipbindings -appargs path/to/game   # live reload of the frontend
```

The app icon is `cmd/savegame-editor/build/appicon.png`, rendered from
`appicon.svg` next to it (`inkscape appicon.svg -w 1024 -h 1024 -o appicon.png`).

Without the Wails CLI, build it with `go build -tags desktop,production`
(on macOS additionally with `CGO_LDFLAGS="-framework UniformTypeIdentifiers"`).
The desktop build needs cgo and the platform's webview (WebKit on macOS and
Linux, WebView2 on Windows), so it is not cross-compiled by goreleaser.
