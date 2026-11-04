meta:
  id: alien_dat
  title: Format of ALIEN.DAT
  file-extension: DAT
  endian: le
doc: |
  ALIEN.DAT contains alien activity scores for countries and regions, used to
  populate the graphs screen. The same format is used for XCOM.DAT (which
  records X-COM activity instead). There are 12 records of 124 bytes (31
  signed 4-byte integers) each; the first record is the current month, the
  second is last month, and so on. The file is 1,488 bytes long.

  Each score increments by 1 every 30 minutes that a craft is over that
  territory. For countries, the matching region is also incremented.
doc-ref: https://www.ufopaedia.org/index.php/ALIEN.DAT
seq:
  - id: months
    type: activity_month
    size: 124
    repeat: expr
    repeat-expr: 12
types:
  activity_month:
    doc: |
      Activity scores for all countries and regions for a single month. Values
      are signed and may go negative in XCOM.DAT.
    seq:
      - id: usa
        type: s4
        doc: 0x00 (country)
      - id: russia
        type: s4
        doc: 0x04 (country)
      - id: uk
        type: s4
        doc: 0x08 (country)
      - id: france
        type: s4
        doc: 0x0C (country)
      - id: germany
        type: s4
        doc: 0x10 (country)
      - id: italy
        type: s4
        doc: 0x14 (country)
      - id: spain
        type: s4
        doc: 0x18 (country)
      - id: china
        type: s4
        doc: 0x1C (country)
      - id: japan
        type: s4
        doc: 0x20 (country)
      - id: india
        type: s4
        doc: 0x24 (country)
      - id: brazil
        type: s4
        doc: 0x28 (country)
      - id: australia
        type: s4
        doc: 0x2C (country)
      - id: nigeria
        type: s4
        doc: 0x30 (country)
      - id: south_africa
        type: s4
        doc: 0x34 (country)
      - id: egypt
        type: s4
        doc: 0x38 (country)
      - id: canada
        type: s4
        doc: 0x3C (country)
      - id: north_america
        type: s4
        doc: 0x40 (region)
      - id: arctic
        type: s4
        doc: 0x44 (region)
      - id: antarctica
        type: s4
        doc: 0x48 (region)
      - id: south_america
        type: s4
        doc: 0x4C (region)
      - id: europe
        type: s4
        doc: 0x50 (region)
      - id: north_africa
        type: s4
        doc: 0x54 (region)
      - id: southern_africa
        type: s4
        doc: 0x58 (region)
      - id: central_asia
        type: s4
        doc: 0x5C (region)
      - id: southeast_asia
        type: s4
        doc: 0x60 (region)
      - id: siberia
        type: s4
        doc: 0x64 (region)
      - id: australasia
        type: s4
        doc: 0x68 (region)
      - id: pacific
        type: s4
        doc: 0x6C (region)
      - id: north_atlantic
        type: s4
        doc: 0x70 (region)
      - id: south_atlantic
        type: s4
        doc: 0x74 (region)
      - id: indian_ocean
        type: s4
        doc: 0x78 (region)
