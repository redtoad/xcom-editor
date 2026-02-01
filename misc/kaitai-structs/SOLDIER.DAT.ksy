meta:
  id: soldier_dat
  title: Format of SOLDIER.DAT
  file-extension: DAT
  endian: le
doc: |
  SOLDIER.DAT has a fixed structure of 250 entries of 68 bytes each. (Only a 
  maximum of 250 soldiers can be had.) Thus it is always 17,000 bytes long 
  (250x68).
seq:
  - id: soldier
    type: soldier_data
    repeat: eos
enums:
  ranks:
    0: rookie
    1: squaddie
    2: sergeant
    3: captain
    4: colonel
    5: commander
  bools:
    0: false
    1: true
  sex:
    0: male
    1: female
  appearance:
    0: blonde
    1: brown_hair
    2: oriental
    3: african
  armor:
    0: none
    1: personal_armor
    2: power_suit
    3: flying_suit
types:
  soldier_data:
    seq:
      - id: rank
        type: u2
        enum: ranks
        doc: |
          Integer set to Rank or FFFF if soldier is dead (or slot not yet used). 
          Ranks are:          

            0 Rookie
            1 Squaddie
            2 Sergeant (SGT)
            3 Captain (CPT)
            4 Colonel (COL)
            5 Commander (CDR)
            FF Dead
      - id: base
        type: u2
      - id: craft
        type: u2
      - id: craft_before
        type: u2
      - id: missions
        type: u2
      - id: kills
        type: u2
      - id: recovery_days
        type: u2
      - id: soldier_value
        type: u2
      - id: name
        type: strz
        size: 25
        encoding: ascii
      - id: destination_base
        type: u1
      - id: initial_time_units
        type: u1
      - id: initial_health
        type: u1
      - id: initial_energy
        type: u1
      - id: initial_reactions
        type: u1
      - id: initial_strength
        type: u1
      - id: initial_firing_accuracy
        type: u1
      - id: initial_throwing_accuracy
        type: u1
      - id: initial_melee_accuracy
        type: u1
      - id: initial_psionic_strength
        type: u1
      - id: initial_psionic_skill
        type: u1
      - id: initial_bravery
        type: u1
      - id: time_unit_improvement
        type: u1
      - id: health_improvement
        type: u1
      - id: energy_improvement
        type: u1
      - id: reactions_improvement
        type: u1
      - id: strength_improvement
        type: u1
      - id: firing_accuracy_improvement
        type: u1
      - id: throwing_accuracy_improvement
        type: u1
      - id: melee_accuracy_improvement
        type: u1
      - id: bravery_improvement
        type: u1
      - id: armor
        type: u1
        enum: armor
      - id: most_recent_psi_lab_training
        type: u1
      - id: in_psy_lab_training
        type: b8
        enum: bools
      - id: promotion
        type: b8
        enum: bools
      - id: sex
        type: u1
        enum: sex
      - id: appearance
        type: u1
        enum: appearance
  