package battlescape

import "encoding/binary"

// RoutesFile holds the route nodes used for spawning units and AI path
// finding.
//
// The file is 24 bytes per node long (see WglobFile.Routes).
//
// https://www.ufopaedia.org/index.php/ROUTES.DAT
type RoutesFile struct {
	Nodes []RouteNode
}

// RouteNode is a single route node.
type RouteNode struct {
	// 0-2: Position of the node.
	X uint8
	Y uint8
	Z uint8
	// 3: Map module the node belongs to.
	Module uint8
	// 4-18: Links to other nodes.
	Links [5]RouteLink
	// 19: Unit type that can use the node.
	UnitType uint8
	// 20: Rank of the unit that can spawn here.
	Rank uint8
	// 21: Flags (attack priority of base modules or patrol importance).
	Flags uint8
	// 22: Unknown, usually 0.
	Unknown uint8
	// 23: Spawn priority, 0 = no spawning.
	Spawn uint8
}

// RouteLink is a link from one route node to another.
type RouteLink struct {
	// Target node, 0xFF if unused, 0xFB-0xFE for exits to the west, south,
	// east and north.
	Target uint8
	// Distance to the target node.
	Distance uint8
	// Unit type that can use the link.
	UnitType uint8
}

const routeNodeSize = 24

// SizeOf implements restruct.Sizer.
func (f RoutesFile) SizeOf() int {
	return len(f.Nodes) * routeNodeSize
}

// Unpack implements restruct.Unpacker.
func (f *RoutesFile) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	nodes, err := unpackRecords[RouteNode](buf, order, routeNodeSize)
	if err != nil {
		return nil, err
	}
	f.Nodes = nodes
	return buf[f.SizeOf():], nil
}

// Pack implements restruct.Packer.
func (f RoutesFile) Pack(buf []byte, order binary.ByteOrder) ([]byte, error) {
	return buf, packRecords(buf, order, routeNodeSize, f.Nodes)
}
