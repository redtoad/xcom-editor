<script lang="ts">
  import { onMount } from 'svelte';
  import { getGame, saveGame, reloadGame, healAll, completeConstructions, speedupDeliveries } from './lib/api';
  import type { GameDetail } from './lib/types';
  import SoldierList from './SoldierList.svelte';
  import BaseList from './BaseList.svelte';
  import CraftList from './CraftList.svelte';
  import TransferList from './TransferList.svelte';
  import FinancialView from './FinancialView.svelte';
  import GlobeView from './GlobeView.svelte';
  import HemisphereMap from './HemisphereMap.svelte';

  // Svelte reserves "slot" as a keyword, use a different prop name
  export let slot: string;

  let game: GameDetail | null = null;
  let activeTab = 'soldiers';
  let status = '';
  let error = '';

  const tabs = [
    { id: 'soldiers', label: 'Soldiers' },
    { id: 'bases', label: 'Bases' },
    { id: 'craft', label: 'Crafts' },
    { id: 'transfers', label: 'Transfers' },
    { id: 'financials', label: 'Financials' },
    { id: 'map', label: 'Map' },
    { id: 'globe', label: 'Globe' },
  ];

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

  async function doAction(name: string, fn: () => Promise<void>) {
    try {
      status = name + '...';
      await fn();
      status = name + ' done!';
      await load();
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
    {#if activeTab === 'soldiers'}
      <SoldierList gameSlot={slot} />
    {:else if activeTab === 'bases'}
      <BaseList gameSlot={slot} />
    {:else if activeTab === 'craft'}
      <CraftList gameSlot={slot} />
    {:else if activeTab === 'transfers'}
      <TransferList gameSlot={slot} />
    {:else if activeTab === 'financials'}
      <FinancialView gameSlot={slot} />
    {:else if activeTab === 'map'}
      <HemisphereMap gameSlot={slot} />
    {:else if activeTab === 'globe'}
      <GlobeView gameSlot={slot} />
    {/if}
  </div>

  <div class="action-bar">
    <button class="action" on:click={() => doAction('Healing', () => healAll(slot))}>Heal All</button>
    <button class="action" on:click={() => doAction('Completing', () => completeConstructions(slot))}>Complete Constructions</button>
    <button class="action" on:click={() => doAction('Speeding up', () => speedupDeliveries(slot))}>Speed Up Deliveries</button>
    <div class="spacer"></div>
    <button class="action reload" on:click={() => doAction('Reloading', () => reloadGame(slot))}>Reload</button>
    <button class="action save" on:click={() => doAction('Saving', () => saveGame(slot))}>Save</button>
    {#if status}
      <span class="status">{status}</span>
    {/if}
  </div>
{/if}

<style>
  .header {
    margin-bottom: 4px;
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
