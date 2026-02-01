import type {
  GameSummary, GameDetail, SoldierSummary, SoldierDetail,
  BaseSummary, BaseDetail, CraftSummary, TransferSummary, Financials,
  GlobeLocation
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

export async function listCraft(slot: string): Promise<CraftSummary[]> {
  return fetchJSON(`/api/games/${slot}/craft`);
}

export async function listTransfers(slot: string): Promise<TransferSummary[]> {
  return fetchJSON(`/api/games/${slot}/transfers`);
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

export async function completeConstructions(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/complete-constructions`, { method: 'POST' });
}

export async function speedupDeliveries(slot: string): Promise<void> {
  await fetchJSON(`/api/games/${slot}/actions/speedup-deliveries`, { method: 'POST' });
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
