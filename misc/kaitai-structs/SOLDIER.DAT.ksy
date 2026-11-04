meta:
  id: soldier_dat
  title: Format of SOLDIER.DAT
  file-extension: DAT
  endian: le
doc: |
  SOLDIER.DAT contains all information about the soldiers. It has a fixed
  structure of 250 entries of 68 bytes each. (Only a maximum of 250 soldiers
  can be had.) Thus it is always 17,000 bytes long (250x68).
doc-ref: https://www.ufopaedia.org/index.php/SOLDIER.DAT
seq:
  - id: soldiers
    type: soldier_data
    size: 68
    repeat: expr
    repeat-expr: 250
types:
  soldier_data:
    seq:
      # Basic information
      - id: rank
        type: s2
        enum: rank
        doc: |
          0x00-0x01 (0-5, 0xFFFF): Rank, or 0xFFFF (-1) if soldier is dead (or
          slot not yet used). Values are actually a pointer within ENGLISH.DAT.
      - id: base
        type: s2
        doc: |
          0x02-0x03 (various, 0xFF): Base the soldier is at, using reference
          values from LOC.DAT. First base is always 0 (first entry in LOC.DAT).
          Value of 0xFF means the soldier is being transferred.
      - id: craft
        type: s2
        doc: |
          0x04-0x05 (various, 0xFFFF): Craft soldier is assigned to, using
          values from LOC.DAT. 0xFFFF (-1) means not assigned to a craft.
      - id: craft_before
        type: s2
        doc: |
          0x06-0x07 (various, 0xFFFF): Craft soldier was assigned to before
          being wounded and spending time in the infirmary.
      - id: missions
        type: s2
        doc: |
          0x08-0x09 (0+): Missions. Signed integer. (The soldier display shows
          a negative value if highest/16th bit is set!)
      - id: kills
        type: s2
        doc: 0x0A-0x0B (0+). Kills. Signed integer.
      - id: recovery_days
        type: s2
        doc: 0x0C-0x0D (0+). Wound recovery days. Signed integer.
      - id: soldier_value
        type: s2
        doc: |
          0x0E-0x0F (20+): Soldier value a.k.a. victory point loss if dies on a
          mission. Signed integer. Is equal to 20 + Missions + Rank Bonus
          (Rookie 0, Squaddie +1, SGT +1, CPT +3, COL +6, CDR +10).

          Note: Hex-editing Missions or Rank and performing a combat (while not
          touching this field) does not cause this field to reflect the edited
          values. This field simply adds 1 for each mission and a delta when
          promoted (e.g. COL to CDR +4).
      - id: name
        type: strz
        size: 25
        encoding: ASCII
        doc: |
          0x10-0x28 (text): Name, 25 characters long. Only approx. 21
          characters can be entered when editing the field within the game,
          but longer names can be hex edited and will be seen in the game. The
          end of the current name is a null byte; garbage can be present after
          that (ends of longer previous names, etc.).
      - id: destination_base
        type: s1
        doc: |
          0x29 (various): Always 0 except for existing soldiers being
          transferred, in which case it equals the LOC.DAT value for
          destination base. New recruits on their way to a base will also have
          this set to 0, so if byte 3 is set to 0xFF and this byte is 0, you
          can't tell (from these two bytes alone) if it's an existing soldier
          going to your first base (LOC.DAT=0), or a new recruit going to any
          base. But the game knows where they're going, via TRANSFER.DAT.
      # Initial recruit stats. All stat bytes (42-51, 53-60, 63) are unsigned.
      # Values above 127 occur in hacked or edited savegames.
      - id: initial_time_units
        type: u1
        doc: 0x2A (50-60). Initial Time Units.
      - id: initial_health
        type: u1
        doc: 0x2B (25-40). Initial Health.
      - id: initial_energy
        type: u1
        doc: 0x2C (40-70). Initial Energy a.k.a. Stamina.
      - id: initial_reactions
        type: u1
        doc: 0x2D (30-60). Initial Reactions.
      - id: initial_strength
        type: u1
        doc: 0x2E (20-40). Initial Strength.
      - id: initial_firing_accuracy
        type: u1
        doc: 0x2F (40-70). Initial Firing Accuracy.
      - id: initial_throwing_accuracy
        type: u1
        doc: 0x30 (50-80). Initial Throwing Accuracy.
      - id: initial_melee_accuracy
        type: u1
        doc: 0x31 (20-40). Initial Melee Accuracy a.k.a. Close Combat Accuracy.
      - id: initial_psionic_strength
        type: u1
        doc: 0x32 (0-100). Psionic Strength. Never changes.
      - id: initial_psionic_skill
        type: u1
        doc: |
          0x33 (0-105+): Current Psionic Skill. 0 = not yet trained / psi stats
          not visible. Unlike other stats, Psi Skill always starts at 0, so the
          current skill is stored in what would otherwise be the initial skill
          byte.
      - id: initial_bravery
        type: s1
        doc: |
          0x34 (5-10): Initial Bravery. Computed as 110 - (10 * this byte).
          E.g. if this byte is 9, Bravery is 20. So the lower this byte is, the
          better your Bravery is. Initial Bravery can be 10-60, and it only
          uses increments of 10.
      # Stat improvement. Add to initial values to get current total.
      - id: time_unit_improvement
        type: u1
        doc: 0x35 (0-31). Time Unit Improvement. (Max total TUs possible is 81.)
      - id: health_improvement
        type: u1
        doc: 0x36 (0-36). Health Improvement. (Max total 61)
      - id: energy_improvement
        type: u1
        doc: 0x37 (0-61). Stamina Improvement. (Max total 101)
      - id: reactions_improvement
        type: u1
        doc: 0x38 (0-75). Reaction Improvement. (Max total 105)
      - id: strength_improvement
        type: u1
        doc: |
          0x39 (0-51): Strength Improvement. (Max total 71) Note that very high
          (hacked) values can cause "unable to throw here".
      - id: firing_accuracy_improvement
        type: u1
        doc: 0x3A (0-85). Firing Accuracy Improvement. (Max total 125)
      - id: throwing_accuracy_improvement
        type: u1
        doc: 0x3B (0-75). Throwing Accuracy Improvement. (Max total 125)
      - id: melee_accuracy_improvement
        type: u1
        doc: 0x3C (0-105). Melee/Close Combat Accuracy Improvement. (Max total 125)
      - id: bravery_improvement
        type: s1
        doc: 0x3D (0-9). Bravery Improvement (multiply by 10). (Max total 9; 100 Bravery)
      # Odds and ends
      - id: armor
        type: s1
        enum: armor
        doc: |
          0x3E (0-3): Armor. Can be hacked to alien skins (4-11); higher values
          are civilians and parts of big units. The CE version crashes when
          starting a mission with such a value.
      - id: most_recent_psi_lab_training
        type: u1
        doc: |
          0x3F (0-24): Most-recent month's Psi Lab training increase. Largest
          values for newly trained (16-24).
      - id: in_psi_lab_training
        type: u1
        enum: bools
        doc: 0x40 (0,1). In Psi Lab training.
      - id: promotion
        type: u1
        enum: bools
        doc: |
          0x41 (0,1): Promotion flag (as of most recent combat). Reset as of
          next combat.
      - id: sex
        type: s1
        enum: gender
        doc: |
          0x42 (0,1): Editing values greater than 1 does not change the
          appearance in the soldier equip screens. However, the soldier image
          will always take the outward appearance of a female during a mission.
      - id: appearance
        type: s1
        enum: appearance
        doc: |
          0x43 (0-3): Appearance (for loadout, no armor). Refers to the
          UFOGRAPH/MAN_*.SPK images. Values greater than 3 result in an error
          because the image file cannot be found.
enums:
  rank:
    -1:
      id: dead_or_unused
      doc: Dead or Unused (stored as 0xFFFF)
    0: rookie
    1: squaddie
    2:
      id: sergeant
      doc: SGT
    3:
      id: captain
      doc: CPT
    4:
      id: colonel
      doc: COL
    5:
      id: commander
      doc: CDR
  bools:
    0: false
    1: true
  gender:
    0: male
    1: female
  appearance:
    0: blonde
    1: brown_hair
    2: oriental
    3: african
  armor:
    0: no_armor
    1: personal_armour
    2: power_suit
    3: flying_suit
    4:
      id: sectoid_armour
      doc: hacked value
    5:
      id: snakeman_armour
      doc: hacked value
    6:
      id: ethereal_armour
      doc: hacked value
    7:
      id: muton_armour
      doc: hacked value
    8:
      id: floater_armour
      doc: hacked value
    9:
      id: celatid_armour
      doc: hacked value
    10:
      id: silacoid_armour
      doc: hacked value
    11:
      id: chryssalid_armour
      doc: hacked value
