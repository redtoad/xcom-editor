meta:
  id: scang_dat
  title: Format of SCANG.DAT
  file-extension: DAT
  endian: le
doc: |
  SCANG.DAT is a game data file found in the game's GEODATA/ directory (it
  is not part of a savegame). It is used along with SCANBORD.PCK to present
  the overhead map. The title presumably stands for SCANner Graphics.

  The file holds 4x4 images of 16 bytes each (one palette index per pixel,
  row by row); UFO has 454 of them (7,264 bytes). The first 36 images are
  units and equipment; terrain parts select their image via
  MCD[20] + MCD[21] * 256 + 35.

  `internal/geodata` has no Go struct for this file yet (only a package
  comment), so this layout follows UFOpaedia.
doc-ref: https://www.ufopaedia.org/index.php/SCANG.DAT
seq:
  - id: images
    type: image
    repeat: eos
types:
  image:
    seq:
      - id: pixels
        size: 4 * 4
        doc: 4x4 palette indexes, row by row.
