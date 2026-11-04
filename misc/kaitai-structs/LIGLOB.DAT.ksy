meta:
  id: liglob_dat
  title: Format of LIGLOB.DAT
  file-extension: DAT
  endian: le
doc: |
  LIGLOB.DAT contains financial data for a savegame. This file is used by
  GEOSCAPE.EXE and its structure is very simple. Every record is a 4 byte
  signed long integer. Probably the most useful offset is the first 4 bytes
  where your current money is stored. The rest of the bytes are used for the
  Finance graphs: Expenditure, Maintenance, and Balance (the others are stored
  elsewhere). The file is 148 bytes long.

  Note: This file is only updated once a month, except for the first four
  bytes.
doc-ref: https://www.ufopaedia.org/index.php/LIGLOB.DAT
seq:
  - id: current_balance
    type: s4
    doc: |
      0x00-0x03: The most useful and probably the most hacked, these 4 bytes
      hold your current funds.
  - id: expenditure
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: |
      0x04-0x33: Expenditure for the months (exactly 12 of these 4 byte values)
      going in order January to December. The graph is adjusted to the current
      date so these essentially contain a year's worth of data.
  - id: maintenance
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: 0x34-0x63. Same as above only for maintenance costs.
  - id: balance
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: |
      0x64-0x93: Same as above only for balance (hence why you will sometimes
      see similar numbers here to your funds).
