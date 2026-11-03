import type { GeoProjection } from 'd3-geo';
import type { WorldPolygon } from './types';

/** Resolution of the terrain map in pixels per degree. */
export const PIXELS_PER_DEGREE = 1.5;

/** Loads images and resolves once all of them are available. */
export function loadImages(urls: string[]): Promise<HTMLImageElement[]> {
  return Promise.all(
    urls.map(
      url =>
        new Promise<HTMLImageElement>((resolve, reject) => {
          const img = new Image();
          img.onload = () => resolve(img);
          img.onerror = () => reject(new Error(`could not load ${url}`));
          img.src = url;
        }),
    ),
  );
}

/**
 * Converts the points of a polygon into a ring that can be drawn on an
 * equirectangular map: longitudes are made continuous (so polygons crossing
 * the antimeridian are not stretched across the map) and points on a pole are
 * expanded along the edge of the map.
 */
export function unwrapRing(points: [number, number][]): [number, number][] {
  const isPole = (lat: number) => Math.abs(lat) >= 90;
  const near = (lon: number, ref: number) => lon - 360 * Math.round((lon - ref) / 360);

  // make the longitudes of all points (except poles) continuous
  const lons: number[] = [];
  let prev: number | null = null;
  for (const [lon, lat] of points) {
    if (isPole(lat)) {
      lons.push(NaN);
      continue;
    }
    prev = prev === null ? lon : near(lon, prev);
    lons.push(prev);
  }

  const n = points.length;
  const neighbour = (i: number, step: number) => {
    for (let j = (i + step + n) % n; j !== i; j = (j + step + n) % n) {
      if (!isNaN(lons[j])) return lons[j];
    }
    return 0;
  };

  const ring: [number, number][] = [];
  points.forEach(([, lat], i) => {
    if (isPole(lat)) {
      // a pole touches the whole edge of the map between its neighbours
      const before = neighbour(i, -1);
      ring.push([before, lat], [near(neighbour(i, 1), before), lat]);
    } else {
      ring.push([lons[i], lat]);
    }
  });
  return ring;
}

/**
 * Draws the terrain into an equirectangular map. The textures are tiled in
 * map coordinates, so they stick to the surface of the globe.
 */
export function buildTerrainMap(
  polygons: WorldPolygon[],
  textures: HTMLImageElement[],
  ocean: string,
  pixelsPerDegree = PIXELS_PER_DEGREE,
): ImageData {
  const width = Math.round(360 * pixelsPerDegree);
  const height = Math.round(180 * pixelsPerDegree);
  const canvas = document.createElement('canvas');
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d')!;
  ctx.imageSmoothingEnabled = false;
  ctx.fillStyle = ocean;
  ctx.fillRect(0, 0, width, height);

  const patterns = textures.map(img => ctx.createPattern(img, 'repeat')!);
  for (const poly of polygons) {
    const ring = unwrapRing(poly.points);
    // draw the polygon (and its copies left and right of the map)
    for (const shift of [-width, 0, width]) {
      ctx.beginPath();
      ring.forEach(([lon, lat], i) => {
        const x = (lon + 180) * pixelsPerDegree + shift;
        const y = (90 - lat) * pixelsPerDegree;
        if (i === 0) ctx.moveTo(x, y);
        else ctx.lineTo(x, y);
      });
      ctx.closePath();
      ctx.fillStyle = patterns[poly.terrain];
      ctx.strokeStyle = patterns[poly.terrain];
      ctx.lineWidth = 1; // closes gaps between neighbouring polygons
      ctx.fill();
      ctx.stroke();
    }
  }
  return ctx.getImageData(0, 0, width, height);
}

/**
 * Renders the visible hemisphere of the terrain map into ctx (pixel by
 * pixel) using the inverse of the given (orthographic) projection.
 */
export function renderGlobe(
  ctx: CanvasRenderingContext2D,
  projection: GeoProjection,
  terrain: ImageData,
  pixelsPerDegree = PIXELS_PER_DEGREE,
) {
  const { width, height } = ctx.canvas;
  const [cx, cy] = projection.translate();
  const r = projection.scale();
  const out = ctx.createImageData(width, height);
  for (let py = 0; py < height; py++) {
    for (let px = 0; px < width; px++) {
      const dx = px + 0.5 - cx;
      const dy = py + 0.5 - cy;
      if (dx * dx + dy * dy > r * r) continue;
      const p = projection.invert!([px + 0.5, py + 0.5]);
      if (!p) continue;
      let mx = Math.floor((p[0] + 180) * pixelsPerDegree) % terrain.width;
      if (mx < 0) mx += terrain.width;
      const my = Math.min(terrain.height - 1, Math.max(0, Math.floor((90 - p[1]) * pixelsPerDegree)));
      const src = (my * terrain.width + mx) * 4;
      const dst = (py * width + px) * 4;
      out.data[dst] = terrain.data[src];
      out.data[dst + 1] = terrain.data[src + 1];
      out.data[dst + 2] = terrain.data[src + 2];
      out.data[dst + 3] = 255;
    }
  }
  ctx.putImageData(out, 0, 0);
}
