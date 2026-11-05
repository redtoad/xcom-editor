import type {
  GameSummary, GameDetail, SoldierSummary, SoldierDetail,
  BaseSummary, BaseDetail, CraftSummary, TransferSummary, Financials,
  ResearchProject, ProductionProject,
  GlobeLocation, WorldData, BattleMap, BattleTerrain
} from './types';

async function fetchJSON<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(url, options);
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(body.error || res.statusText);
  }
  return res.json();
}

export async function listGames(): Promise<GameSummary[]> {
  return fetchJSON('/api/games');
}

export async function getGame(slot: string): Promise<GameDetail> {
  return fetchJSON(`/api/games/${slot}`);
}

export async function listSoldiers(slot: string): Promise<SoldierSummary[]> {
  return fetchJSON(`/api/games/${slot}/soldiers`);
}

export async function getSoldier(slot: string, idx: number): Promise<SoldierDetail> {
  return fetchJSON(`/api/games/${slot}/soldiers/${idx}`);
}

export async function updateSoldier(slot: string, idx: number, updates: Partial<SoldierDetail>): Promise<void> {
  await fetchJSON(`/api/games/${slot}/soldiers/${idx}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(updates),
  });
}

export async function resurrectSoldier(slot: string, idx: number): Promise<void> {
  await fetchJSON(`/api/games/${slot}/soldiers/${idx}/resurrect`, { method: 'POST' });
}

export async function healSoldier(slot: string, idx: number): Promise<void> {
  await fetchJSON(`/api/games/${slot}/soldiers/${idx}/heal`, { method: 'POST' });
}

export async function listBases(slot: string): Promise<BaseSummary[]> {
  return fetchJSON(`/api/games/${slot}/bases`);
}

export async function getBase(slot: string, idx: number): Promise<BaseDetail> {
  return fetchJSON(`/api/games/${slot}/bases/${idx}`);
}

export async function updateBase(slot: string, idx: number, updates: Record<string, any>): Promise<void> {
  await fetchJSON(`/api/games/${slot}/bases/${idx}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(updates),
  });
}

// adjustInventory changes the quantity of an item in a base by delta and
// returns the new quantity.
export async function adjustInventory(slot: string, base: number, item: number, delta: number): Promise<number> {
  const res = await fetchJSON<{ count: number }>(`/api/games/${slot}/bases/${base}/inventory/${item}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ delta }),
  });
  return res.count;
}

export async function listCraft(slot: string): Promise<CraftSummary[]> {
  return fetchJSON(`/api/games/${slot}/craft`);
}

export async function listTransfers(slot: string): Promise<TransferSummary[]> {
  return fetchJSON(`/api/games/${slot}/transfers`);
}

export async function listResearch(slot: string): Promise<ResearchProject[]> {
  return fetchJSON(`/api/games/${slot}/research`);
}

export async function speedupResearchProject(slot: string, base: number, topic: number): Promise<void> {
  await fetchJSON(`/api/games/${slot}/research/${base}/${topic}/speedup`, { method: 'POST' });
}

export async function listProduction(slot: string): Promise<ProductionProject[]> {
  return fetchJSON(`/api/games/${slot}/production`);
}

export async function speedupProductionProject(slot: string, base: number, product: number): Promise<void> {
  await fetchJSON(`/api/games/${slot}/production/${base}/${product}/speedup`, { method: 'POST' });
}

export async function getFinancials(slot: string): Promise<Financials> {
  return fetchJSON(`/api/games/${slot}/financials`);
}

export async function updateFinancials(slot: string, updates: Partial<Financials>): Promise<void> {
  await fetchJSON(`/api/games/${slot}/financials`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(updates),
  });
}

export async function healAll(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/heal-all`, { method: 'POST' });
}

export async function equipPersonalArmour(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/personal-armour`, { method: 'POST' });
}

export async function repairRefuelCrafts(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/repair-refuel-crafts`, { method: 'POST' });
}

export async function completeConstructions(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/complete-constructions`, { method: 'POST' });
}

export async function speedupDeliveries(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/speedup-deliveries`, { method: 'POST' });
}

export async function speedupResearch(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/speedup-research`, { method: 'POST' });
}

export async function speedupProduction(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/speedup-production`, { method: 'POST' });
}

export async function saveGame(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/save`, { method: 'POST' });
}

export async function reloadGame(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/reload`, { method: 'POST' });
}

export async function listLocations(slot: string): Promise<GlobeLocation[]> {
  return fetchJSON(`/api/games/${slot}/locations`);
}

export async function getBattleMap(slot: string): Promise<BattleMap> {
  return fetchJSON(`/api/games/${slot}/battle`);
}

/** Returns the sprites of the map, fails if the game graphics are not available. */
export async function getBattleTerrain(slot: string): Promise<BattleTerrain> {
  return fetchJSON(`/api/games/${slot}/battle/terrain`);
}

export function terrainSheetURL(name: string): string {
  return `/api/terrain/${name}.png`;
}

export function unitSheetURL(name: string): string {
  return `/api/unitsprites/${name}.png`;
}

export async function getWorld(): Promise<WorldData> {
  return fetchJSON('/api/world');
}

export function textureURL(set: number, terrain: number): string {
  return `/api/world/textures/${set}/${terrain}.png`;
}

export function baseBitURL(sprite: number): string {
  return `/api/basebits/${sprite}.png`;
}

export function soldierImageURL(name: string): string {
  return `/api/soldier-images/${name}.png`;
}
