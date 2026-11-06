
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
savegames saved during a mission (available as `Savegame.Battle`; Battlescape files
left over from an earlier mission are ignored, see `TacticalSave` in `SAVEINFO.DAT`). Every file is
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
units and score). For those savegames the *Mission* tab shows a map of the
Battlescape, either top-down level by level or isometric like in the game
(levels above the selected one are hidden). If the folder is the game folder,
the isometric map is drawn with the terrain graphics of the game
(`TERRAIN/*.MCD`, `*.PCK`); as the savegame does not name the terrain files of
the map, they are determined by matching the tile types of the MCD records
against the map (see `savegame.Battle.MatchTerrainSets`). Soldiers, tanks
(with their turret), Sectoids, Snakemen, Mutons, Floaters, Chryssalids, Reapers,
Cyberdiscs and civilians are drawn with their sprites (`UNITS/*.PCK`), other
units as markers. Dead Reapers are drawn with their corpse (`UNITS/FLOOROB.PCK`),
which lies on all 4 tiles of the unit. The sprites use the
Battlescape palette of `PALETTES.DAT` with the last 16 colors replaced by the
grey gradient the game sets at runtime (see `resources.BattlescapePalette`). Otherwise the map is drawn
with simple shapes. The map shows the floors, walls and objects of the
map (`MAP.DAT`), the fog of war (`SEEMAP.DAT`), all soldiers, aliens and
civilians (`UNITPOS.DAT`, `UNITREF.DAT`) and the area of the X-COM craft and
the UFO (`GEODATA.DAT`, the size of the crafts is approximate).

Edits to soldiers, bases and the balance are applied with "Apply Changes" (or
"Update"); "Save" writes the savegame to disk and also applies edits in an
open form that were not applied yet, while "Reload" refreshes open forms with
the values read from disk. It has a map of both hemispheres of the globe with all
bases, X-COM craft and UFOs, their flight paths and details on hover. Both the
map and the 3D globe show the detection range of the radars of each base
(small radar 1696 nm, large radar 2578 nm, hyperwave decoder 2759 nm, as
used by the detection routine of `GEOSCAPE.EXE`; see `savegame.Base.Radars`),
only completed facilities count. If the
folder is the game folder (containing `GEODATA` and `GEOGRAPH`), the map is
drawn with the terrain and textures of the game (`WORLD.DAT`, `TEXTURE.DAT`)
and bases are shown like on the base screen of the game, including the crafts
in their hangars (`BASEBITS.PCK`), and each soldier is shown like on the
inventory screen, without armour next to the selected armour (`MAN_*.SPK`),
otherwise a generic world map and a simple facility grid are used.
Like on the base screen of the game, bases are selected with the miniatures of
all eight base slots at the top (built facilities in green, facilities under
construction in red); edits not applied yet are applied when switching bases.
The editor watches the savegame folder: savegames that are saved in the game
appear in the sidebar, and if a loaded savegame is overwritten (e.g. saved
again in the game), it is marked as changed and can be reloaded.
Each tab offers quick actions: on *Soldiers* "Heal All" and "All With Personal
Armour" (equips every living soldier without armour), a "Heal" button for
each wounded soldier and a "Resurrect" button for each dead soldier, on
*Bases* "Complete
Constructions" and "-"/"+" buttons to change the quantity of each item in
the stores of a base (by ten with the shift key; "Show all items" lists items
not in stores, so they can be added; the stores capacity is not checked), on
*Craft* "Repair & Refuel All" (repairs, refuels and rearms all
X-COM craft immediately and puts enough ammunition into the stores of their
bases to rearm them once more), on *Research* and *Manufacture* (which list
the ongoing research and production of all bases) a "Speed Up" button for
each project and "Speed Up All Research"/"Speed Up All Production" (research
is completed on the next day, the next unit of an item within the next hour;
the remaining units are built at the normal pace; projects that are finished
by then anyway are shown as "Ready", and projects without scientists or
engineers, which do not progress, are marked red), on *Transfers* "Speed Up
Deliveries", and on *Financials* a
"Max" button sets the current balance to its maximum ($2,147,483,647). Dead
soldiers come back with the rank they had before. The game does not keep it,
so it is estimated from the soldier value (see `savegame.Soldier.LastRank`;
Rookies and Squaddies cannot be told apart and are revived as Rookies). The user
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

It uses [Wails v3](https://v3.wails.io). The build is done by the tasks in
`cmd/savegame-editor/Taskfile.yml` and `cmd/savegame-editor/build/` (generated
with `wails generate build-assets`), which call the CLI as `wails`, as
installed by Homebrew (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`
installs it as `wails3`, which then has to be linked as `wails`):

```sh
brew install wails
cd cmd/savegame-editor
wails build                   # builds frontend and app into bin/
wails package                 # additionally creates bin/savegame-editor.app on macOS
wails dev                     # live reload of the frontend and the Go code
```

The tasks add the build tag `desktop` (`EXTRA_TAGS` in `Taskfile.yml`).
Platform settings such as the bundle identifier, version and minimum macOS
version are in `build/config.yml` and `build/darwin/Info.plist` (and the
other platform folders). `wails dev` builds the frontend unminified into
`frontend/dist`, which is checked in, so restore it afterwards
(`git checkout frontend/dist`).

The app icon is `cmd/savegame-editor/build/appicon.png`, rendered from
`appicon.svg` next to it (`inkscape appicon.svg -w 1024 -h 1024 -o appicon.png`);
the platform icons are generated from it during the build.

The desktop build needs cgo and the platform's webview (WebKit on macOS and
Linux, WebView2 on Windows), so it is not cross-compiled by goreleaser.
