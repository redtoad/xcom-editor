meta:
  id: smokref_dat
  title: Format of SMOKREF.DAT
  file-extension: DAT
  endian: le
doc: |
  SMOKREF.DAT holds the smoke and fire on the Battlescape map.

  The file is 3600 bytes long (400 entries * 9 bytes).
doc-ref: https://www.ufopaedia.org/index.php/SMOKREF.DAT
seq:
  - id: clouds
    type: smoke_cloud
    repeat: expr
    repeat-expr: 400
types:
  smoke_cloud:
    doc: Smoke or fire on a single tile.
    seq:
      - id: y
        type: u1
      - id: x
        type: u1
      - id: z
        type: u1
      - id: smoke
        type: u1
        doc: Smoke intensity (and remaining time).
      - id: fire
        type: u1
        doc: Remaining time of the fire.
      - id: spread
        type: u1
        doc: Intensity of new spread (only used at runtime).
      - id: type
        type: u1
        enum: cloud_type
        doc: "Type: 0 = empty, 1 = fire, 2 = smoke."
      - id: visible
        type: u1
        doc: Visibility flag.
      - id: frame
        type: u1
        doc: Animation frame offset (0-3).
enums:
  cloud_type:
    0: empty
    1: fire
    2: smoke
