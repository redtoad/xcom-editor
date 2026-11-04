meta:
  id: obdata_dat
  title: Format of OBDATA.DAT
  file-extension: DAT
  endian: le
doc: |
  OBDATA.DAT is a game data file found in the game's GEODATA/ directory (it
  is not part of a savegame). It contains general information on objects
  found in the game, such as their accuracy and weight. It also contains
  template information regarding clip size.

  The file consists of 80 records of 54 bytes each (4,320 bytes): 33 corpse
  entries, 41 objects and 6 obsolete/unused entries.

  `internal/geodata` has no Go struct for this file yet (only a package
  comment), so this layout follows UFOpaedia.
doc-ref: https://www.ufopaedia.org/index.php/OBDATA.DAT
seq:
  - id: objects
    type: object
    size: 54
    repeat: eos
enums:
  damage_type:
    0: armor_piercing
    1: incendiary
    2: high_explosive
    3: laser
    4: plasma
    5: stun
    6: melee
    7: acid
    255: none
  grenade_type:
    0: none
    1: explosive
    2: proximity
    3: smoke
    4: gas_4
    5: gas_5
    6: gas_6
    7: gas_7
types:
  object:
    seq:
      - id: name
        type: strz
        size: 20
        encoding: ASCII
        doc: Internal name of the object.
      - id: bigob
        type: u1
        doc: Image in BIGOBS.PCK (shown in the inventory).
      - id: floorob
        type: u1
        doc: Image in FLOOROB.PCK (shown when lying on the ground).
      - id: power
        type: u1
        doc: |
          Average weapon strength (for weapons using clips the value is
          stored with the clip).
      - id: stack_priority
        type: u1
        doc: Display priority when several items are stacked on one tile.
      - id: width
        type: u1
        doc: Width in the inventory.
      - id: height
        type: u1
        doc: Height in the inventory.
      - id: ammo
        type: u1
        repeat: expr
        repeat-expr: 3
        doc: Indexes of usable ammo types in OBDATA.DAT (255 = none).
      - id: handob
        type: u1
        doc: Image in HANDOB.PCK (shown when held).
      - id: firing_graphics
        type: u1
        doc: Graphics type for firing (appears to be unused).
      - id: damage_type
        type: u1
        enum: damage_type
        doc: Damage type of the ammo (255 if not ammo).
      - id: accuracy_auto
        type: u1
      - id: accuracy_snap
        type: u1
      - id: accuracy_aimed
        type: u1
      - id: tu_auto
        type: u1
        doc: Time units for an auto shot (percentage).
      - id: tu_snap
        type: u1
        doc: Time units for a snap shot (percentage).
      - id: tu_aimed
        type: u1
        doc: Time units for an aimed shot (percentage).
      - id: clip_size
        type: u1
        doc: |
          Multi-purpose: clip size for ammo, melee strength for weapons, turn
          of detonation for armed grenades.
      - id: melee_damage_type
        type: u1
        enum: damage_type
        doc: |
          Damage type of melee attacks (6 and 7 are TFTD vibroblade/electric
          and apparently unused in UFO).
      - id: melee_accuracy
        type: u1
      - id: melee_tu
        type: u1
        doc: Time units for a melee attack (percentage).
      - id: armor
        type: u1
        doc: Explosive strength needed to destroy the object ("constitution").
      - id: weight
        type: u1
        doc: Weight (used for throwing and encumbrance).
      - id: grenade_type
        type: u1
        enum: grenade_type
      - id: ai_priority
        type: u1
        doc: Priority for aliens picking up the object (only used when above 5).
      - id: is_ranged_weapon
        type: u1
      - id: is_melee_weapon
        type: u1
      - id: is_non_waypoint_weapon
        type: u1
        doc: Ranged weapon that does not use waypoints.
      - id: is_ammo
        type: u1
      - id: is_two_handed
        type: u1
      - id: is_waypoint_weapon
        type: u1
        doc: Guided missile launcher (uses waypoints).
      - id: is_alien_artefact
        type: u1
      - id: score
        type: u1
        doc: Points scored when recovered.
