meta:
  id: routes_dat
  title: Format of ROUTES.DAT
  file-extension: DAT
  endian: le
doc: |
  ROUTES.DAT holds the route nodes used for spawning units and AI path
  finding, 24 bytes per node.

  The number of used nodes is stored in WGLOB.DAT (`routes`). The game does
  not truncate the file, so stale nodes of a previous, larger map may follow
  the used part. As a single description cannot read WGLOB.DAT, all records
  up to the end of the file are parsed.
doc-ref: https://www.ufopaedia.org/index.php/ROUTES.DAT
seq:
  - id: nodes
    type: route_node
    repeat: eos
types:
  route_node:
    seq:
      - id: x
        type: u1
      - id: y
        type: u1
      - id: z
        type: u1
      - id: module
        type: u1
        doc: Map module the node belongs to.
      - id: links
        type: route_link
        repeat: expr
        repeat-expr: 5
        doc: Links to other nodes.
      - id: unit_type
        type: u1
        doc: Unit type that can use the node.
      - id: rank
        type: u1
        doc: Rank of the unit that can spawn here.
      - id: flags
        type: u1
        doc: Flags (attack priority of base modules or patrol importance).
      - id: unknown
        type: u1
        doc: Unknown, usually 0.
      - id: spawn
        type: u1
        doc: Spawn priority, 0 = no spawning.
  route_link:
    seq:
      - id: target
        type: u1
        doc: |
          Target node, 0xFF if unused, 0xFB-0xFE for exits to the west,
          south, east and north.
      - id: distance
        type: u1
        doc: Distance to the target node.
      - id: unit_type
        type: u1
        doc: Unit type that can use the link.
