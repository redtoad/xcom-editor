meta:
  id: map_dat
  title: Format of MAP.DAT
  file-extension: DAT
  endian: le
doc: |
  MAP.DAT holds the tiles of the Battlescape map, 4 bytes per tile. Tiles are
  ordered along X, then Y, then from the top level to the bottom level.

  The used part of the file is size_x * size_y * size_z * 4 bytes, with the
  map dimensions taken from WGLOB.DAT (or modules_x * 10, modules_y * 10 and
  size_z from GEODATA.DAT). The game does not truncate the file, so stale
  tiles of a previous, larger map may follow the used part. As a single
  description cannot read WGLOB.DAT, all records up to the end of the file
  are parsed.
doc-ref: https://www.ufopaedia.org/index.php/MAP.DAT
seq:
  - id: tiles
    type: map_tile
    repeat: eos
types:
  map_tile:
    doc: |
      The parts of a single tile. Each part is an index into the MCD records
      of the map, 0 if the part does not exist.
    seq:
      - id: ground
        type: u1
      - id: west_wall
        type: u1
      - id: north_wall
        type: u1
      - id: object
        type: u1
