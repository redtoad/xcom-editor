import type { BattleMap } from './types';

// Flags of the tiles in BattleMap.tiles (see battleMapJS in api.go).
export const TILE_GROUND = 1;
export const TILE_WEST_WALL = 2;
export const TILE_NORTH_WALL = 4;
export const TILE_OBJECT = 8;
export const TILE_SEEN = 16;

/** Returns the flags of a tile. Z is the level in the order of MAP.DAT (0 = top level). */
export function tileAt(m: BattleMap, x: number, y: number, z: number): number {
  return m.tiles[z * m.sizeX * m.sizeY + y * m.sizeX + x];
}

/**
 * Returns the name of level z. Levels are stored from the top level down,
 * level 1 is the ground level.
 */
export function levelLabel(m: BattleMap, z: number): string {
  return `Level ${m.sizeZ - z}`;
}

/** Darkens a #rrggbb color (factor 0 = black, 1 = unchanged). */
export function darken(color: string, factor: number): string {
  const m = /^#([0-9a-f]{2})([0-9a-f]{2})([0-9a-f]{2})$/i.exec(color.trim());
  if (!m) return color;
  const c = m.slice(1).map(h => Math.round(parseInt(h, 16) * factor));
  return `rgb(${c[0]}, ${c[1]}, ${c[2]})`;
}
