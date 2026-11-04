meta:
  id: facil_dat
  title: Format of FACIL.DAT
  file-extension: DAT
  endian: le
doc: |
  FACIL.DAT holds the stats of the 17 base facility types, in the order of
  the facility enum (the same order as the facility layout in BASE.DAT,
  where the hangar is split into four quarters 16-19; FACIL.DAT has a single
  hangar entry 16).

  The file is 272 bytes long (17 facilities * 16 bytes). A new game copies the
  table from GEOSCAPE.EXE (offset 374240) and it never changes during the
  game.
doc-ref: https://www.ufopaedia.org/index.php/FACIL.DAT
seq:
  - id: facilities
    type: facility_stats
    size: 16
    repeat: expr
    repeat-expr: 17
    doc: Stats of each facility type (indexed by facility).
enums:
  facility:
    0: access_lift
    1: living_quarters
    2: laboratory
    3: workshop
    4: small_radar_system
    5: large_radar_system
    6: missile_defense
    7: general_stores
    8: alien_containment
    9: laser_defense
    10: plasma_defense
    11: fusion_ball_defense
    12: grav_shield
    13: mind_shield
    14: psionic_laboratory
    15: hyperwave_decoder
    16: hangar
types:
  facility_stats:
    doc: The stats of a single base facility type.
    seq:
      - id: name_index
        type: u2
        doc: Index of the facility name in ENGLISH.DAT.
      - id: build_cost
        type: s2
        doc: Build cost in $1000.
      - id: build_days
        type: u1
        doc: Build time in days.
      - id: upkeep
        type: u1
        doc: Monthly upkeep in $1000.
      - id: capacity
        type: u1
        doc: Capacity (e.g. personnel, storage or containment space).
      - id: unknown1
        type: u1
        doc: Unknown, always 0.
      - id: defense
        type: s2
        doc: Defense strength.
      - id: hit_ratio
        type: u1
        doc: Defense hit ratio in percent.
      - id: unknown2
        type: u1
        doc: Unknown, always 0.
      - id: flags
        type: u1
        doc: |
          Build flags: 0x01 = special facility, 0x02 = defense facility,
          0x04 = can be built now, 0x08 = requires research.
      - id: unknown3
        size: 3
        doc: Unknown, always 0.
    instances:
      is_special:
        value: (flags & 0x01) != 0
      is_defense:
        value: (flags & 0x02) != 0
      can_be_built:
        value: (flags & 0x04) != 0
      requires_research:
        value: (flags & 0x08) != 0
