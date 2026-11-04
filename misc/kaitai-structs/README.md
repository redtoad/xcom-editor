# Kaitai format descriptions

[Kaitai Struct](http://kaitai.io/#quick-start) is a declarative language used to describe various binary data structures, laid out in files or in memory: i.e. binary file formats, network stream packet formats, etc.

You can try out any of the files in this directory using the [Kaitai Web IDE](https://ide.kaitai.io).

The descriptions mirror the Go structs in `internal/` (which are the source of
truth, as they round-trip real savegames byte-identically):

- Geoscape savegame files (`internal/geoscape`): ACTS, AKNOW, ALIEN, ASTORE,
  BASE, BPROD, CRAFT, DIPLOM, FACIL, IGLOB, INTER, LEASE, LIGLOB, LOC,
  MISSIONS, PRODUCT, PROJECT, PURCHASE, RESEARCH, SAVEINFO, SOLDIER, TRANSFER,
  UIGLOB, UP, XBASES, XCOM and ZONAL.
- Battlescape savegame files (`internal/battlescape`): BGLOB, GEODATA, MAP,
  MISDATA, OBPOS, ROUTES, SEEMAP, SMOKBIT, SMOKREF, SOURCEMP, TERMP, UNITPOS,
  UNITREF and WGLOB.
- Game data files (`internal/geodata`): INTERWIN, LOFTEMPS, OBDATA, SCANG and
  WORLD from the game's `GEODATA/` directory, TEXTURE from `GEOGRAPH/`. Do not
  confuse `GEODATA.DAT.ksy` (a Battlescape savegame file) with the `GEODATA/`
  directory. INTERWIN, LOFTEMPS, OBDATA and SCANG have no Go struct yet; their
  descriptions are based on UFOpaedia and checked against the original files.

When changing a Go struct, update the matching `.ksy` file as well.

## Verification

`verify.py` compiles every `.ksy` file to Python and parses all files of the
same name in `testdata/`. Each file must parse without errors and all bytes
must be consumed. Files without samples in `testdata/` are only compiled.

```bash
brew install kaitai-struct-compiler
pip install kaitaistruct
python3 misc/kaitai-structs/verify.py              # all descriptions
python3 misc/kaitai-structs/verify.py CRAFT.DAT.ksy  # a single one
```
