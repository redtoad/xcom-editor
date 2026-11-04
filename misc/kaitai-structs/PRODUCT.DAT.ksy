meta:
  id: product_dat
  title: Format of PRODUCT.DAT
  file-extension: DAT
  endian: le
doc: |
  PRODUCT.DAT holds the manufacturing data of all manufacturable items.

  The file is 630 bytes long (35 items * 18 bytes). A new game copies the
  table from GEOSCAPE.EXE (offset 373608); only not_researched changes during
  the game, all unknown fields are constants.
doc-ref: https://www.ufopaedia.org/index.php/PRODUCT.DAT
seq:
  - id: products
    type: product
    size: 18
    repeat: expr
    repeat-expr: 35
types:
  product:
    doc: The manufacturing data of a single item.
    seq:
      - id: cost
        type: s4
        doc: Manufacturing cost in $.
      - id: hours
        type: u2
        doc: Engineer hours needed to build one item.
      - id: work_space
        type: u2
        doc: Workshop space needed.
      - id: requirements
        type: u1
        repeat: expr
        repeat-expr: 3
        doc: Required items (inventory index in BASE.DAT), 0xFF if none.
      - id: required_quantities
        type: u1
        repeat: expr
        repeat-expr: 3
        doc: Quantity needed of each required item.
      - id: item
        type: u1
        doc: Manufactured item (inventory index in BASE.DAT).
      - id: unknown1
        type: u1
        doc: Unknown, always 1.
      - id: not_researched
        type: u1
        doc: 0 if the item can be built, 1 if it still needs to be researched.
      - id: unknown2
        type: u1
        doc: Unknown, always 0.
