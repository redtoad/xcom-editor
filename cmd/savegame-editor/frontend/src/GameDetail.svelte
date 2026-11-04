<script lang="ts">
  import { onMount, createEventDispatcher } from 'svelte';
  import { getGame, saveGame, reloadGame, healAll, equipPersonalArmour, completeConstructions, speedupDeliveries } from './lib/api';
  import { applyPendingChanges, refreshForms } from './lib/pending';
  import type { GameDetail } from './lib/types';
  import SoldierList from './SoldierList.svelte';
  import BaseList from './BaseList.svelte';
  import CraftList from './CraftList.svelte';
  import TransferList from './TransferList.svelte';
  import FinancialView from './FinancialView.svelte';
  import GlobeView from './GlobeView.svelte';
  import HemisphereMap from './HemisphereMap.svelte';
  import MissionMap from './MissionMap.svelte';

  // Svelte reserves "slot" as a keyword, use a different prop name
  export let slot: string;
  // the savegame was changed on disk (e.g. saved in the game), as reported
  // by the list of savegames in the sidebar
  export let stale = false;

  const dispatch = createEventDispatcher<{ changed: void }>();

  let game: GameDetail | null = null;
  let activeTab = 'soldiers';
  let status = '';
  let error = '';
  // incremented after every action, re-creates the active tab so that it
  // shows the current data
  let version = 0;

  // The sidebar and the detail view poll at different times, the savegame
  // is stale if either of them saw it changed.
  $: changedOnDisk = stale || (game?.stale ?? false);

  const allTabs = [
    { id: 'mission', label: 'Mission' },
    { id: 'soldiers', label: 'Soldiers' },
    { id: 'bases', label: 'Bases' },
    { id: 'craft', label: 'Crafts' },
    { id: 'transfers', label: 'Transfers' },
    { id: 'financials', label: 'Financials' },
    { id: 'map', label: 'Map' },
    { id: 'globe', label: 'Globe' },
  ];
  // the mission map is only available for savegames saved during a mission
  $: tabs = allTabs.filter(t => t.id !== 'mission' || game?.mission);
  $: if (game && !tabs.some(t => t.id === activeTab)) activeTab = 'soldiers';

  async function load() {
    try {
      game = await getGame(slot);
      error = '';
    } catch (e: any) {
      error = e.message;
    }
  }

  onMount(load);
  $: if (slot) load();

  // Apply changes in open edit forms first, so values that were edited but
  // not applied with "Apply Changes" are not lost.
  // The slot is captured first, so selecting another savegame while the
  // changes are applied does not save the wrong game.
  async function saveWithPendingChanges() {
    const target = slot;
    await applyPendingChanges();
    await saveGame(target);
  }

  // After reloading from disk, refresh open edit forms so they neither show
  // nor later re-apply discarded values.
  async function reloadAndRefresh() {
    await reloadGame(slot);
    await refreshForms();
  }

  // Changes to the savegame (e.g. "Heal All") re-create the active tab, so
  // apply changes in open edit forms first instead of discarding them.
  function withPendingChanges(fn: (slot: string) => Promise<void>): () => Promise<void> {
    return async () => {
      const target = slot;
      await applyPendingChanges();
      await fn(target);
    };
  }

  function save() {
    if (changedOnDisk && !confirm('This savegame was changed on disk. Saving overwrites these changes. Save anyway?')) {
      return;
    }
    doAction('Saving', saveWithPendingChanges);
  }

  async function doAction(name: string, fn: () => Promise<void>) {
    try {
      status = name + '...';
      await fn();
      status = name + ' done!';
      await load();
      version++;
      dispatch('changed');
      setTimeout(() => status = '', 2000);
    } catch (e: any) {
      status = '';
      error = e.message;
    }
  }
</script>

