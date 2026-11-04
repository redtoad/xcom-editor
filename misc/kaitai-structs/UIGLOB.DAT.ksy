meta:
  id: uiglob_dat
  title: Format of UIGLOB.DAT
  file-extension: DAT
  endian: le
doc: |
  UIGLOB.DAT stores counters used for naming new crafts, UFOs, missions and
  waypoints, the current year and the monthly scores.

  The file is 48 bytes long.
doc-ref: https://www.ufopaedia.org/index.php/UIGLOB.DAT
seq:
  - id: next_ufo
    type: u2
    doc: Number of the next spotted UFO.
  - id: next_terror_site
    type: u2
    doc: Number of the next terror mission.
  - id: next_alien_base
    type: u2
    doc: Number of the next alien base.
  - id: minimized_interception_windows
    type: u2
    doc: Number of minimized interception windows.
  - id: months_elapsed
    type: u2
    doc: |
      Months elapsed since the start of the game. This controls the alien
      stages and which (early/mid/late) weapons the aliens carry.
  - id: next_waypoint
    type: u2
    doc: Number of the next waypoint.
  - id: next_skyranger
    type: u2
    doc: Number of the next Skyranger.
  - id: next_lightning
    type: u2
    doc: Number of the next Lightning.
  - id: next_avenger
    type: u2
    doc: Number of the next Avenger.
  - id: next_interceptor
    type: u2
    doc: Number of the next Interceptor.
  - id: next_firestorm
    type: u2
    doc: Number of the next Firestorm.
  - id: year
    type: u2
    doc: Current year (e.g. 1999).
  - id: monthly_scores
    type: s2
    repeat: expr
    repeat-expr: 12
    doc: |
      Score for each month since the start of the game (index 0 is the first
      month). Scores can become negative.
