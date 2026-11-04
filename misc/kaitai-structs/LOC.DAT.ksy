meta:
  id: loc_dat
  title: Format of LOC.DAT
  file-extension: DAT
  endian: le
doc: |
  LOC.DAT provides location data for bases, crafts and other objects on the
  Geoscape. It has a row width of 20 bytes. There are a total of 50 records
  (not all of them necessarily used) for a fixed file size of 1,000 bytes.
doc-ref: https://www.ufopaedia.org/index.php/LOC.DAT
seq:
  - id: objects
    type: location
    size: 20
    repeat: expr
    repeat-expr: 50
types:
  location:
    seq:
      - id: type
        type: u1
        enum: location_type
        doc: 0x00. Object type.
      - id: table_reference
        type: u1
        doc: |
          0x01: Object table reference (0x00 to 0xFF). This just shows how many
          there are of this type on the geoscape. If the object is either a UFO
          or X-COM craft, then this is the index into CRAFT.DAT. If the object
          is an X-COM base, then this is the index into BASE.DAT and if it is
          an alien base this byte contains the race.
      - id: x
        type: s2
        doc: |
          0x02-0x03: Horizontal coordinate or longitude. Value range: 0 - 2880
          (1/8 degree).
      - id: y
        type: s2
        doc: |
          0x04-0x05: Vertical coordinate or latitude. Value range: -720 - 720
          (1/8 degree).
      - id: timer
        type: s2
        doc: |
          0x06-0x07: For crash site or terror site - countdown timer (in
          hours). For moving objects - how many game ticks (5s) have to pass
          until craft moved to next globe coordinate (cell_size div speed).
          Note that ground UFOs are treated as moving objects, except for
          speed = 0.
      - id: fraction
        type: s2
        doc: |
          0x08-0x09: Fractional part of how much is left to the next globe
          coordinate (cell_size mod speed), used only for moving objects.
      - id: count_suffix
        type: s2
        doc: |
          0x0A-0x0B: Count suffix of the item, eg: Skyranger-1 or Crash
          Site-47. It appears to have no meaning for X-COM bases, but for other
          types where it is set, 0x0B is the high byte for when you go over 255
          UFOs or crafts, etc.
      - id: unused1
        type: s2
        doc: 0x0C-0x0D. Unused.
      - id: transfer_mode
        type: s1
        doc: 0x0E. Craft transfer mode.
      - id: unused2
        type: s1
        doc: 0x0F. Unused.
      - id: visibility
        type: s4
        doc: 0x10-0x13. Globe object visibility/mobility bitfield.
    instances:
      lat_float:
        value: -y / 8.0
        doc: Latitude in degrees (north positive).
      lon_float:
        value: x / 8.0
        doc: Longitude in degrees.
enums:
  location_type:
    0: unused
    1: alien_ship
    2: xcom_ship
    3: xcom_base
    4: alien_base
    5: crash_site
    6: landed_ufo
    7: waypoint
    8: terror_site
