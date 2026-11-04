meta:
  id: unitref_dat
  title: Format of UNITREF.DAT
  file-extension: DAT
  endian: le
doc: |
  UNITREF.DAT holds the stats of all units on the map. Each entry belongs to
  the entry with the same index in UNITPOS.DAT.

  The file is 9920 bytes long (80 units * 124 bytes).

  MISSDAT/UNIREF.DAT (written by GEOSCAPE.EXE before a mission, only the
  units used) and MISSDAT/UNIREF2.DAT (written by TACTICAL.EXE at the end of
  a mission) use the same record format.
doc-ref: https://www.ufopaedia.org/index.php/UNITREF.DAT
seq:
  - id: units
    type: unit
    repeat: expr
    repeat-expr: 80
types:
  armor_values:
    doc: Armor values (front, left, right, rear, under).
    seq:
      - id: front
        type: u1
      - id: left
        type: u1
      - id: right
        type: u1
      - id: rear
        type: u1
      - id: under
        type: u1
  unit:
    seq:
      - id: type
        type: u1
        doc: Unit type, 0xFF if the slot is unused.
      - id: paper_doll
        type: u1
        doc: Inventory paper doll image.
      - id: pck_pointer
        type: u4
        doc: Memory address of the unit's PCK sprites (reset on load).
      - id: tab_pointer
        type: u4
        doc: Memory address of the unit's TAB sprites (reset on load).
      - id: facing
        type: u1
        doc: Facing (0-7).
      - id: movement_type
        type: u1
        enum: movement_type
        doc: "Movement type: 0 = walk, 1 = slide, 2 = fly."
      - id: time_units
        type: u1
      - id: health
        type: u1
      - id: stun
        type: u1
      - id: energy
        type: u1
      - id: reactions
        type: u1
      - id: strength
        type: u1
      - id: armor
        type: armor_values
        doc: Current armor.
      - id: firing_accuracy
        type: u1
      - id: throwing_accuracy
        type: u1
      - id: max_time_units
        type: u1
      - id: max_health
        type: u1
      - id: max_energy
        type: u1
      - id: max_reactions
        type: u1
      - id: max_armor
        type: armor_values
        doc: Maximum armor.
      - id: health_recharge
        type: u1
        doc: Health recharge (unused).
      - id: energy_recharge
        type: s1
        doc: Energy recharge per turn.
      - id: armor_recharge
        type: u1
        doc: Armor recharge (unused).
      - id: psi_skill
        type: u1
        doc: Psionic skill.
      - id: corpse
        type: u1
        doc: Item created when the unit dies or falls unconscious (OBDATA.DAT).
      - id: value
        type: u1
        doc: Score value of the unit.
      - id: soldier
        type: u1
        doc: Index in SOLDIER.DAT, 0xFF if the unit is not a soldier.
      - id: animation_frame
        type: u1
        doc: |
          Animation frame counter. TACTICAL.EXE resets it for all units after
          loading MISSDAT/UNIREF.DAT.
      - id: rank
        type: u1
        doc: Rank, 0xFF if dead or unused.
      - id: deployment
        type: u1
        doc: "Deployment: 0 = normal, 1 = civilian spawn."
      - id: aggression
        type: u1
        doc: Alien aggression.
      - id: energy_loss_divisor
        type: u1
      - id: turret_image
        type: u1
        doc: Tank turret image or inventory layout.
      - id: floor
        type: s1
        doc: Copy of the MCD floor offset of the tile.
      - id: lof_template
        type: u1
        doc: Line of fire template.
      - id: standing_height
        type: u1
      - id: kneeling_height
        type: u1
      - id: floating_height
        type: u1
      - id: lof_radius
        type: u1
        doc: Line of fire template radius.
      - id: left_hand_image
        type: u1
        doc: Image of the item in the left hand, 0xFF if none.
      - id: right_hand_image
        type: u1
        doc: Image of the item in the right hand, 0xFF if none.
      - id: damage_modifiers
        type: u1
        doc: Damage modifier category.
      - id: melee_accuracy
        type: u1
      - id: psi_strength
        type: u1
      - id: morale
        type: u1
        doc: Current morale.
      - id: bravery
        type: u1
        doc: Bravery (110 - 10 * value).
      - id: panic
        type: u1
        enum: panic
        doc: "Panic: 0 = none, 1 = freeze, 2 = run, 3 = berserk."
      - id: motion_scanner_step
        type: u1
      - id: motion_scanner_visibility
        type: u1
      - id: fatal_wounds
        type: u1
        repeat: expr
        repeat-expr: 6
        doc: |
          Fatal wounds (head, torso, right arm, left arm, right leg,
          left leg).
      - id: route_start
        type: u1
        doc: Patrol start node (ROUTES.DAT).
      - id: route_end
        type: u1
        doc: Patrol end node (ROUTES.DAT).
      - id: lighting
        type: u1
        doc: Light emitted by the unit.
      - id: painkiller_morale
        type: u1
        doc: Morale restored by painkillers.
      - id: intelligence
        type: u1
        doc: Intelligence (2-8).
      - id: ai_mode
        type: u1
        doc: AI mode, 0xFF to recalculate.
      - id: light_direction
        type: u1
        doc: Light direction (unused).
      - id: missions
        type: s2
      - id: kills
        type: s2
      - id: experience
        type: u1
        repeat: expr
        repeat-expr: 6
        doc: |
          Experience (reactions, firing, melee, throwing, psionics,
          bravery).
      - id: name
        type: strz
        size: 25
        encoding: ASCII
        doc: Name (null terminated).
      - id: big_fire
        type: u1
        doc: 1 while firing.
      - id: finished
        type: u1
        doc: 1 if the AI is done with the unit this turn.
      - id: tank
        type: u1
        doc: 1 for tanks.
      - id: burning_turns
        type: u1
        doc: Turns the unit is burning.
      - id: gender
        type: u1
        enum: gender
        doc: "Gender: 0 = male, 1 = female."
      - id: appearance
        type: u1
        doc: Appearance (0-3).
      - id: turret_weapon
        type: u1
        doc: Tank turret weapon (0xFF if unused).
      - id: turret_ammo
        type: u1
        doc: Tank turret ammunition.
      - id: unknown
        type: u1
        doc: Unknown, always 0.
      - id: flags
        type: u4
        doc: |
          Flags: 1 = moving, 2 = can fly, 4 = flying, 8 = selected this turn,
          16 = selection disabled, 32 = right hand item selected,
          64 = kneeling, 128 = wears power or flying suit.
enums:
  movement_type:
    0: walk
    1: slide
    2: fly
  panic:
    0: none
    1: freeze
    2: run
    3: berserk
  gender:
    0: male
    1: female
