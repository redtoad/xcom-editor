meta:
  id: xcom_dat
  title: Format of XCOM.DAT
  file-extension: DAT
  endian: le
doc: |
  XCOM.DAT contains X-COM activity scores for countries and regions, used to
  populate the graphs screen. It uses the same format as ALIEN.DAT (which
  records alien activity instead). Values may go negative.

  The file holds 12 months of 31 signed 32-bit scores (1488 bytes). The first
  record is the current month, the second is last month, and so on.
doc-ref: https://www.ufopaedia.org/index.php/XCOM.DAT
seq:
  - id: months
    type: activity_month
    size: 124
    repeat: expr
    repeat-expr: 12
types:
  activity_month:
    doc: Activity scores for all countries and regions for a single month.
    seq:
      # Countries (offsets 0x00-0x3C)
      - id: usa
        type: s4
      - id: russia
        type: s4
      - id: uk
        type: s4
      - id: france
        type: s4
      - id: germany
        type: s4
      - id: italy
        type: s4
      - id: spain
        type: s4
      - id: china
        type: s4
      - id: japan
        type: s4
      - id: india
        type: s4
      - id: brazil
        type: s4
      - id: australia
        type: s4
      - id: nigeria
        type: s4
      - id: south_africa
        type: s4
      - id: egypt
        type: s4
      - id: canada
        type: s4
      # Regions (offsets 0x40-0x78)
      - id: north_america
        type: s4
      - id: arctic
        type: s4
      - id: antarctica
        type: s4
      - id: south_america
        type: s4
      - id: europe
        type: s4
      - id: north_africa
        type: s4
      - id: southern_africa
        type: s4
      - id: central_asia
        type: s4
      - id: south_east_asia
        type: s4
      - id: siberia
        type: s4
      - id: australasia
        type: s4
      - id: pacific
        type: s4
      - id: north_atlantic
        type: s4
      - id: south_atlantic
        type: s4
      - id: indian_ocean
        type: s4
