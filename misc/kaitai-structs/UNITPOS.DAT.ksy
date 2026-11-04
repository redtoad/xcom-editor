meta:
  id: unitpos_dat
  title: Format of UNITPOS.DAT
  file-extension: DAT
  endian: le
doc: |
  UNITPOS.DAT holds the positions and owners of all units on the map. Each
  entry belongs to the entry with the same index in UNITREF.DAT.

  The file is 1120 bytes long (80 units * 14 bytes).

  MISSDAT/UNIPOS.DAT (written by GEOSCAPE.EXE before a mission, only the
  units used) uses the same record format.
doc-ref: https://www.ufopaedia.org/index.php/UNITPOS.DAT
seq:
  - id: units
    type: unit_position
    repeat: expr
    repeat-expr: 80
types:
  unit_position:
    seq:
      - id: y
        type: u1
        doc: |
          Position. Unused slots (unit type is 0xFF in UNITREF.DAT) may
          contain arbitrary values here.
      - id: x
        type: u1
      - id: z
        type: u1
      - id: unknown1
        type: u1
        doc: Unknown, always 0.
      - id: unit_ref
        type: u4
        doc: |
          Memory address of the UNITREF.DAT entry (not used on disk;
          TACTICAL.EXE recalculates it after loading MISSDAT/UNIPOS.DAT; the
          four entries of a large unit all point to the same UNITREF.DAT
          entry).
      - id: turns_since_spotted
        type: u1
        doc: Turns since the unit was last spotted.
      - id: owner
        type: u1
        enum: owner
        doc: |
          Owner: 0 = X-COM, 1 = alien, 2 = civilian or 3 for X-COM units left
          behind on a previous stage.
      - id: flags
        type: unit_flags
        doc: Flags.
      - id: unknown2
        size: 3
        doc: Unknown, always 0.
  unit_flags:
    seq:
      - id: damaged
        type: b1
        doc: Bit 7.
      - id: mind_controlled
        type: b1
        doc: Bit 6.
      - id: unknown
        type: b1
        doc: Bit 5, unknown.
      - id: lower_section
        type: b1
        doc: Bit 4, lower section of a large unit.
      - id: right_section
        type: b1
        doc: Bit 3, right section of a large unit.
      - id: large
        type: b1
        doc: Bit 2, large (2x2) unit.
      - id: active
        type: b1
        doc: Bit 1, unit is active (e.g. not dead).
      - id: visible
        type: b1
        doc: Bit 0, unit is visible.
enums:
  owner:
    0: xcom
    1: alien
    2: civilian
    3: xcom_left_behind
