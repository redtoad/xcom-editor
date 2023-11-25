meta:
  id: libglob_dat
  file-extension: LIBGLOB.DAT
  endian: le
doc: |
  This file is used by GEOSCAPE.EXE and it's structure is very simple. Every 
  record is a 4 byte signed long integer. Probably the most useful offset is the 
  first 4 bytes where your current money is stored. The rest of the bytes are 
  used for the Finance graphs: Expenditure, Maintenance, and Balance (the others 
  are stored elsewhere).
  
  Note: This file is only updated once a month, except for the first four bytes.
seq:
  - id: current_money
    type: s4
    doc: |
      The most useful and probably the most hacked, these 4 bytes hold your 
      current funds.
  - id: expenditure
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: |
      This is the expenditure for the months (exactly 12 of these 4 byte values) 
      going in order January to December. The graph is adjusted to the current 
      date so these essentially contain a years worth of data.
  - id: maintenance
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: |
      Same as above only for Maintenance.
  - id: balance
    type: s4
    repeat: expr
    repeat-expr: 12
    doc: |
      Same as above only for Balance (hence why you will sometimes see similar 
      numbers here to your funds).