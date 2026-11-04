meta:
  id: inter_dat
  title: Format of INTER.DAT
  file-extension: DAT
  endian: le
doc: |
  INTER.DAT details the specifics of any interceptions currently occurring on
  the Geoscape. It has a fixed structure of 4 entries (interception windows)
  of 142 bytes each and is thus always 568 bytes long.

  Apart from the first byte, the Go implementation keeps each record as raw
  bytes. The (unverified) meaning of the remaining offsets is listed in the
  doc of `raw_data`.
doc-ref: https://www.ufopaedia.org/index.php/INTER.DAT
seq:
  - id: interceptions
    type: interception_data
    size: 142
    repeat: expr
    repeat-expr: 4
types:
  interception_data:
    seq:
      - id: active_window
        type: s1
        enum: window_state
        doc: |
          0x00: Active/inactive window. 0 = active window (aircraft exists and
          is in intercept range), 1 = no active window (aircraft does not exist
          or does exist but is not in intercept range now). Documented as a
          word, but only this first byte is decoded.
      - id: raw_data
        size: 141
        doc: |
          0x01-0x8D: Not decoded (kept as raw bytes). Documented offsets (all
          relative to the start of the record, words unless noted):

          Geoscape window/icon related:
          0x02 aircraft type for icon's picture;
          0x04 byte A (craft number) from LOC.DAT for icon of minimized window;
          0x06 craft LOC.DAT entry number;
          0x08 CRAFT.DAT entry number;
          0x0A UFO LOC.DAT entry number.

          Active window:
          0x0C UFO size (5 - entry from UFO data: vs=0, s=1, m=2, l=3, vl=4);
          0x0E window minimized (=1) or active (=0);
          0x12 y position of window;
          0x14 x position of window;
          0x16 status message (starting at offset 0x397 in ENGLISH.DAT);
          0x18 count down before the status message is cleared;
          0x1A attack mode (0 = standoff, 1 = cautious, 2 = standard,
          3 = aggressive);
          0x1C current distance to target (km x8);
          0x1E requested distance to target (km x8), e.g. 560 = 70km;
          0x20 approach speed (change in position for projectiles);
          0x22 current air speed;
          0x24-0x41 15 positions of projectiles from the left weapon;
          0x42-0x5F 15 positions of projectiles from the right weapon;
          0x60 left weapon type; 0x62 right weapon type;
          0x64 left weapon reload timer; 0x66 right weapon reload timer;
          0x68 explosion/hit radius value (and window close delay timer);
          0x6A timer preventing firing, set when the UFO is downed/destroyed;
          0x6C UFO return fire intense level;
          0x6E timer preventing firing, set when an X-COM craft is destroyed;
          0x70 UFO side view preview opened;
          0x72 X-COM craft type; 0x74 UFO craft type;
          0x76 (dword) pointer to CRAFT.DAT offset of the attacking craft;
          0x7A (dword) pointer to attacking craft stats offset;
          0x7E (dword) pointer to attacking craft left weapon stats offset;
          0x82 (dword) pointer to attacking craft right weapon stats offset;
          0x86 (dword) pointer to CRAFT.DAT offset of the attacked UFO;
          0x8A (dword) pointer to UFO stats offset.
enums:
  window_state:
    0: active
    1: inactive
