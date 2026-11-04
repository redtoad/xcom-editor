meta:
  id: obpos_dat
  title: Format of OBPOS.DAT
  file-extension: DAT
  endian: le
doc: |
  OBPOS.DAT holds all objects on the map that can be picked up (e.g. weapons,
  ammunition and corpses), whether on the ground or carried by units.

  The file is 2720 bytes long (170 objects * 16 bytes).

  MISSDAT/OBPOSREF.DAT (written by GEOSCAPE.EXE before a mission, only the
  objects used, e.g. 83 * 16 = 1328 bytes) and MISSDAT/OBPOS2.DAT (written by
  TACTICAL.EXE at the end of a mission) use the same record format.
doc-ref: https://www.ufopaedia.org/index.php/OBPOS.DAT
seq:
  - id: objects
    type: battle_object
    repeat: expr
    repeat-expr: 170
types:
  battle_object:
    seq:
      - id: type
        type: u1
        doc: Object type (index in OBDATA.DAT), 0xFF if unused or destroyed.
      - id: y
        type: u1
        doc: Position on the ground (Z is not updated while carried).
      - id: x
        type: u1
      - id: z
        type: u1
      - id: owner
        type: u1
        doc: |
          Index of the unit carrying the object (UNITPOS.DAT), 0xFF if it is
          on the ground.
      - id: slot
        type: u1
        doc: |
          Inventory slot: 0 = right hand, 1 = left hand, 2-3 = right leg,
          4-5 = left leg, 6-7 = right shoulder, 8-9 = left shoulder,
          10-18 = backpack, 19-24 = belt.
      - id: container
        type: u2
        doc: |
          Index of the object this object is loaded into (e.g. a clip in a
          rifle), 0xFFFF if none.
      - id: value1
        type: u1
        doc: |
          Ammunition, painkillers, primed turn of grenades or UNITPOS.DAT
          index of an unconscious unit (depending on the object type).
      - id: value2
        type: u1
        doc: Stimulants, last thrower of grenades or unit type of corpses.
      - id: value3
        type: u1
        doc: Heals of medikits or rank of corpses.
      - id: floor
        type: s1
        doc: Copy of the MCD floor offset of the tile.
      - id: flags
        type: u1
        doc: |
          Flags (e.g. 8 = unconscious unit). TACTICAL.EXE tests and sets bits
          0-3 of this byte. Only when creating an object does it access bytes
          12-15 as a 32-bit bitfield, changing only bits 0-2.
      - id: unknown
        size: 3
        doc: |
          Unknown, upper bytes of the 32-bit flags bitfield (see flags). They
          are only preserved by read-modify-write accesses and never set
          explicitly. Byte 13 is always 0, bytes 14-15 are usually 0 but some
          entries contain stale values like 0x0004, 0x0008 or 0x0288, 0x028c.
