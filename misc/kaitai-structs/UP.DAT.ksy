meta:
  id: up_dat
  title: Format of UP.DAT
  file-extension: DAT
  endian: le
doc: |
  UP.DAT holds the UFOpaedia entries.

  The file is 1392 bytes long (116 entries * 12 bytes). A new game copies the
  table from GEOSCAPE.EXE (offset 386176); only research changes during the
  game, all unknown fields are constants.
doc-ref: https://www.ufopaedia.org/index.php/UP.DAT
seq:
  - id: entries
    type: up_entry
    size: 12
    repeat: expr
    repeat-expr: 116
enums:
  research_status:
    0: not_researched
    2: researched
    3: available_from_start
types:
  up_entry:
    doc: A single UFOpaedia entry.
    seq:
      - id: category
        type: u1
        doc: Category (0-8).
      - id: unknown1
        type: u1
        doc: Unknown, always 0.
      - id: sub_index
        type: u1
        doc: |
          Index within the category. For categories 2 and 3 this is the index
          in OBDATA.DAT, armours use 100-102.
      - id: unknown2
        size: 2
        doc: Unknown, always 0.
      - id: image
        type: u1
        doc: Image number (GEOGRAPH/UP0##.SPK).
      - id: string_index
        type: u2
        doc: Index of the title in ENGLISH.DAT (1-based).
      - id: research
        type: u1
        enum: research_status
        doc: "Research status: 0 = not researched, 2 = researched, 3 = available from the start."
      - id: unknown3
        size: 3
        doc: Unknown, always 0.
