package savegame

import (
	"fmt"
	"math"
)

// Location is a real-world GPS coordinate.
type Location struct {
	// Horizontal coordinates or longitude
	Lon float32
	// Vertical coordinates or latitude
	Lat float32
}

func NewLocation(x int, y int) Location {
	return Location{
		Lon: 180.0 - Longitude(x),
		Lat: 180.0 - Latitude(y),
	}
}

func (l Location) String() string {
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
