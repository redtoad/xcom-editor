meta:
  id: loc_dat
  file-extension: DAT
  endian: le
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
seq:
  - id: object
    type: object
    size: 2
    repeat: eos
types:
  object:
    seq:
      - id: type
        type: u1
        enum: object_type
