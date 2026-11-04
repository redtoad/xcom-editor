meta:
  id: xbases_dat
  title: Format of XBASES.DAT
  file-extension: DAT
  endian: le
doc: |
  XBASES.DAT holds information about alien attacks on X-COM bases, one record
  for each zone (in the order of the mission_zone enum).

  The format is only partially documented. The file is 48 bytes long
  (12 zones * 4 bytes).
doc-ref: https://www.ufopaedia.org/index.php/XBASES.DAT
seq:
  - id: zones
    type: xbase_target
    size: 4
    repeat: expr
    repeat-expr: 12
    doc: Alien targeting information of each zone (indexed by mission_zone).
enums:
  mission_zone:
    0: north_america
    1: arctic
    2: antarctica
    3: south_america
    4: europe
    5: north_africa
    6: southern_africa
    7: central_asia
    8: south_east_asia
    9: siberia
    10: australasia
    11: pacific
types:
  xbase_target:
    doc: The alien targeting information for a single zone.
    seq:
      - id: detected
        type: u2
        doc: Non-zero if and only if an X-COM base has been detected.
      - id: target
        type: u2
        doc: LOC.DAT reference of the X-COM base to be targeted.
