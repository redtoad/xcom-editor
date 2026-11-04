meta:
  id: craft_dat
  title: Format of CRAFT.DAT
  file-extension: DAT
  endian: le
doc: |
  This file contains information specific to the crafts in the game. Both X-COM
  and alien crafts (including hidden ones) are in this file, though the values
  for the alien craft are still shrouded in mystery.

  The file has a fixed structure of 50 entries of 104 bytes each and is thus
  always 5,200 bytes long.
doc-ref: https://www.ufopaedia.org/index.php/CRAFT.DAT
seq:
  - id: crafts
    type: craft
    size: 104
    repeat: expr
    repeat-expr: 50
types:
  craft:
    seq:
      - id: type
        type: s1
        enum: craft_type
        doc: |
          0x00: Craft type. Possible values are the same as GEODATA.DAT. Unused
          entries have the value 0xFF (-1).
      - id: left_weapon
        type: s1
        enum: weapon_type
        doc: |
          0x01: Weapon placed in the left slot. (The Lightning does not have a
          center weapon type, only a left.)
      - id: left_ammo
        type: s2
        doc: 0x02-0x03
      - id: flight_mode
        type: s1
        enum: flight_mode
        doc: 0x04
      - id: right_weapon
        type: s1
        enum: weapon_type
        doc: 0x05
      - id: right_ammo
        type: s2
        doc: 0x06-0x07
      - id: unused
        type: s2
        doc: 0x08-0x09. Unused.
      - id: damage
        type: s2
        doc: |
          0x0A-0x0B: Damage, that is the amount it currently has taken. This
          value divided by the crafts damage capacity gives the percentage
          shown in-game.
      - id: altitude
        type: s2
        enum: altitude
        doc: |
          0x0C-0x0D: Altitude of craft. Is an index within ENGLISH.DAT for the
          string. (If craft is airborne and you change it to GROUND the altitude
          will remain the same. Speed must be edited to 0 for the change to
          hold.)
      - id: speed
        type: s2
        doc: 0x0E-0x0F
      - id: destination
        type: s2
        doc: |
          0x10-0x11: Index into LOC.DAT referencing the destination - for
          example, waypoints for X-COM craft, or X-COM bases for alien craft.
      - id: interception_ref
        type: s2
        doc: |
          0x12-0x13: Index into INTER.DAT when the ship is in interception mode.
      - id: next_ufo_waypoint_lon
        type: s2
        doc: 0x14-0x15. Next UFO waypoint coordinate X (longitude).
      - id: next_ufo_waypoint_lat
        type: s2
        doc: 0x16-0x17. Next UFO waypoint coordinate Y (latitude).
      - id: fuel_amount_remaining
        type: s2
        doc: |
          0x18-0x19: Fuel, amount remaining. This value divided by the crafts
          total fuel capacity gives the percentage shown in-game.
      - id: base_reference
        type: s2
        doc: 0x1A-0x1B. Base reference as an index to LOC.DAT.
      - id: mission_type
        type: s2
        enum: mission_type
        doc: |
          0x1C-0x1D: Mission type craft is on. Is an index within ENGLISH.DAT
          for the string (558 + this value).
      - id: mission_zone
        type: s2
        enum: mission_zone
        doc: |
          0x1E-0x1F: Zone where mission is being carried out. Is an index
          within ENGLISH.DAT for string (543 + this value).
      - id: ufo_trajectory_segment
        type: s2
        doc: 0x20-0x21. UFO trajectory segment (ranges from 0-7).
      - id: ufo_trajectory_type
        type: s2
        doc: 0x22-0x23. UFO trajectory type (ranges from 0-9).
      - id: alien_race
        type: s2
        enum: alien_race
        doc: |
          0x24-0x25: Alien race found on craft. Is index within ENGLISH.DAT for
          the string (466 + this value).
      - id: ufo_attack_timer
        type: s2
        doc: 0x26-0x27
      - id: ufo_escape_maneuver_timer
        type: s2
        doc: 0x28-0x29
      - id: status
        type: s2
        enum: craft_status
        doc: |
          0x2A-0x2B: Craft status. Is an index within ENGLISH.DAT for the string
          (268 + this value).
      - id: items
        type: u1
        repeat: expr
        repeat-expr: 56
        doc: |
          0x2C-0x63: Quantities of the items on board of the craft (the Go code
          keeps bytes 0x2C-0x67 as one raw 60 byte array `Cargo`).

          44 Tank/Cannon, 45 Tank/Rocket Launcher, 46 Tank/Laser Cannon,
          47 Hover Tank/Plasma, 48 Hover Tank/Launcher.
          Offsets 49-98 refer to offsets 0-49 in OBDATA.DAT:
          49 Pistol, 50 Pistol Clip, 51 Rifle, 52 Rifle Clip,
          53 Heavy Cannon, 54 Cannon AP-Ammo, 55 Cannon HE-Ammo,
          56 Cannon I-Ammo, 57 Auto-Cannon, 58 Auto-Cannon AP-Ammo,
          59 Auto-Cannon HE-Ammo, 60 Auto-Cannon I-Ammo, 61 Rocket Launcher,
          62 Small Rocket, 63 Large Rocket, 64 Incendiary Rocket,
          65 Laser Pistol, 66 Laser Gun, 67 Heavy Laser, 68 Grenade,
          69 Smoke Grenade, 70 Proximity Grenade, 71 High Explosive,
          72 Motion Scanner, 73 Medi-Kit, 74 Psi-Amp, 75 Stun Rod, 76 Flare,
          77-79 empty, 80 Corpse, 81 Corpse & Armour, 82 Corpse & Power Suit,
          83 Heavy Plasma, 84 Heavy Plasma Clip, 85 Plasma Rifle,
          86 Plasma Rifle Clip, 87 Plasma Pistol, 88 Plasma Pistol Clip,
          89 Blaster Launcher, 90 Blaster Bomb, 91 Small Launcher,
          92 Stun Missile, 93 Alien Grenade, 94 Elerium-115, 95 Mind Probe,
          96 undefined, 97-98 empty, 99 no known use.
      - id: flags
        type: u4
        doc: |
          0x64-0x67: Bitfield (only the 1st byte is used), see the `is_*`
          instances.
    instances:
      is_in_hangar:
        value: (flags & 0x01) != 0
        doc: Bit 0. Craft is in hangar (0 = craft has flown away).
      wants_to_go_home:
        value: (flags & 0x02) != 0
        doc: |
          Bit 1. Want to go home (set after mission or when out of fuel),
          blocks craft from re-targeting.
      is_out_of_elerium:
        value: (flags & 0x04) != 0
        doc: Bit 2. [runtime] Out of elerium (0 = fueled).
      is_out_of_left_ammo:
        value: (flags & 0x08) != 0
        doc: Bit 3. [runtime] Out of left ammo (0 = left weapon rearmed).
      is_out_of_right_ammo:
        value: (flags & 0x10) != 0
        doc: Bit 4. [runtime] Out of right ammo (0 = right weapon rearmed).
      is_processed_in_interception:
        value: (flags & 0x20) != 0
        doc: |
          Bit 5. [runtime] Flag of processed UFO in interception window, to
          avoid multiple escape/attack timers decrementing in case when UFO is
          pursued in more than 1 window at once.
      is_hyperwaved:
        value: (flags & 0x40) != 0
        doc: Bit 6. Hyperwaved extra info.
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
    -1:
      id: entry_not_used
      doc: Stored as 0xFF.
  weapon_type:
    0: stingray
    1: avalanche
    2: cannon
    3: fusion_ball
    4: laser_cannon
    5: plasma_beam
    -1:
      id: no_weapon
      doc: Stored as 0xFF (255).
  flight_mode:
    0:
      id: no_destination
      doc: No destination set (at base)
    1:
      id: single_destination
      doc: Single destination
    2:
      id: multiple_destinations
      doc: Multiple destinations (UFO only)
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
    12:
      id: north_atlantic
      doc: unused
    13:
      id: south_atlantic
      doc: unused
    14:
      id: indian_ocean
      doc: unused
  alien_race:
    0:
      id: sectoid
      doc: Aquatoid in TFTD
    1:
      id: snakeman
      doc: Gillman in TFTD
    2:
      id: ethereal
      doc: Lobsterman in TFTD
    3:
      id: muton
      doc: Tasoth in TFTD
    4:
      id: floater
      doc: Mixed Crew (Type I) in TFTD
    5:
      id: final_mission_mix
      doc: Mixed Crew (Type II) in TFTD
  craft_status:
    0: ready
    1: out
    2: repairs
    3: refueling
    4: rearming
