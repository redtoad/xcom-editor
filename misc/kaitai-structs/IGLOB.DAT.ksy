meta:
  id: iglob_dat
  title: Format of IGLOB.DAT
  file-extension: DAT
  endian: le
doc: |
  IGLOB.DAT stores the current date, time, and game state counters. The file
  is 60 bytes in the original DOS EU version, or 64 bytes in patched versions
  which also include a difficulty setting.

  Note that changing the days, minutes, hours or seconds do not change the
  values seen in the Load/Save game screen, but do change the values when you
  actually load the game. The values seen in the Load/Save screen are read
  from SAVEINFO.DAT.
doc-ref: https://www.ufopaedia.org/index.php/IGLOB.DAT
seq:
  - id: month
    type: s4
    doc: 0x00. Current month (1-12).
  - id: weekday
    type: s4
    doc: 0x04. Current day of week (0-6).
  - id: day
    type: s4
    doc: 0x08. Current day of month (1-31).
  - id: hour
    type: s4
    doc: 0x0C. Current hour (0-23).
  - id: minute
    type: s4
    doc: 0x10. Current minute (0-59).
  - id: second
    type: s4
    doc: 0x14. Current second (0-55, multiples of 5).
  - id: hidden_craft
    type: s4
    doc: 0x18. Number of craft in hidden status.
  - id: airborne_xcraft
    type: s4
    doc: 0x1C. Number of airborne X-COM craft.
  - id: refueling_craft
    type: s4
    doc: 0x20. Number of craft refueling.
  - id: damaged_craft
    type: s4
    doc: 0x24. Number of damaged craft.
  - id: repairing_craft
    type: s4
    doc: 0x28. Number of craft being repaired.
  - id: airborne_craft
    type: s4
    doc: 0x2C. Total number of airborne craft.
  - id: ufos_on_retaliation
    type: s4
    doc: 0x30. Number of UFOs on retaliation missions.
  - id: alien_bases
    type: s4
    doc: 0x34. Number of alien bases.
  - id: current_midi
    type: s4
    doc: 0x38. Current MIDI track number.
  - id: difficulty
    type: s4
    enum: difficulty
    if: not _io.eof
    doc: |
      0x3C: (optional) Difficulty setting. Only present in patched versions
      (64 byte file); absent in the 60 byte DOS EU version.
enums:
  difficulty:
    0: beginner
    1: experienced
    2: veteran
    3: genius
    4: superhuman
