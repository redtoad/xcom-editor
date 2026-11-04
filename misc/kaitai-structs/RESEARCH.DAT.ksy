meta:
  id: research_dat
  title: Format of RESEARCH.DAT
  file-extension: DAT
  endian: le
doc: |
  RESEARCH.DAT holds all research topics.

  The file is 2112 bytes long (96 topics * 22 bytes). A new game copies the
  table from GEOSCAPE.EXE (offset 371496); only unknown1, completed and
  prerequisites_met change during the game.
doc-ref: https://www.ufopaedia.org/index.php/RESEARCH.DAT
seq:
  - id: topics
    type: research_topic
    size: 22
    repeat: expr
    repeat-expr: 96
types:
  research_topic:
    doc: The data of a single research topic.
    seq:
      - id: time
        type: u2
        doc: |
          Base research time in scientist days. It is randomized by +/-50% when
          the research starts and then copied to PROJECT.DAT.
      - id: topic
        type: u1
        repeat: expr
        repeat-expr: 3
        doc: |
          Topic code. For live aliens, byte 3 is the alien race and byte 4 the
          alien rank (as in ASTORE.DAT). If byte 2 or 3 is not 0, the item must
          be in stores before the topic can be researched.
      - id: unknown1
        type: u1
        doc: |
          Unknown, 0 in the initial table. For live alien topics this is the
          index in ASTORE.DAT of the first matching alien (same race and rank)
          in the current base. GEOSCAPE.EXE stores it while building the
          research list of a base (0x19fd3) and does not update it afterwards,
          so it may be stale (e.g. after the alien was researched). No code
          reading it has been found.
      - id: ufopaedia_category
        type: u1
        doc: UFOpaedia category.
      - id: ufopaedia_entry
        type: u1
        doc: UFOpaedia entry within the category.
      - id: points
        type: u1
        doc: Score points awarded on completion.
      - id: unknown2
        type: u1
        doc: Unknown, always 0. Not accessed by any code.
      - id: completed
        type: s4
        doc: Non-zero if the research is completed.
      - id: prerequisite
        type: s4
        doc: Prerequisite flag (exact meaning unclear). 0 or 1, never changes during the game.
      - id: prerequisites_met
        type: s4
        doc: Non-zero if all prerequisites are met.
