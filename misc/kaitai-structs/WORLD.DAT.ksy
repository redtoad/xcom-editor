meta:
  id: world_dat
  title: Format of WORLD.DAT
  file-extension: DAT
  endian: le
doc: |
  WORLD.DAT is a game data file found in the game's GEODATA/ directory (it is
  not part of a savegame). It describes the terrain on the geoscape screen
  using quadrilateral polygons and triangles.

  The file is a sequence of 20-byte records. The first 16 bytes of a record
  contain the points of the polygon: 4 sets of 2 short (2-byte) integers,
  designating the 'X' and 'Y' coordinate (or longitude and latitude
  respectively, if you prefer). If the last set has an x value of -1 then it
  is to be rendered as a triangle, otherwise it is a quad.

  The last 4 bytes in the record contain the terrain type. This could be a
  long integer or 2 short integers as the last 2 bytes in each record are 0.

  Mirrors `geodata.WorldFile` / `geodata.Polygon` in `internal/geodata`.
doc-ref: https://www.ufopaedia.org/index.php/WORLD.DAT
seq:
  - id: polygons
    type: polygon
    size: 20
    repeat: eos
enums:
  polygon_type:
    0: triangle
    1: quadrilateral_polygon
types:
  polygon:
    seq:
      - id: x0
        type: s2
        doc: First X coordinate/longitude
      - id: y0
        type: s2
        doc: First Y coordinate/latitude
      - id: x1
        type: s2
        doc: Second X coordinate/longitude
      - id: y1
        type: s2
        doc: Second Y coordinate/latitude
      - id: x2
        type: s2
        doc: Third X coordinate/longitude
      - id: y2
        type: s2
        doc: Third Y coordinate/latitude
      - id: x3
        type: s2
        doc: Fourth X coordinate/longitude, -1 if the polygon is a triangle
      - id: y3
        type: s2
        doc: Fourth Y coordinate/latitude (ignored for triangles)
      - id: terrain
        type: s4
        doc: |
          Terrain type/texture 0-12, used as index into the texture sets of
          GEOGRAPH/TEXTURE.DAT.
    instances:
      polygon_type:
        value: 'x3 == -1 ? 0 : 1'
        enum: polygon_type
        doc: Triangle if the fourth X coordinate is -1, quadrilateral otherwise.
