meta:
  id: purchase_dat
  title: Format of PURCHASE.DAT
  file-extension: DAT
  endian: le
doc: |
  PURCHASE.DAT holds prices and transfer times of everything that can be
  bought or sold.

  The file is 1820 bytes long (91 entries * 20 bytes). A new game copies the
  table from GEOSCAPE.EXE (offset 374512); only last_purchased changes during
  the game, all unknown fields are constants.
doc-ref: https://www.ufopaedia.org/index.php/PURCHASE.DAT
seq:
  - id: items
    type: purchase_item
    size: 20
    repeat: expr
    repeat-expr: 91
enums:
  category:
    0: item
    1: craft
    2: soldier
    3: scientist
    4: engineer
types:
  purchase_item:
    doc: The purchase data of a single item.
    seq:
      - id: category
        type: u1
        enum: category
        doc: "Category: 0 = item, 1 = craft, 2 = soldier, 3 = scientist, 4 = engineer."
      - id: index
        type: u1
        doc: Index of the item within its category (e.g. inventory index or craft type).
      - id: unknown1
        type: s4
        doc: |
          Unknown, 1000 for all items except the crafts X-COM can only
          manufacture (Lightning, Avenger, Firestorm), which have 0.
      - id: buy_price
        type: s4
        doc: Buy price in $. Also stored in LEASE.DAT when a craft is bought.
      - id: sell_price
        type: s4
        doc: Sell price in $.
      - id: transfer_hours
        type: u1
        doc: Transfer time in hours.
      - id: unknown2
        type: u1
        doc: Unknown, always 0.
      - id: last_purchased
        type: u1
        doc: |
          Quantity last purchased. The purchase screen of GEOSCAPE.EXE
          (0x1c8fa) reads bytes 16-17 as a uint16 order quantity (limited to
          255), so unknown3 is its high byte.
      - id: unknown3
        type: u1
        doc: Unknown, always 0 (high byte of last_purchased, see above).
      - id: purchasable
        type: u1
        doc: 1 if the item can be bought, 0 otherwise.
      - id: unknown4
        type: u1
        doc: Unknown, always 0.
