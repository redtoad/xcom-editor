meta:
  id: astore_dat
  title: Format of ASTORE.DAT
  file-extension: DAT
  endian: le
doc: |
  ASTORE.DAT holds the live aliens kept in alien containment.

  The file is 600 bytes long (50 entries * 12 bytes).
doc-ref: https://www.ufopaedia.org/index.php/ASTORE.DAT
seq:
  - id: aliens
    type: stored_alien
    size: 12
    repeat: expr
    repeat-expr: 50
enums:
  alien_race:
    0: none
    1: sectoid
    2: snakeman
    3: ethereal
    4: muton
    5: floater
    6: celatid
    7: silacoid
    8: chryssalid
    9: reaper
    10: sectopod
    11: cyberdisc
  alien_rank:
    1: commander
    2: leader
    3: engineer
    4: medic
    5: navigator
    6: soldier
    7: terrorist
types:
  stored_alien:
    doc: A single live alien in alien containment.
    seq:
      - id: race
        type: u1
        enum: alien_race
        doc: Alien race. `none` if the entry is unused (e.g. after the alien was researched).
      - id: rank
        type: u1
        enum: alien_rank
      - id: base
        type: u1
        doc: Base the alien is kept in (index in LOC.DAT). 0xFF while in transfer.
      - id: unknown
        size: 8
        doc: |
          Unknown, usually 0. Neither GEOSCAPE.EXE nor TACTICAL.EXE (which only
          copies the file) reads or writes these bytes: a new game only clears
          bytes 0 and 11 of each entry (GEOSCAPE.EXE 0x21db0) and storing an
          alien only sets bytes 0, 1, 2 and 11. The bytes therefore keep
          uninitialized memory, e.g. 00 a4 3b 20 00 00 00 00 in entry 0.
      - id: transfer_state
        type: u1
        doc: "Transfer state: 1 = in transit, 2 = selected for transfer."
