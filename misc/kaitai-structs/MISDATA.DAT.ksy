meta:
  id: misdata_dat
  title: Format of MISDATA.DAT
  file-extension: DAT
  endian: le
doc: |
  MISDATA.DAT holds the mission data and the score of the mission.

  The file is 170 bytes long. MISSDAT/MISSION.DAT and MISSDAT/MISSION2.DAT
  use the same format.
doc-ref: https://www.ufopaedia.org/index.php/MISDATA.DAT
seq:
  - id: recovered
    type: s2
    repeat: expr
    repeat-expr: 10
    doc: |
      Recovered UFO components (only set if the mission was won): power
      source, navigation, construction, alien food, reproduction,
      entertainment, surgery, examination room, alloys and habitat.
  - id: mission_type
    type: s2
    doc: Mission type.
  - id: scores
    type: mission_score
    repeat: expr
    repeat-expr: 13
    doc: |
      Score of the mission, in this order: aliens killed, alien corpses, live
      aliens, artefacts, alien base control, civilians killed by aliens,
      civilians killed by X-COM, civilians saved, X-COM soldiers killed,
      soldiers retired, soldiers missing in action, tanks lost and crafts
      lost.
  - id: time_of_day
    type: s2
    doc: Time of day (index into a table in the executable).
  - id: ufo_location
    type: u1
    doc: LOC.DAT index of the UFO crash site.
  - id: craft
    type: u1
    doc: |
      LOC.DAT index of the X-COM craft (or of the X-COM base for base
      defence). Note: UFOpaedia describes it as a CRAFT.DAT index, but in the
      savegames it refers to LOC.DAT.
  - id: result
    type: u1
    enum: mission_result
    doc: |
      Result: 0 = aborted, 1 = won, 2 = lost, 3 = aborted without units in
      the exit area.
  - id: ufo_site
    type: u1
    doc: 1 for UFO sites, 0 otherwise.
  - id: base_modules
    type: s2
    repeat: expr
    repeat-expr: 36
    doc: "Base defence: number of targets left in each base module."
  - id: equipment
    type: u1
    repeat: expr
    repeat-expr: 16
    doc: Equipment flags (1 = enabled).
  - id: language
    type: s2
    enum: language
    doc: "Language: 0 = English, 1 = German, 2 = French."
types:
  mission_score:
    doc: The score for one category of a mission.
    seq:
      - id: count
        type: s2
      - id: score
        type: s2
enums:
  mission_result:
    0: aborted
    1: won
    2: lost
    3: aborted_without_units_in_exit_area
  language:
    0: english
    1: german
    2: french
