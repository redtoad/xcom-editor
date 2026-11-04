meta:
  id: texture_dat
  title: Format of TEXTURE.DAT
  file-extension: DAT
  endian: le
doc: |
  TEXTURE.DAT is a game data file (it is not part of a savegame). Note that,
  unlike the other files in internal/geodata, it lives in the game's
  GEOGRAPH/ directory (GEOGRAPH/TEXTURE.DAT), not in GEODATA/.

  It holds the textures used to draw the terrain of the globe. There are 3
  sets of 13 textures (one per terrain type, see the terrain field of
  WORLD.DAT), one set for each zoom level. Each texture is 32x32 pixels, one
  byte (palette index) per pixel, stored row by row.

  Texture `i` belongs to set `i / 13` and terrain type `i % 13`.

  The file is 39936 bytes long (3 * 13 * 32 * 32).

  Mirrors `geodata.TextureFile` in `internal/geodata`.
doc-ref: https://www.ufopaedia.org/index.php/TEXTURE.DAT
seq:
  - id: textures
    type: texture
    repeat: expr
    repeat-expr: 3 * 13
types:
  texture:
    seq:
      - id: pixels
        size: 32 * 32
        doc: 32x32 palette indexes, row by row.
