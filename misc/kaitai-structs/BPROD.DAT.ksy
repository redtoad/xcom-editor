meta:
  id: bprod_dat
  title: Format of BPROD.DAT
  file-extension: DAT
  endian: le
doc: |
  BPROD.DAT holds the manufacturing progress of each base.

  The file is 2800 bytes long (8 bases * 350 bytes).
doc-ref: https://www.ufopaedia.org/index.php/BPROD.DAT
seq:
  - id: bases
    type: base_production
    size: 350
    repeat: expr
    repeat-expr: 8
types:
  base_production:
    doc: |
      The manufacturing progress of a single base. All arrays have 35 entries
      and are indexed in the same order as PRODUCT.DAT. Completed projects are
      reset to 0 in all fields.
    seq:
      - id: hours_stored
        type: s4
        repeat: expr
        repeat-expr: 35
        doc: |
          Hours of work done towards the next unit. An item is only shown as in
          production if this value is not 0.
      - id: engineers
        type: s2
        repeat: expr
        repeat-expr: 35
        doc: Engineers allocated.
      - id: total
        type: s2
        repeat: expr
        repeat-expr: 35
        doc: Total number of units to produce.
      - id: produced
        type: s2
        repeat: expr
        repeat-expr: 35
        doc: Number of units produced.
