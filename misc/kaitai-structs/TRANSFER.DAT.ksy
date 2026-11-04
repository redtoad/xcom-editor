meta:
  id: transfer_dat
  title: Format of TRANSFER.DAT
  file-extension: DAT
  endian: le
doc: |
  TRANSFER.DAT contains all information about items in transit. Each record
  is 8 bytes long and the file is fixed at 100 entries, thus a fixed size of
  800 bytes.
doc-ref: https://www.ufopaedia.org/index.php/TRANSFER.DAT
seq:
  - id: transfers
    type: transfer_data
    size: 8
    repeat: expr
    repeat-expr: 100
types:
  transfer_data:
    seq:
      - id: origin
        type: u1
        doc: |
          Base the item is coming from (as indexed in LOC.DAT). 255 if the item
          is purchased and thus has no base of origin.
      - id: destination
        type: u1
        doc: Destination the item is going to (again from LOC.DAT). 255 should not be used here.
      - id: hours_left
        type: u1
        doc: |
          Hours left in transit. NOTE: Setting this to 0 will make the game
          think it has been completed already.
      - id: type
        type: u1
        doc: Item type. This also affects what can be used in the next field.
      - id: reference_number
        type: s2
        doc: Reference number. The meaning of this value depends on the item type.
      - id: quantity
        type: u1
        doc: |
          Quantity. The entry is ignored if this value is 0, thus there can be
          invalid data in the other fields but they will always have this byte
          set to 0.
      - id: unused
        type: u1
        doc: Unused padding byte. Not kept by the Go code, which writes it as 0.
