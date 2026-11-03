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
}

export interface SoldierSummary {
  index: number;
  name: string;
  rank: string;
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
  inventory: Record<string, number>;
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
}

export interface TransferSummary {
  index: number;
  origin: number;
  destination: number;
  hoursLeft: number;
  type: number;
  quantity: number;
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
}
