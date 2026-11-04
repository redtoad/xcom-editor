meta:
  id: project_dat
  title: Format of PROJECT.DAT
  file-extension: DAT
  endian: le
doc: |
  PROJECT.DAT holds the research progress of each base.

  The file is 2304 bytes long (8 bases * 288 bytes).
doc-ref: https://www.ufopaedia.org/index.php/PROJECT.DAT
seq:
  - id: bases
    type: base_research
    size: 288
    repeat: expr
    repeat-expr: 8
types:
  base_research:
    doc: |
      The research progress of a single base. Both arrays have 96 entries and
      are indexed in the same order as RESEARCH.DAT.
    seq:
      - id: days_remaining
        type: u2
        repeat: expr
        repeat-expr: 96
        doc: Scientist days remaining. 0 if the topic is not being researched.
      - id: scientists
        type: u1
        repeat: expr
        repeat-expr: 96
        doc: Scientists assigned.
