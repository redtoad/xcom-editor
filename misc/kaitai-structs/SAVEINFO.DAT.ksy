meta:
  id: saveinfo_dat
  title: Format of SAVEINFO.DAT
  file-extension: DAT
  endian: le
doc: |
  SAVEINFO.DAT stores the game time and save game title. This file is present
  in any save. It is 40 bytes long, no separate entries.

  The values seen in the Load/Save game screen are read from this file; the
  actual game time is stored in IGLOB.DAT.
doc-ref: https://www.ufopaedia.org/index.php/SAVEINFO.DAT
seq:
  - id: missdat_flag
    type: u2
    enum: bools
    doc: |
      0x00-0x01: Ignore this if the file is not in the missdat folder. If 0,
      then this is a savegame made on the beginning of a new battlescape game.
      If 1, then check DIRECT.DAT to see which save slot to load from.
  - id: name
    type: strz
    size: 26
    encoding: ASCII
    doc: |
      0x02-0x1B: This is a 26 byte null terminated string, which details the
      name of the save file. The name may be 25 characters long; the final byte
      is always of value 0. A 0 also marks the end of the save name, should it
      not use all 25 characters.
  - id: year
    type: s2
    doc: 0x1C-0x1D. The current year.
  - id: month
    type: s2
    enum: months
    doc: 0x1E-0x1F. The current month. Note that 0 is January.
  - id: day_of_month
    type: s2
    doc: 0x20-0x21. The current day of the month.
  - id: hour
    type: s2
    doc: 0x22-0x23. The current hour (24 hour time).
  - id: minute
    type: s2
    doc: 0x24-0x25. The current minute.
  - id: tactical_save
    type: u2
    enum: bools
    doc: 0x26-0x27. 0 for geoscape save, 1 for tactical save.
enums:
  months:
    0: jan
    1: feb
    2: mar
    3: apr
    4: may
    5: jun
    6: jul
    7: aug
    8: sep
    9: oct
    10: nov
    11: dec
  bools:
    0: false
    1: true
