meta:
  id: acts_dat
  title: Format of ACTS.DAT
  file-extension: DAT
  endian: le
doc: |
  ACTS.DAT holds the priorities of the alien missions ("acts") for each zone
  (in the order of the mission_zone enum). A priority of 0 means the act will
  not be selected. When an act is selected in a month, its priority is set to
  0. If it was the last non-zero act of the zone, the zone's priority in
  ZONAL.DAT is set to 0 too.

  The file is 84 bytes long (12 zones * 7 bytes).
doc-ref: https://www.ufopaedia.org/index.php/ACTS.DAT
seq:
  - id: zones
    type: zone_acts
    size: 7
    repeat: expr
    repeat-expr: 12
    doc: Mission priorities of each zone (indexed by mission_zone).
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
  zone_acts:
    doc: The mission priorities for a single zone.
    seq:
      - id: research
        type: u1
        doc: Alien Research.
      - id: harvest
        type: u1
        doc: Alien Harvest.
      - id: abduction
        type: u1
        doc: Alien Abduction (not tried in some zones, e.g. Antarctica).
      - id: infiltration
        type: u1
        doc: Alien Infiltration.
      - id: base
        type: u1
        doc: Alien Base.
      - id: terror
        type: u1
        doc: Alien Terror. Normally unused, terror missions are generated separately.
      - id: retaliation
        type: u1
        doc: Alien Retaliation. Normally unused, retaliation missions are generated separately.
