meta:
  id: lease_dat
  title: Format of LEASE.DAT
  file-extension: DAT
  endian: le
doc: |
  LEASE.DAT holds the view of the Geoscape globe. Rotations are stored in
  1/8 degree (0-2880 = 0-360 degrees).

  Only the first bytes are documented, the rest is kept as raw bytes. The
  file is 180 bytes long.

  GEOSCAPE.EXE treats the file as 30 records of 6 bytes (uint16 + uint32),
  apparently left over from a craft leasing feature: a new game only clears
  the uint16 of each record, so the uint32 parts keep uninitialized memory
  (which looks like x86 machine code). When saving, the globe view is written
  into the uint16 of the first three records (offsets 0, 6 and 12). The
  purchase screen writes (36, buy price) into a free record (uint16 = 0) when
  crafts are bought, but due to a bug it only checks the record indexed by
  the number of crafts bought so far. Nothing else reads the records.

  References in GEOSCAPE.EXE (linear addresses): allocation 0x216d3, new game
  0x21c9b, save 0x14970, load 0x1585b, purchase screen 0x1c6f9.
doc-ref: https://www.ufopaedia.org/index.php/LEASE.DAT
seq:
  - id: rotation_y
    type: s2
    doc: Y rotation of the globe (uint16 of record 0).
  - id: unknown1
    size: 4
    doc: Unknown (uint32 of record 0, uninitialized memory).
  - id: rotation_x
    type: s2
    doc: X rotation of the globe (uint16 of record 1).
  - id: unknown2
    size: 4
    doc: Unknown (uint32 of record 1, uninitialized memory).
  - id: zoom
    type: s2
    doc: |
      Zoom (globe radius in pixels: 90, 120, 180, 360, 450 or 720; uint16 of
      record 2).
  - id: unknown3
    size: 166
    doc: |
      Unknown: uint32 of record 2 (uninitialized memory), followed by records
      3-29. Their uint16 is 0 (free) unless crafts were bought, then it may be
      36 followed by the buy price as uint32.
