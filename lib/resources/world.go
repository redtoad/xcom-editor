package resources

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/go-restruct/restruct"
)

// WorldData describes the terrain on the geoscape screen using quadrilateral polygons and triangles
// loaded from GEODATA/WORLD.DAT.
//
// The first 16 bytes of file contain the points for the polygon. 4 sets of 2 short (2-byte) integers,
// designating the 'X' and 'Y' coordinate (or longitude and latitude respectively, if you prefer). If
// the last set has an x value of -1 then it is to be rendered as a triangle, otherwise it is a quad.
//
// The last 4 bytes in the record contain the terrain type. This could be a long integer or 2 short
// integers as the last 2 bytes in each record are 0.
//
// See https://www.ufopaedia.org/index.php/WORLD.DAT for more information.
type WorldData struct {
	Polygons []Polygon
}

func (w *WorldData) Unpack(buf []byte, order binary.ByteOrder) ([]byte, error) {

	reader := bytes.NewReader(buf)
	for {
		data := make([]byte, 20)
		noBytes, err := reader.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not read next data chunk: %w", err)
		}
		var poly Polygon
		if noBytes != poly.SizeOf() {
			return nil, fmt.Errorf("not enough to read all polygon data")
		}
		if err := restruct.Unpack(data, order, &poly); err != nil {
			return nil, fmt.Errorf("could not unpack polygon: %w", err)
		}
		w.Polygons = append(w.Polygons, poly)
	}
	return []byte{}, nil
}

// LoadWorldData loads data from WORLD.DATA.
func LoadWorldData(gameDir string) (WorldData, error) {
	fp := path.Join(gameDir, "./GEODATA/WORLD.DAT")
	buf, err := os.ReadFile(fp)
	if err != nil {
		return WorldData{}, fmt.Errorf("could not open WORLD.DAT: %w", err)
	}
	var world WorldData
	if err = restruct.Unpack(buf, binary.LittleEndian, &world); err != nil {
		return WorldData{}, fmt.Errorf("could not parse WORLD.DAT: %w", err)
	}
	return world, nil
}

type Polygon struct {

	// First X coordinate/longitude
	X0 int `struct:"int16"`
	// First Y coordinate/latitude
	Y0 int `struct:"int16"`
	// Second X coordinate/longitude
	X1 int `struct:"int16"`
	// Second Y coordinate/latitude
	Y1 int `struct:"int16"`
	// Third X coordinate/longitude
	X2 int `struct:"int16"`
	// Third Y coordinate/latitude
	Y2 int `struct:"int16"`
	// Fourth* X coordinate/longitude
	X3 int `struct:"int16"`
	// Fourth* Y coordinate/latitude
	Y3 int `struct:"int16"`

	// Terrain Type/Texture	0-12
	Terrain int `struct:"int32"`
}

func (p *Polygon) Type() PolygonType {
	if p.X3 == -1 {
		return Triangle
	}
	return QuadrilateralPolygon
}

func (p Polygon) SizeOf() int { return 20 }

func (p Polygon) String() string {
	return fmt.Sprintf(
		"P{(%d,%d) (%d,%d) (%d,%d) (%d,%d) terrain=%d}",
		p.X0, p.Y0, p.X1, p.Y1, p.X2, p.Y2, p.X3, p.Y3, p.Terrain,
	)
}

type PolygonType int

const (
	Triangle PolygonType = iota
	QuadrilateralPolygon
)

// Longitude converts an X coordinate to longitude east.
//
// The X coordinate starts with X = 0 at 0° longitude (the Prime Meridian or Greenwich Meridian) and
// increases going eastward. Unlike real-world longitude, there is only "East" longitude, from 0°E
// to 359.875°E. This is the result of forcing the coordinate system to be positive-only, for
// algorithmic purposes. So, for example, the equivalent of 90°W longitude would be "270°E longitude",
// with a game X coordinate of 270 x 8 = 2160.
func Longitude(x int) float32 {
	return float32(x) / (2880.0 / 360.0)
}

// Latitude converts a Y coordinate to latitude.
//
// A Y coordinate value of 720 corresponds to 90.0 "90° S" latitude (the South Pole); a Y coordinate
// value of -720 corresponds to "90° N" latitude (the North Pole).
func Latitude(y int) float32 {
	return float32(y) / (720.0 / 90.0)
}

var TerrainHexColorsXCom = []string{
	"#7FFE11", // Forest / Jungle
	"#6CD911", // Farm
	"#5AB411", // Farm
	"#478F11", // Farm
	"#356A11", // Farm
	"#224511", // Mountain
	"#464645", // Forest / Jungle
	"#6B6B45", // Desert
	"#909045", // Desert
	"#B5B545", // Polar Ice
	"#DADA45", // Forest / Jungle
	"#FFFF45", // Forest / Jungle
	"#FFFFFF", // Polar Seas w/Icebergs
}
