meta:
  id: loftemps_dat
  title: Format of LOFTEMPS.DAT
  file-extension: DAT
  endian: le
doc: |
  LOFTEMPS.DAT is a game data file found in the game's GEODATA/ directory (it
  is not part of a savegame). It contains "Line of Fire Templates", used to
  create 3D representations of the creatures and objects within the
  battlescape. These are used for Line of sight (LOS) (to see units, terrain
  uses MCD[31]) and line of fire (LOF) purposes.

  The file is a sequence of 16x16 bitmaps of 32 bytes each, one bit per
  pixel (1 is solid, 0 is empty). UFO has 112 templates (3,584 bytes).

  `internal/geodata` has no Go struct for this file yet (only a package
  comment), so this layout follows UFOpaedia. Each template is read as 16
  little-endian 16-bit rows; OpenXcom tests bit (15 - x) of row y for the
  voxel at (x, y), i.e. the most significant bit is the leftmost column.
doc-ref: https://www.ufopaedia.org/index.php/LOFTEMPS.DAT
seq:
  - id: templates
    type: template
    repeat: eos
types:
  template:
    seq:
      - id: rows
        type: u2
        repeat: expr
        repeat-expr: 16
        doc: One 16-bit row per y coordinate, 1 bits are solid.
