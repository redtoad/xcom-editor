meta:
  id: missions_dat
  title: Format of MISSIONS.DAT
  file-extension: DAT
  endian: le
doc: |
  MISSIONS.DAT holds the ongoing alien missions, indexed by zone (in the
  order of the mission_zone enum) and mission type (in the order of the
  mission_type enum). Alien supply missions are not included.

  The file is 672 bytes long (12 zones * 7 mission types * 8 bytes).
doc-ref: https://www.ufopaedia.org/index.php/MISSIONS.DAT
seq:
  - id: zones
    type: zone_missions
    size: 56
    repeat: expr
    repeat-expr: 12
    doc: Missions of each zone (indexed by mission_zone).
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
  mission_type:
    0: alien_research
    1: alien_harvest
    2: alien_abduction
    3: alien_infiltration
    4: alien_base
    5: alien_terror
    6: alien_retaliation
  alien_race:
    0: sectoid
    1: snakeman
    2: ethereal
    3: muton
    4: floater
    5: final_mission_mix
types:
  zone_missions:
    seq:
      - id: missions
        type: alien_mission
        size: 8
        repeat: expr
        repeat-expr: 7
        doc: Missions of the zone (indexed by mission_type).
  alien_mission:
    doc: A single ongoing alien mission.
    seq:
      - id: ufo_counter
        type: u2
        doc: Number of UFOs sent so far, 0xFFFF if the mission is not active.
      - id: sub_counter
        type: u2
        doc: Secondary UFO counter, only used by some missions.
      - id: timer
        type: u2
        doc: Time until the next UFO is sent (counts down every 30 minutes).
      - id: race
        type: u2
        enum: alien_race
        doc: |
          Alien race (same values as the alien race in CRAFT.DAT). Undefined for
          inactive missions.
    instances:
      is_active:
        value: ufo_counter != 0xffff
