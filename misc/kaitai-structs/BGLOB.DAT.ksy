meta:
  id: bglob_dat
  title: Format of BGLOB.DAT
  file-extension: DAT
  endian: le
doc: |
  BGLOB.DAT holds the ambient lighting of the Battlescape map. Light values
  range from 0 (bright) to 16 (dark), e.g. 0 for day missions and 16 for
  night missions.

  The file is 2 bytes long. It is only present in savegames that were saved
  during a mission.
doc-ref: https://www.ufopaedia.org/index.php/BGLOB.DAT
seq:
  - id: ambient
    type: u1
    doc: Overall ambient light.
  - id: shadow
    type: u1
    doc: |
      Ambient light in shadowed areas (e.g. inside crafts), usually
      ambient + 4.
