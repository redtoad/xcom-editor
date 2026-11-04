meta:
  id: wglob_dat
  title: Format of WGLOB.DAT
  file-extension: DAT
  endian: le
doc: |
  WGLOB.DAT holds the map dimensions, the current turn, the mission score and
  the speed options of the Battlescape.

  The file is 44 bytes long. Bytes 32-43 (scroll_speed to scroll_mode) have
  the same layout as MISSDAT/OPTION.DAT, which TACTICAL.EXE reads at startup
  (0x2f430) into its option variables and writes back at the end of a
  mission (0x2b7a2).

  The map dimensions (size_x * size_y * size_z tiles) determine how much of
  MAP.DAT, SEEMAP.DAT, SOURCEMP.DAT and TERMP.DAT is used; `routes`
  determines the used part of ROUTES.DAT.
doc-ref: https://www.ufopaedia.org/index.php/WGLOB.DAT
seq:
  - id: turn_times_two
    type: u2
    doc: Current turn * 2.
  - id: turn
    type: u2
    doc: Current turn (the first turn is 1).
  - id: size_x
    type: u2
    doc: Map size along the X axis.
  - id: size_y
    type: u2
    doc: Map size along the Y axis.
  - id: size_z
    type: u2
    doc: Number of map levels (Z axis).
  - id: area
    type: u2
    doc: size_x * size_y.
  - id: unknown1
    type: u2
    doc: Unknown, appears to be size_z - 1.
  - id: routes
    type: u2
    doc: Number of nodes in ROUTES.DAT.
  - id: mcd_records
    type: u2
    doc: Number of MCD records.
  - id: aliens_killed
    type: u2
    doc: Aliens killed.
  - id: soldiers_lost
    type: u2
    doc: X-COM soldiers lost.
  - id: tanks_lost
    type: u2
    doc: X-COM tanks lost.
  - id: alien_kill_score
    type: s2
    doc: Score for killing aliens.
  - id: xcom_loss_score
    type: s2
    doc: Score for losing X-COM units (negative).
  - id: civilians_killed_by_aliens
    type: u2
    doc: Civilians killed by aliens.
  - id: civilians_killed_by_xcom
    type: u2
    doc: Civilians killed by X-COM.
  - id: scroll_speed
    type: u2
    doc: Scroll speed (6 - value).
  - id: xcom_speed
    type: u2
    doc: X-COM movement speed (6 - value).
  - id: unknown2
    type: u2
    doc: Unknown, appears to mirror xcom_speed.
  - id: alien_speed
    type: u2
    doc: Alien movement speed (6 - value).
  - id: fire_speed
    type: u2
    doc: Fire speed.
  - id: scroll_mode
    type: u2
    enum: scroll_mode
    doc: "Scrolling: 0 = trigger scroll, 1 = auto scroll."
instances:
  tiles:
    value: size_x * size_y * size_z
    doc: Number of tiles of the map.
enums:
  scroll_mode:
    0: trigger_scroll
    1: auto_scroll
