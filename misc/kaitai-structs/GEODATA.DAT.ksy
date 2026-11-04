meta:
  id: geodata_dat
  title: Format of GEODATA.DAT (Battlescape savegame file)
  file-extension: DAT
  endian: le
doc: |
  GEODATA.DAT is a Battlescape savegame file (not to be confused with the
  game's GEODATA/ directory) and describes how the map was put together from
  map modules. Maps consist of up to 6x6 modules of 10x10 tiles.

  The file is 92 bytes long.
doc-ref: https://www.ufopaedia.org/index.php/GEODATA.DAT
seq:
  - id: modules_y
    type: u2
    doc: Number of modules along the Y axis (map size / 10).
  - id: modules_x
    type: u2
    doc: Number of modules along the X axis (map size / 10).
  - id: size_z
    type: u2
    doc: Number of map levels.
  - id: xcom_craft
    type: u2
    doc: X-COM craft type (ignored for base missions).
  - id: alien_craft
    type: u2
    doc: Alien craft type (ignored if not a UFO mission).
  - id: terrain
    type: u2
    doc: Terrain type.
  - id: modules
    type: u1
    repeat: expr
    repeat-expr: 36
    doc: |
      Module IDs, west to east and north to south. 253-255 mark modules
      covered by a 2x2 module.
  - id: load_routes
    type: u1
    repeat: expr
    repeat-expr: 36
    doc: 1 if the routes (RMP) of the module are loaded, 0 otherwise.
  - id: crafts
    type: u1
    enum: crafts_on_map
    doc: "Crafts on the map: 0 = none, 1 = X-COM only, 2 = X-COM and alien."
  - id: xcom_craft_copy
    type: u1
    doc: Copy of xcom_craft.
  - id: xcom_craft_y
    type: u1
    doc: Y position of the X-COM craft.
  - id: xcom_craft_x
    type: u1
    doc: X position of the X-COM craft.
  - id: alien_craft_copy
    type: u1
    doc: Copy of alien_craft.
  - id: alien_craft_y
    type: u1
    doc: Y position of the alien craft.
  - id: alien_craft_x
    type: u1
    doc: X position of the alien craft.
  - id: unknown
    type: u1
    doc: Unknown, always 0.
enums:
  crafts_on_map:
    0: none
    1: xcom_only
    2: xcom_and_alien