{#if error}
  <div class="error">{error}</div>
{/if}

{#if game}
  {#if changedOnDisk}
    <div class="stale-notice">
      <span>This savegame was changed on disk (e.g. saved in the game). Reload it to see the changes; saving in the editor would overwrite them.</span>
      <button class="action reload" on:click={() => doAction('Reloading', reloadAndRefresh)}>Reload</button>
    </div>
  {/if}
  <div class="header">
    <h2>{game.title}</h2>
    <div class="meta">
      <span>{game.time}</span>
      <span><span class="label">Funds&gt;</span> <span class="value">${game.balance?.toLocaleString()}</span></span>
    </div>
    {#if game.mission}
      <div class="mission dialog">
        <strong>Mission in progress</strong>
        <span>Turn {game.mission.turn}</span>
        {#if game.mission.craftName}<span>{game.mission.craftName}</span>{/if}
        <span>X-COM {game.mission.xcomAlive}/{game.mission.xcomTotal}</span>
        <span>Aliens {game.mission.aliensAlive}/{game.mission.aliensTotal}</span>
        <span>Score {game.mission.score}</span>
      </div>
    {:else}
      <div class="mission">On Geoscape &middot; not in a mission</div>
    {/if}
  </div>

  <nav class="tabs">
    {#each tabs as tab}
      <button
        class="tab"
        class:active={activeTab === tab.id}
        on:click={() => activeTab = tab.id}
      >{tab.label}</button>
    {/each}
  </nav>

  <div class="tab-content">
    {#key version}
    {#if activeTab === 'mission'}
      <MissionMap gameSlot={slot} />
    {:else if activeTab === 'soldiers'}
      <SoldierList gameSlot={slot} />
    {:else if activeTab === 'bases'}
      <BaseList gameSlot={slot} />
    {:else if activeTab === 'craft'}
      <CraftList gameSlot={slot} />
    {:else if activeTab === 'transfers'}
      <TransferList gameSlot={slot} />
    {:else if activeTab === 'financials'}
      <FinancialView gameSlot={slot} on:saved={load} />
    {:else if activeTab === 'map'}
      <HemisphereMap gameSlot={slot} />
    {:else if activeTab === 'globe'}
      <GlobeView gameSlot={slot} />
    {/if}
    {/key}
  </div>

  <div class="action-bar">
    {#if activeTab === 'soldiers'}
      <button class="action" on:click={() => doAction('Healing', withPendingChanges(healAll))}>Heal All</button>
      <button class="action" on:click={() => doAction('Equipping', withPendingChanges(equipPersonalArmour))}>All With Personal Armour</button>
    {:else if activeTab === 'bases'}
      <button class="action" on:click={() => doAction('Completing', withPendingChanges(completeConstructions))}>Complete Constructions</button>
    {:else if activeTab === 'transfers'}
      <button class="action" on:click={() => doAction('Speeding up', withPendingChanges(speedupDeliveries))}>Speed Up Deliveries</button>
    {/if}
    <div class="spacer"></div>
    <button class="action reload" on:click={() => doAction('Reloading', reloadAndRefresh)}>Reload</button>
    <button class="action save" on:click={save}>Save</button>
    {#if status}
      <span class="status">{status}</span>
    {/if}
  </div>
{/if}

<style>
  .header {
    margin-bottom: 4px;
  }
  .stale-notice {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
    padding: 6px 8px;
    color: var(--warning);
    background: var(--bg);
    border: 2px solid var(--warning);
  }
  .stale-notice span {
    flex: 1;
  }
  .meta {
    display: flex;
    gap: 24px;
    color: var(--muted);
  }
  .meta .label {
    color: var(--danger);
    text-transform: uppercase;
  }
  .meta .value {
    color: var(--value);
  }
  .mission {
    display: inline-flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-top: 8px;
    color: var(--muted);
    text-transform: uppercase;
  }
  .mission.dialog {
    color: var(--dialog-text);
  }
  .mission.dialog strong {
    color: var(--value);
  }
  .tabs {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
    margin: 12px 0 0 0;
    padding-bottom: 8px;
    border-bottom: 2px solid var(--border);
  }
  .tab.active {
    background: var(--button-dark);
    color: var(--value);
    border-color: var(--bg) var(--button-light) var(--button-light) var(--bg);
  }
  .tab-content {
    margin-top: 12px;
    min-height: 300px;
    /* the pixel font is wide, tables may need to scroll */
    overflow-x: auto;
  }
  .action-bar {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 20px;
    padding-top: 12px;
    border-top: 2px solid var(--border);
    flex-wrap: wrap;
  }
  .action.save {
    background: var(--dialog-border);
    color: var(--bg);
    border-color: var(--dialog-text) var(--dialog) var(--dialog) var(--dialog-text);
  }
  .action.reload {
    background: var(--warning);
    color: var(--bg);
    border-color: #fcd49c #8c5c1c #8c5c1c #fcd49c;
  }
  .spacer { flex: 1; }
  .status {
    color: var(--value);
    text-transform: uppercase;
  }
</style>
