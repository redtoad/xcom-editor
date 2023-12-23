package savegame

import (
	"fmt"
	"math"
)

// Coord is a real-world GPS coordinate.
type Coord struct {
	// Horizontal coordinates or longitude
	Lon float32
	// Vertical coordinates or latitude
	Lat float32
}

func NewCoord(x int, y int) Coord {
	return Coord{
		// The X coordinate starts with X = 0 at 0° longitude (the Prime Meridian or Greenwich
		// Meridian) and increases going eastward. Unlike real-world longitude, there is only
		// "East" longitude, from 0°E to 359.875°E. This is the result of forcing the
		// coordinate system to be positive-only, for algorithmic purposes. So, for example,
		// the equivalent of 90°W longitude would be "270°E longitude", with a game X
		// coordinate of 270 x 8 = 2160.
		Lon: float32(x) / (2880.0 / 360.0),
		// A Y coordinate value of 720 corresponds to 90.0 "90° S" latitude (the South Pole);
		// a Y coordinate value of -720 corresponds to "90° N" latitude (the North Pole).
		Lat: float32(y) / (720.0 / 90.0),
	}
}

func (l Coord) String() string {
	latDir := "N"
	if l.Lat < 0 {
		latDir = "S"
	}
	lonDir := "E"
	if l.Lon < 0 {
		lonDir = "W"
	}
	return fmt.Sprintf(
		"%.5f° %s %.5f° %s",
		math.Abs(float64(l.Lon)), lonDir,
		math.Abs(float64(l.Lat)), latDir,
	)
}

// TerrainHexColorsXCom stores a color map for the world terrain types
// (see geodata.Polygon#Terrain).
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
