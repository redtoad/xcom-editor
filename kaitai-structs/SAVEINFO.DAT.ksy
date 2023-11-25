meta:
  id: saveinfo_dat
  title: Format of SAVEINFO.DAT
  file-extension: DAT
  endian: le
doc: |
  SAVEINFO.DAT stores the game time and save game title. 
  This file is present in any save.
seq:
  - id: battlescape_game
    type: u2
    enum: bools
    doc: |
      Ignore this if the file is not in the missdat folder. If 0, then 
      this is a savegame made on the beginning of a new battlescape game. 
      If 1, then check DIRECT.DAT to see where which save slot to load 
      from.
  - id: name
    type: strz
    size: 26
    encoding: ascii
  - id: year
    type: u2
  - id: month
    type: u2
    enum: months
  - id: day
    type: u2
  - id: hour
    type: u2
  - id: minute
    type: u2
  - id: tactical_save
    type: u2
    enum: bools
    doc: |
      0 for geoscape save, 1 for tactical save.
enums:
  months:
    0: jan
    1: feb
    3: mar
    4: apr
    5: may
    6: jun
    7: jul
    8: aug
    9: sep
    10: oct
    11: nov
    12: dec
  bools:
    0: false
    1: true
    