meta:
  id: craft_dat
  file-extension: CRAFT.DAT
  endian: le
doc: |
  This file contains information specific to the crafts in the game. Both X-COM 
  and Alien crafts (including hidden ones) are in this file, though the values 
  for the alien craft are still shrouded in mystery.
seq:
  - id: crafts
    type: craft
    size: 104
    repeat: eos
enums:
  craft_type:
    0: skyranger           
    1: lightning     
    2: avenger       
    3: interceptor   
    4: firestorm     
    5: small_scout
    6: medium_scout
    7: large_scout
    8: harvester
    9: abductor
    10: terror_ship
    11: battleship
    12: supply_ship
    255: not_used
  weapon_type:
    0: stingray
    1: avalanche
    2: cannon
    3: fusion_ball
    4: laser_cannon
    5: plasma_beam
    255: no_weapon
  flight_mode:
    0: no_destination
    1: single_destination
    2: multiple_destinations
  altitude:
    0: ground
    1: very_low
    2: low
    3: high
    4: very_high
  mission_type:
    0: alien_research
    1: alien_harvest
    2: alien_abduction
    3: alien_infiltration
    4: alien_base
    5: alien_terror
    6: alien_retaliation
    7: alien_supply
  mission_zone:
    0: morth_america
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
    12: north_atlantic
    13: south_atlantic
    14: indian_ocean
  alien_race:
    0: sectoid_aquatoid
    1: snakeman_gillman
    2: ethereal_lobsterman
    3: muton_tasoth
    4: floater__mixed_crew_type_i
    5: final_mission_mix__mixed_crew_type_ii
  craft_status:
    0: ready
    1: out
    2: repairs
    3: refueling
    4: rearming
  status:
    0: craft_is_in_hangar
    2: want_to_go_home
    4: out_of_elerium
    8: out_of_left_ammo
    16: out_of_right_ammo
    32: ufo_in_interception_window
    64: hyperwaved_extra_info
types:
  craft:
    seq: 
      - id: type
        type: u1
        enum: craft_type
      - id: left_weapon_type
        type: u1
        enum: weapon_type
      - id: left_ammo
        type: u2
      - id: flight_mode
        type: u1
        enum: flight_mode
      - id: right_weapon_type
        type: u1
        enum: weapon_type
      - id: right_ammo
        type: u2
      - id: unused_
        type: u2
      - id: damage
        type: u2
        doc: |
          Damage, that is the amount it currently has taken. This value divided 
          by the crafts damage capacity gives the percentage shown in-game.
      - id: altitude
        type: u2
        enum: altitude
      - id: speed
        type: u2
      - id: loc_index
        type: u2
        doc: |
          Index into LOC.DAT referencing the destination - for example, 
          waypoints for X-COM craft, or X-COM bases for alien craft.
      - id: inter_index
        type: u2
        doc: |
          Index into INTER.DAT when the ship is in interception mode.
      - id: next_waypoint_lon
        type: u2
        doc: |
          Next UFO waypoint coordinate X (longitude).
      - id: next_waypoint_lat
        type: u2
        doc: |
          Next UFO waypoint coordinate Y (latitude).
      - id: fuel
        type: u2
        doc: |
          Fuel, amount remaining. This value divided by the crafts total fuel 
          capacity gives the percentage shown in-game.
      - id: base
        type: u2
        doc: |
          Base reference as an index to LOC.DAT.
      - id: mission_type
        type: u2
        enum: mission_type
        doc: |
          Mission type craft is on. Is an index within ENGLISH.DAT for the 
          string (558 + this value).
      - id: mission_zone
        type: u2
        enum: mission_zone
        doc: |
          Zone where mission is being carried out. Is an index within 
          ENGLISH.DAT for string (543 + this value).
      - id: ufo_trajectory_segment 
        type: u2
        doc: UFO trajectory segment (ranges from 0-7).
      - id: ufo_trajectory_type
        type: u2
        doc: UFO trajectory type (ranges from 0-9).
      - id: alien_race
        type: u2
        enum: alien_race
        doc: |
          Alien Race found on craft. Is index within ENGLISH.DAT for the string 
          (466 + this value).
      - id: ufo_attack_timer
        type: u2
      - id: ufo_escape_maneuver_timer
        type: u2
      - id: craft_status
        type: u2
        enum: craft_status
      - id: cargo
        type: u1
        repeat: expr
        repeat-expr: 55
        doc: |
          The rest of the known values detail the items on board of the craft. 
          49-98 refer to offsets 0-49 in OBDATA.DAT.
      - id: status
        type: b6
        enum: status

        
      
        
