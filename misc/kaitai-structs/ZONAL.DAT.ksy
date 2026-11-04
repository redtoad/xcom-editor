meta:
  id: zonal_dat
  title: Format of ZONAL.DAT
  file-extension: DAT
  endian: le
doc: |
  ZONAL.DAT is a dynamic list of the chances that the aliens will conduct a
  mission in a particular zone. There is one priority byte for each of the 12
  zones, in the order of the mission_zone enum (North America, Arctic,
  Antarctica, South America, Europe, North Africa, Southern Africa, Central
  Asia, South East Asia, Siberia, Australasia, Pacific).

  A zone with priority 0 will only get terror attacks. Once research, harvest,
  abduction, infiltration and base creation missions have all been performed
  in a zone (see ACTS.DAT), the value for that zone is set to 0.

  The file is 12 bytes long.
doc-ref: https://www.ufopaedia.org/index.php/ZONAL.DAT
seq:
  - id: priorities
    type: u1
    repeat: expr
    repeat-expr: 12
    doc: Mission priority of each zone (indexed by mission_zone).
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
