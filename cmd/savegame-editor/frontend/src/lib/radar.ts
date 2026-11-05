import { geoCircle } from 'd3-geo';
import type { GlobeLocation, Radar } from './types';

// Colors of the radar facilities.
export const radarColors: Record<string, string> = {
  'Small Radar System': '#b6ff6e',
  'Large Radar System': '#00ff41',
  'Hyperwave Decoder': '#d070ff',
};

export function radarColor(radar: Radar): string {
  return radarColors[radar.facility] || '#ffffff';
}

// A nautical mile is one arc minute of a great circle.
export function rangeToDegrees(nauticalMiles: number): number {
  return nauticalMiles / 60;
}

export interface RadarCoverage {
  base: GlobeLocation;
  radar: Radar;
  // the area within range as a GeoJSON polygon
  area: GeoJSON.Polygon;
}

// radarCoverage returns the area within range of each radar of the bases,
// ordered by range (longest first) so that shorter ranges are drawn on top.
export function radarCoverage(locations: GlobeLocation[]): RadarCoverage[] {
  return locations
    .flatMap(base => (base.radars ?? []).map(radar => ({
      base,
      radar,
      area: geoCircle()
        .center([base.coord.lon, base.coord.lat])
        .radius(rangeToDegrees(radar.range))
        .precision(2)() as GeoJSON.Polygon,
    })))
    .sort((a, b) => b.radar.range - a.radar.range);
}
