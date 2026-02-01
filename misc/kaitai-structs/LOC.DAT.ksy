meta:
  id: loc_dat
  file-extension: DAT
  endian: le
doc-ref: https://www.ufopaedia.org/index.php?title=LOC.DAT
seq:
  - id: location
    type: location
    repeat: eos
types:
  location:
    seq:
      - id: type
        type: u1
        enum: object_type
      - id: reference
        type: u1
      - id: lon
        type: u2
      - id: lat
        type: s2
      - id: tbd
        size: 14
    instances:
      lat_float:
        value: -lat / 8.0
      lon_float:
        value:  (lon / 8.0)
enums:
  object_type:
    0x00: unused
    0x01: alien_ship
    0x02: xcom_ship
    0x03: xcom_base
    0x04: alien_base
    0x05: crash_site
    0x06: landed_ufo
    0x07: waypoint
    0x08: terror_site
    # Extras for TFTD only:
    0x51: port_attack
    0x52: island_attack
    0x53: passenger_cargo_ship
    0x54: artefact_site