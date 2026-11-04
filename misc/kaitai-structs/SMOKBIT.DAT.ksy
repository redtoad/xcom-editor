meta:
  id: smokbit_dat
  title: Format of SMOKBIT.DAT
  file-extension: DAT
  endian: le
doc: |
  SMOKBIT.DAT is a bit field with one bit per tile (in the same order as in
  MAP.DAT, least significant bit first: tile i is bit i % 8 of byte i / 8).
  A set bit means that the tile has smoke or fire (see SMOKREF.DAT). The game
  recalculates it every turn.

  The file is 1250 bytes (10000 bits) long, regardless of the map size.
doc-ref: https://www.ufopaedia.org/index.php/SMOKBIT.DAT
seq:
  - id: bits
    size: 1250
    doc: One bit per tile, least significant bit first.
