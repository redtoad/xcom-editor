meta:
  id: diplom_dat
  title: Format of DIPLOM.DAT
  file-extension: DAT
  endian: le
doc: |
  DIPLOM.DAT holds the relations to the 16 funding countries (in the same
  order as in ALIEN.DAT).

  The file is 576 bytes long (16 countries * 36 bytes).
doc-ref: https://www.ufopaedia.org/index.php/DIPLOM.DAT
seq:
  - id: countries
    type: country_relation
    size: 36
    repeat: expr
    repeat-expr: 16
enums:
  satisfaction:
    0: pact_with_aliens
    1: unhappy
    2: satisfied
    3: happy
types:
  country_relation:
    doc: The relation to a single funding country.
    seq:
      - id: name_index
        type: s2
        doc: Index of the country name in ENGLISH.DAT (starts at 602).
      - id: satisfaction
        type: s2
        enum: satisfaction
        doc: |
          Satisfaction: 3 = happy, 2 = satisfied, 1 = unhappy, 0 = signed a
          pact with the aliens (stops funding).
      - id: funding
        type: s2
        repeat: expr
        repeat-expr: 12
        doc: Monthly funding in $1000 (January to December).
      - id: funding_change
        type: s2
        doc: Change in funding (in $1000) between last month and this month.
      - id: pact
        type: s2
        doc: Non-zero if the country signs a pact with the aliens at the end of the month.
      - id: starting_funding
        type: s2
        doc: Funding at the start of the game (in $1000). No effect after the first save.
      - id: funding_cap
        type: s2
        doc: Maximum funding (in $1000).
