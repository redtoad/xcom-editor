meta:
  id: seemap_dat
  title: Format of SEEMAP.DAT
  file-extension: DAT
  endian: le
doc: |
  SEEMAP.DAT holds the fog of war: which parts of each tile X-COM has seen.
  One byte per tile, in the same order as in MAP.DAT.

  The used part of the file is size_x * size_y * size_z bytes, with the map
  dimensions taken from WGLOB.DAT (or GEODATA.DAT). The game does not
  truncate the file, so stale data of a previous, larger map may follow the
  used part. As a single description cannot read WGLOB.DAT, all bytes up to
  the end of the file are parsed as tiles.
doc-ref: https://www.ufopaedia.org/index.php/SEEMAP.DAT
seq:
  - id: tiles
    type: seen_flags
    repeat: eos
types:
  seen_flags:
    seq:
      - id: reserved
        type: b5
        doc: Unused bits 3-7.
      - id: content
        type: b1
        doc: Bit 2 (SeenContent) - the tile content has been seen.
      - id: north_wall
        type: b1
        doc: Bit 1 (SeenNorthWall) - the north wall has been seen.
      - id: west_wall
        type: b1
        doc: Bit 0 (SeenWestWall) - the west wall has been seen.
