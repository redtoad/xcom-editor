meta:
  id: sourcemp_dat
  title: Format of SOURCEMP.DAT
  file-extension: DAT
  endian: le
doc: |
  SOURCEMP.DAT is a light map holding the light level of each tile, from 0
  (bright) to 15 (dark). Soldiers can only see tiles with a light level of 9
  or less. The effective light level is the lower value of TERMP.DAT (fixed
  and ambient light) and SOURCEMP.DAT (mobile light sources such as units).
  Tiles are stored in the same order as in MAP.DAT.

  The file is one byte per tile plus 4 unknown trailing bytes long. The used
  part of the tiles is size_x * size_y * size_z bytes, with the map
  dimensions taken from WGLOB.DAT (or GEODATA.DAT). The game does not
  truncate the file, so stale data of a previous, larger map may precede the
  trailer. As a single description cannot read WGLOB.DAT, every byte except
  the last 4 is parsed as a tile (as the Go code does).
doc-ref: https://www.ufopaedia.org/index.php/SOURCEMP.DAT
seq:
  - id: tiles
    size: _io.size - 4
    doc: Light level of each tile (one byte per tile).
  - id: trailer
    size: 4
    doc: Unknown trailing bytes.
