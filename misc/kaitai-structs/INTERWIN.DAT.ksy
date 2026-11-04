meta:
  id: interwin_dat
  title: Format of INTERWIN.DAT
  file-extension: DAT
  endian: le
doc: |
  INTERWIN.DAT is a game data file found in the game's GEODATA/ directory (it
  is not part of a savegame).

  Ten images used by the UFO Interception display. Direct palette indexes
  (one byte per pixel, row by row), no compression. 88,960 bytes in total.
  All images are 160 pixels wide: the first is 96 pixels high, the second 44
  and the remaining eight (one per UFO type) 52.

  `internal/geodata` has no Go struct for this file yet (only a package
  comment), so this layout follows UFOpaedia.
doc-ref: https://www.ufopaedia.org/index.php/INTERWIN.DAT
seq:
  - id: window
    type: image(96)
    doc: Interception window background (160x96).
  - id: window_bottom
    type: image(44)
    doc: Second interception window image (160x44).
  - id: ufos
    type: image(52)
    repeat: expr
    repeat-expr: 8
    doc: |
      UFO images (160x52 each), in this order: small scout, medium scout,
      large scout, abductor, harvester, terror ship, battleship, supply ship.
types:
  image:
    params:
      - id: height
        type: u2
    seq:
      - id: pixels
        size: 160 * height
        doc: Palette indexes, 160 pixels per row.
    instances:
      width:
        value: 160
