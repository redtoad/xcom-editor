# Test data

Real savegames of UFO: Enemy Unknown used for round-trip tests (load and save
must result in byte-identical files).

- `GAME_1`, `GAME_2`, `GAME_3`: complete Geoscape savegames (all 27 files).
- `GAME_4`, `GAME_5`: savegames saved during a mission to an alien base in
  turn 2 and 3 (27 Geoscape and 14 Battlescape files).
- `GAME_6`: savegame saved during a mission with a tank (a large unit in the
  first slots of `UNITPOS.DAT`) and three soldiers still aboard the Skyranger.
- `PARTIAL`: some files of another savegame (named "Test"), used by tests that
  check decoded values.
