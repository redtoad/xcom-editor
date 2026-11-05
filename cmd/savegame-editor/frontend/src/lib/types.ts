export interface MissionInfo {
  turn: number;
  craftName: string;
  xcomAlive: number;
  xcomTotal: number;
  aliensAlive: number;
  aliensTotal: number;
  score: number;
}

export interface GameSummary {
  slot: string;
  title: string;
  time: string;
  soldierCount: number;
  baseCount: number;
  craftCount: number;
  mission: MissionInfo | null;
  // the savegame was changed on disk and needs to be reloaded
  stale: boolean;
}

export interface GameDetail {
  slot: string;
  title: string;
  time: string;
  soldierCount: number;
  baseCount: number;
  craftCount: number;
  balance: number;
  mission: MissionInfo | null;
  stale: boolean;
}

export interface SoldierSummary {
  index: number;
  name: string;
  rank: string;
  // (estimated) rank before being killed in action, see Soldier.LastRank
  lastRank: string;
  baseName: string;
  craftName: string;
  isDead: boolean;
  isWounded: boolean;
  missions: number;
  kills: number;
}

export interface SoldierDetail {
  index: number;
  name: string;
  rank: string;
  // (estimated) rank before being killed in action, see Soldier.LastRank
  lastRank: string;
  baseName: string;
  craftName: string;
  isDead: boolean;
  isWounded: boolean;
  missions: number;
  kills: number;
  recoveryDays: number;
  timeUnits: number;
  health: number;
  energy: number;
  reactions: number;
  strength: number;
  firingAccuracy: number;
  throwingAccuracy: number;
  meleeAccuracy: number;
  psionicStrength: number;
  psionicSkill: number;
  bravery: number;
  armor: string;
  gender: string;
  appearance: string;
  // inventory images of the soldier by armour, missing without game files
  images?: Record<string, string> | null;
  initialTimeUnits: number;
  initialHealth: number;
  initialEnergy: number;
  initialReactions: number;
  initialStrength: number;
  initialFiringAccuracy: number;
  initialThrowingAccuracy: number;
  initialMeleeAccuracy: number;
  initialPsionicStrength: number;
  initialPsionicSkill: number;
  initialBravery: number;
}

export interface BaseSummary {
  index: number;
  name: string;
  active: boolean;
  engineers: number;
  scientists: number;
  coord: string;
}

export interface BaseDetail {
  index: number;
  name: string;
  active: boolean;
  engineers: number;
  scientists: number;
  coord: string;
  tiles: TileInfo[];
  // all items that can be stored in a base, including those not in stores
  items: InventoryItem[];
  // sprites to draw the base with the graphics of the game (ground,
  // outlines, corridors, interiors, crafts), only if the game files are
  // available
  layers?: BaseSprite[][];
}

export interface InventoryItem {
  // index in the inventory of BASE.DAT
  index: number;
  name: string;
  count: number;
}

export interface BaseSprite {
  sprite: number;
  // position in pixels (of the original 32x32 tiles)
  x: number;
  y: number;
}

export interface TileInfo {
  type: string;
  daysToCompletion: number;
}

export interface CraftSummary {
  index: number;
  name: string;
  type: string;
  status: string;
  damage: number;
  fuel: number;
  baseName: string;
  // the entry in CRAFT.DAT is (nearly) empty, its values are meaningless
  damaged: boolean;
}

export interface TransferSummary {
  index: number;
  origin: number;
  destination: number;
  hoursLeft: number;
  type: number;
  quantity: number;
}

export interface ResearchProject {
  base: number;
  baseName: string;
  // index in RESEARCH.DAT
  topic: number;
  name: string;
  // scientist days remaining
  daysRemaining: number;
  scientists: number;
}

export interface ProductionProject {
  base: number;
  baseName: string;
  // index in PRODUCT.DAT
  product: number;
  name: string;
  // engineer hours needed to build one unit
  hoursPerUnit: number;
  // engineer hours of work done towards the next unit
  hoursDone: number;
  engineers: number;
  total: number;
  produced: number;
}

export interface Financials {
  currentBalance: number;
  expenditure: number[];
  maintenance: number[];
  balance: number[];
}

export interface GlobeLocation {
  type: string;
  typeCode: number;
  name: string;
  coord: { lat: number; lon: number };
  destCoord?: { lat: number; lon: number };
  craft?: CraftDetails;
  // detection capabilities of an X-COM base
  radars?: Radar[];
}

export interface Radar {
  facility: string;
  range: number; // nautical miles
  chance: number; // percent
}

export interface CraftDetails {
  type: string;
  status: string;
  altitude: string;
  speed: number;
  fuel: number;
  damage: number;
  baseName?: string;
  mission?: string;
  zone?: string;
  race?: string;
  destination?: string;
  damaged?: boolean;
}

export interface WorldPolygon {
  terrain: number;
  points: [number, number][]; // [lon, lat]
}

export interface WorldData {
  available: boolean;
  ocean?: string;
  polygons?: WorldPolygon[];
  textureSets?: number;
  terrains?: number;
}

export interface BattleUnit {
  index: number;
  name: string;
  owner: 'xcom' | 'alien' | 'civilian' | 'left-behind';
  x: number;
  y: number;
  z: number;
  large: boolean;
  alive: boolean;
  unconscious: boolean;
  visible: boolean;
  health: number;
  maxHealth: number;
  timeUnits: number;
  maxTimeUnits: number;
  morale: number;
  tank: boolean;
}

export interface BattleCraft {
  alien: boolean;
  name: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface BattleMap {
  turn: number;
  sizeX: number;
  sizeY: number;
  sizeZ: number;
  /** Flags per tile (see TILE_* in MissionMap.svelte), ordered along X, then Y, then Z. */
  tiles: number[];
  units: BattleUnit[];
  crafts: BattleCraft[];
}

export interface BattleTerrainPart {
  /** Index of the terrain set in BattleTerrain.sets. */
  set: number;
  /** Sprite in the sheet of the set (32x40 pixels per sprite). */
  frame: number;
  /** Vertical offset of the sprite in pixels. */
  y: number;
}

export interface BattleTerrain {
  /** Terrain sets of the map, their sprites are served at terrainSheetURL(name). */
  sets: string[];
  /** Sprite of each MCD index. */
  parts: BattleTerrainPart[];
  /** MCD indices of floor, west wall, north wall and object of each tile (4 per tile, order of BattleMap.tiles). */
  tiles: number[];
  /** Sprites of the units which can be drawn with the game graphics. */
  units: BattleUnitSprites[];
}

export interface UnitSprite {
  /** Sprite file, served at unitSheetURL(sheet). */
  sheet: string;
  /** Sprite in the sheet (32x40 pixels per sprite). */
  frame: number;
  /** Position relative to the sprite of the upper left tile of the unit, in pixels. */
  dx: number;
  dy: number;
}

export interface BattleUnitSprites {
  /** BattleUnit.index */
  index: number;
  /** Sprites in the order they are drawn. */
  sprites: UnitSprite[];
}
