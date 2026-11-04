meta:
  id: base_dat
  title: Format of BASE.DAT
  file-extension: DAT
  endian: le
doc: |
  BASE.DAT has all of the base layout and contents information, as well as
  base name info. It has a fixed structure of 8 entries of 292 bytes each and
  is thus always 2,336 bytes long.
doc-ref: https://www.ufopaedia.org/index.php/BASE.DAT
seq:
  - id: bases
    type: base_data
    size: 292
    repeat: expr
    repeat-expr: 8
types:
  base_data:
    seq:
      - id: name
        type: strz
        size: 16
        encoding: ASCII
        doc: |
          0x00-0x0F: Base name. 0x0F is presumably the null character if the
          name uses all 15 characters.
      - id: short_range
        type: s2
        doc: |
          0x10-0x11: Base's short range detection capability.

          Logical values for the detection capabilities:
          10 short, 0 long: This base has small radar(s) only.
          20 short, 20 long: This base has large radar(s) only.
          30 short, 20 long: This base has small and large radar(s).
          100 hyperwave: This base has a hyperwave decoder(s).

          The radar values can be set to 100 for perfect short range detection
          (presumably -- it definitely makes UFOs appear more often), but these
          reset to the correct values any time you complete a build in that
          base.
      - id: long_range
        type: s2
        doc: 0x12-0x13. Base's long range detection capability.
      - id: hyperwave
        type: s2
        doc: 0x14-0x15. Base's hyperwave detection capability.
      - id: grid
        type: u1
        enum: facility
        repeat: expr
        repeat-expr: 36
        doc: 0x16-0x39. Facilities in the 6x6 base grid (36 cells).
      - id: days_to_completion
        type: u1
        repeat: expr
        repeat-expr: 36
        doc: |
          0x3A-0x5D: Days until the facility in the matching grid cell is
          completed. Set up the same way as `grid`.
      - id: engineers
        type: s1
        doc: 0x5E. Number of engineers in the base.
      - id: scientists
        type: s1
        doc: 0x5F. Number of scientists in the base.
      - id: inventory
        type: s2
        repeat: expr
        repeat-expr: 96
        doc: |
          0x60-0x11F: Item quantities in base stores. The array index is the
          item, see enum `inventory_item` (indices 43-45 and 62-64 are unused).
      - id: active
        type: u4
        enum: base_state
        doc: |
          0x120-0x123: Active/inactive base, stored as a 4-byte integer.
          Inactive entries have a value of 1. Active entries have a value of 0.
          Creating a new base will overwrite the first inactive entry. If a
          base is dismantled, the only change to the record is this value so it
          is possible to restore a dismantled base (access lift removed) by
          restoring this value to 0.
enums:
  base_state:
    0: active
    1: inactive
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
    16:
      id: hangar_top_left
      doc: Hangar (top left)
    17:
      id: hangar_top_right
      doc: Hangar (top right)
    18:
      id: hangar_bottom_left
      doc: Hangar (bottom left)
    19:
      id: hangar_bottom_right
      doc: Hangar (bottom right)
    255: empty
  inventory_item:
    # Index into the `inventory` array of a base (not a stored value).
    0: stingray_launcher
    1: avalanche_launcher
    2: cannon
    3: fusion_ball_launcher
    4: laser_cannon
    5: plasma_beam
    6: stingray_missile
    7: avalanche_missile
    8: cannon_rounds
    9: fusion_balls
    10: tank_cannon
    11: tank_rocket_launcher
    12: tank_laser_cannon
    13: hovertank_plasma
    14: hovertank_launcher
    15: pistol
    16: pistol_clip
    17: rifle
    18: rifle_clip
    19: heavy_cannon
    20: hc_ap_ammo
    21: hc_he_ammo
    22: hc_in_ammo
    23: auto_cannon
    24: ac_ap_ammo
    25: ac_he_ammo
    26: ac_in_ammo
    27: rocket_launcher
    28: small_rocket
    29: large_rocket
    30: incendiary_rocket
    31: laser_pistol
    32: laser_rifle
    33: heavy_laser
    34: grenade
    35: smoke_grenade
    36: proximity_grenade
    37: high_explosive
    38: motion_scanner
    39: medi_kit
    40: psi_amp
    41: stun_rod
    42: electro_flare
    46: corpse
    47: corpse_armour
    48: corpse_powersuit
    49: heavy_plasma
    50: heavy_plasma_clip
    51: plasma_rifle
    52: plasma_rifle_clip
    53: plasma_pistol
    54: plasma_pistol_clip
    55: blaster_launcher
    56: blaster_bomb
    57: small_launcher
    58: stun_bomb
    59: alien_grenade
    60: elerium_115
    61: mind_probe
    65: sectoid_corpse
    66: snakeman_corpse
    67: ethereal_corpse
    68: muton_corpse
    69: floater_corpse
    70: celatid_corpse
    71: silacoid_corpse
    72: chryssalid_corpse
    73: reaper_corpse
    74: sectopod_corpse
    75: cyberdisc_corpse
    76: hovertank_corpse
    77: tank_corpse
    78: male_civilian_corpse
    79: female_civilian_corpse
    80: ufo_power_source
    81: ufo_navigation
    82: ufo_construction
    83: alien_food
    84: alien_reproduction
    85: alien_entertainment
    86: alien_surgery
    87: examination_room
    88: alien_alloys
    89: alien_habitat
    90: personal_armour
    91: power_suit
    92: flying_suit
    93: hwp_cannon_shell
    94: hwp_rockets
    95: hwp_fusion_bomb
